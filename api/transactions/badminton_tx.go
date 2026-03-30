package transactions

import (
	"context"
	"fmt"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const badmintonSetsToWin = 2
const badmintonMinScore = 21
const badmintonMaxScore = 30

func (store *SQLStore) UpdateBadmintonScoreTx(ctx *gin.Context, matchPublicID, teamPublicID uuid.UUID, setNumber int) (*models.Match, map[string]interface{}, *models.BadmintonSetsPoints, *models.BadmintonScore, error) {
	var badmintonSetScore map[string]interface{}
	var matchResult *models.Match
	var points *models.BadmintonSetsPoints
	var newSet *models.BadmintonScore

	err := store.execTx(ctx, func(q *database.Queries) error {

		var updatedScore *models.BadmintonScore

		// Get match to access internal IDs
		match, err := q.GetMatchModelByPublicId(ctx, matchPublicID)
		if err != nil {
			store.logger.Error("failed to get match: ", err)
			return err
		}

		team, err := q.GetTeamByPublicID(ctx, teamPublicID)
		if err != nil {
			store.logger.Error("failed to get team: ", err)
			return err
		}

		// Increment score for the scoring team
		updatedScore, err = q.UpdateBadmintonScore(ctx, matchPublicID, teamPublicID, setNumber)
		if err != nil {
			store.logger.Error("failed to update badminton score: ", err)
			return err
		}

		if updatedScore != nil {
			pointNumber, err := q.GetBadmintonMaxPointNumber(ctx, updatedScore.MatchID, setNumber)
			if err != nil {
				store.logger.Error("failed to get badminton point number: ", err)
				return err
			}

			// broadcaster to send response to frontend
			points, err = q.AddBadmintonSetsPoints(ctx, updatedScore.MatchID, setNumber, int32(team.ID), updatedScore.HomeScore, updatedScore.AwayScore, *pointNumber)
			if err != nil {
				store.logger.Error("failed to add badminton sets points: ", err)
				return err
			}
		}

		// Check if the current set is finished
		setFinished := isBadmintonSetFinished(updatedScore.HomeScore, updatedScore.AwayScore)
		if !setFinished {
			// Set still in progress — build response and return
			badmintonSetScore = buildBadmintonSetScore(updatedScore, matchPublicID)
			return nil
		} else {
			// Mark the current set as finished
			updatedScore, err = q.UpdateBadmintonSetStatus(ctx, matchPublicID, setNumber, "finished")
			if err != nil {
				store.logger.Error("failed to update badminton set status: ", err)
				return err
			}

			// Get overall sets won by each side
			matchScore, err := q.GetBadmintonMatchScore(ctx, matchPublicID)
			if err != nil {
				store.logger.Error("failed to get match score: ", err)
				return err
			}

			homeSetsWon := derefInt(matchScore.HomeSetsWon)
			awaySetsWon := derefInt(matchScore.AwaySetsWon)

			// Determine if match is decided (best of 3)
			if homeSetsWon >= badmintonSetsToWin || awaySetsWon >= badmintonSetsToWin {
				// Match is over — determine winner
				var winnerTeamID int32
				if homeSetsWon > awaySetsWon {
					winnerTeamID = match.HomeTeamID
				} else {
					winnerTeamID = match.AwayTeamID
				}

				matchResult, err = q.UpdateMatchResult(ctx, int32(match.ID), winnerTeamID)
				if err != nil {
					store.logger.Error("failed to update match result: ", err)
					return err
				}

				// Update badminton player stats for both teams/player
				if err := updateBadmintonStatsOnFinish(ctx, q, store, match, matchPublicID, winnerTeamID, homeSetsWon, awaySetsWon); err != nil {
					store.logger.Error("failed to update badminton player stats: ", err)
					return err
				}
			} else {
				// Match continues — create next set
				nextSetNumber := setNumber + 1
				newSet, err = q.AddBadmintonScore(ctx, int32(match.ID), nextSetNumber)
				if err != nil {
					store.logger.Error("failed to add next set: ", err)
					return err
				}
			}
		}

		// Build response with updated set score
		badmintonSetScore = buildBadmintonSetScore(updatedScore, matchPublicID)
		return nil
	})

	return matchResult, badmintonSetScore, points, newSet, err
}

// isBadmintonSetFinished checks if a badminton set is complete.
// Standard rules: first to 21, win by 2, hard cap at 30.
func isBadmintonSetFinished(homeScore, awayScore int) bool {
	if homeScore == badmintonMaxScore || awayScore == badmintonMaxScore {
		return true
	}

	diff := homeScore - awayScore
	if diff < 0 {
		diff = -diff
	}

	if (homeScore >= badmintonMinScore || awayScore >= badmintonMinScore) && diff >= 2 {
		return true
	}

	return false
}

func buildBadmintonSetScore(score *models.BadmintonScore, matchPublicID uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"public_id":       score.PublicID,
		"match_public_id": matchPublicID,
		"set_number":      score.SetNumber,
		"home_score":      score.HomeScore,
		"away_score":      score.AwayScore,
		"set_status":      score.SetStatus,
	}
}

