package transactions

import (
	"khelogames/database"
	"khelogames/database/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const badmintonSetsToWin = 2
const badmintonMinScore = 21
const badmintonMaxScore = 30

func (store *SQLStore) UpdateBadmintonScoreTx(ctx *gin.Context, matchPublicID, teamPublicID uuid.UUID, setNumber int) (*models.Match, map[string]interface{}, error) {
	var badmintonSetScore map[string]interface{}
	var matchResult *models.Match

	err := store.execTx(ctx, func(q *database.Queries) error {
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
		updatedScore, err := q.UpdateBadmintonScore(ctx, matchPublicID, teamPublicID, setNumber)
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
			_, err = q.AddBadmintonSetsPoints(ctx, updatedScore.MatchID, setNumber, int32(team.ID), updatedScore.HomeScore, updatedScore.AwayScore, *pointNumber)
			if err != nil {
				store.logger.Error("failed to add badminton sets points: ", err)
				return err
			}
		}

		// Check if the current set is finished
		if updatedScore.HomeScore != 0 && updatedScore.AwayScore != 0 {
			setFinished := isBadmintonSetFinished(updatedScore.HomeScore, updatedScore.AwayScore)
			if !setFinished {
				// Set still in progress — build response and return
				badmintonSetScore = buildBadmintonSetScore(updatedScore, matchPublicID)
				return nil
			}

			// Mark the current set as finished
			_, err = q.UpdateBadmintonSetStatus(ctx, matchPublicID, setNumber, "finished")
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
			} else {
				// Match continues — create next set
				nextSetNumber := setNumber + 1
				_, err = q.AddBadmintonScore(ctx, int32(match.ID), nextSetNumber)
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

	return matchResult, badmintonSetScore, err
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