// derefInt safely dereferences an *int, returning 0 if nil.
func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// updateBadmintonStatsOnFinish calculates and upserts player stats for each player in both teams when a match finishes.
// For singles: each team has 1 player → 1 stats upsert per side
// For doubles: each team has 2 players → 2 stats upserts per side (each player gets individual stats)
func updateBadmintonStatsOnFinish(
	ctx context.Context,
	q *database.Queries,
	store *SQLStore,
	match *models.Match,
	matchPublicID uuid.UUID,
	winnerTeamID int32,
	homeSetsWon, awaySetsWon int,
) error {
	// Get all set scores to calculate total points
	sets, err := q.GetBadmintonMatchSetsScore(ctx, matchPublicID)
	if err != nil {
		return fmt.Errorf("failed to get match sets for stats: %w", err)
	}

	// Calculate total points from all sets
	homePointsScored := 0
	homePointsConceded := 0
	for _, set := range sets {
		homePointsScored += set.HomeScore
		homePointsConceded += set.AwayScore
	}
	awayPointsScored := homePointsConceded
	awayPointsConceded := homePointsScored

	// Determine play_type from match type
	playType := "singles"
	if match.Type == "double" {
		playType = "doubles"
	}

	// Home team win/loss
	homeWon := 0
	homeLost := 0
	homeStreak := 0
	if winnerTeamID == match.HomeTeamID {
		homeWon = 1
		homeStreak = 1
	} else {
		homeLost = 1
	}

	// Away team win/loss
	awayWon := 0
	awayLost := 0
	awayStreak := 0
	if winnerTeamID == match.AwayTeamID {
		awayWon = 1
		awayStreak = 1
	} else {
		awayLost = 1
	}

	// Get players from home team
	homePlayerIDs, err := q.GetPlayerIDsByTeamID(ctx, match.HomeTeamID)
	if err != nil {
		return fmt.Errorf("failed to get home team players: %w", err)
	}

	// Get players from away team
	awayPlayerIDs, err := q.GetPlayerIDsByTeamID(ctx, match.AwayTeamID)
	if err != nil {
		return fmt.Errorf("failed to get away team players: %w", err)
	}

	// UPSERT stats for each player in home team
	for _, playerID := range homePlayerIDs {
		_, err = q.AddOrUpdateBadmintonPlayerStats(ctx, database.AddOrUpdateBadmintonPlayerStatsParams{
			PlayerID:       playerID,
			PlayType:       playType,
			Wins:           homeWon,
			Losses:         homeLost,
			SetsWon:        homeSetsWon,
			SetsLost:       awaySetsWon,
			PointsScored:   homePointsScored,
			PointsConceded: homePointsConceded,
			WinPercentage:  float64(homeWon) * 100,
			Streak:         homeStreak,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert home player %d stats: %w", playerID, err)
		}
	}

	// UPSERT stats for each player in away team
	for _, playerID := range awayPlayerIDs {
		_, err = q.AddOrUpdateBadmintonPlayerStats(ctx, database.AddOrUpdateBadmintonPlayerStatsParams{
			PlayerID:       playerID,
			PlayType:       playType,
			Wins:           awayWon,
			Losses:         awayLost,
			SetsWon:        awaySetsWon,
			SetsLost:       homeSetsWon,
			PointsScored:   awayPointsScored,
			PointsConceded: awayPointsConceded,
			WinPercentage:  float64(awayWon) * 100,
			Streak:         awayStreak,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert away player %d stats: %w", playerID, err)
		}
	}

	return nil
}
