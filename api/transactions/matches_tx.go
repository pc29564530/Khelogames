package transactions

import (
	"context"
	"fmt"
	"khelogames/database"
	"khelogames/database/models"

	footballhelper "khelogames/api/sports/football_helper"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// func footballhelper.GetInt32(v interface{}) int32 {
// 	switch val := v.(type) {
// 	case nil:
// 		return 0
// 	case int:
// 		return int32(val)
// 	case int32:
// 		return val
// 	case int64:
// 		return int32(val)
// 	case float32:
// 		return int32(val)
// 	case float64:
// 		return int32(val)
// 	default:
// 		return 0
// 	}
// }

// Update match status transaction
func (store *SQLStore) UpdateMatchStatusTx(ctx *gin.Context, matchPublicID uuid.UUID, statusCode string, gameID models.Game) (models.Match, error) {
	var updatedMatchData models.Match

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		updatedMatchData, err = q.UpdateMatchStatus(ctx, matchPublicID, statusCode)
		if err != nil {
			store.logger.Error("Unable to update match status: ", err)
			return err
		}

		// Handle status-specific logic
		switch updatedMatchData.StatusCode {
		case "finished":
			if gameID.Name == "football" {
				_, err = q.AddORUpdateFootballPlayerStats(ctx, matchPublicID)
				if err != nil {
					return fmt.Errorf("Failed to update player stats: ", err)
				}
				if err := UpdateFootballStatusCode(ctx, updatedMatchData, gameID.ID, q, store); err != nil {
					return fmt.Errorf("Failed to update football status code: ", err)
				}

				homeIncident, err := q.GetFootballIncidentByTeam(ctx, updatedMatchData.ID, updatedMatchData.HomeTeamID)
				if err != nil {
					return fmt.Errorf("Failed to football incident by team: ", err)
				}

				awayIncident, err := q.GetFootballIncidentByTeam(ctx, updatedMatchData.ID, updatedMatchData.AwayTeamID)
				if err != nil {
					return fmt.Errorf("Failed to football incident by team: ", err)
				}

				var homeCurrentStats map[string]interface{}
				var awayCurrentStats map[string]interface{}
				for _, incident := range *homeIncident {
					homeCurrentStats = footballhelper.GetStatisticsUpdateFromIncident(homeCurrentStats, incident.IncidentType)
				}

				_, err = q.UpdateFootballStatistics(ctx,
					int32(updatedMatchData.ID),
					updatedMatchData.HomeTeamID,
					footballhelper.GetInt32(homeCurrentStats["shot_on_target"]),
					footballhelper.GetInt32(homeCurrentStats["total_shots"]),
					footballhelper.GetInt32(homeCurrentStats["corner_kicks"]),
					footballhelper.GetInt32(homeCurrentStats["fouls"]),
					footballhelper.GetInt32(homeCurrentStats["goal_keeper_saves"]),
					footballhelper.GetInt32(homeCurrentStats["free_kicks"]),
					footballhelper.GetInt32(homeCurrentStats["yellow_cards"]),
					footballhelper.GetInt32(homeCurrentStats["red_cards"]),
				)
				if err != nil {
					return fmt.Errorf("Failed to update football statistics: ", err)
				}

				for _, incident := range *awayIncident {
					awayCurrentStats = footballhelper.GetStatisticsUpdateFromIncident(awayCurrentStats, incident.IncidentType)
				}

				_, err = q.UpdateFootballStatistics(ctx,
					int32(updatedMatchData.ID),
					updatedMatchData.AwayTeamID,
					footballhelper.GetInt32(awayCurrentStats["shot_on_target"]),
					footballhelper.GetInt32(awayCurrentStats["total_shots"]),
					footballhelper.GetInt32(awayCurrentStats["corner_kicks"]),
					footballhelper.GetInt32(awayCurrentStats["fouls"]),
					footballhelper.GetInt32(awayCurrentStats["goal_keeper_saves"]),
					footballhelper.GetInt32(awayCurrentStats["free_kicks"]),
					footballhelper.GetInt32(awayCurrentStats["yellow_cards"]),
					footballhelper.GetInt32(awayCurrentStats["red_cards"]),
				)
				if err != nil {
					return fmt.Errorf("Failed to update football statistics: ", err)
				}

			} else if gameID.Name == "cricket" {
				if err := UpdateCricketStatusCode(ctx, updatedMatchData, gameID.ID, q, store); err != nil {
					return fmt.Errorf("Failed to update cricket status code: ", err)
				}
			}

		case "in_progress":
			if gameID.Name == "football" {
				if err := UpdateFootballStatusCode(ctx, updatedMatchData, gameID.ID, q, store); err != nil {
					return fmt.Errorf("Failed to initialize the football score: ", err)
				}
			} else if gameID.Name == "cricket" {
				if err := UpdateCricketStatusCode(ctx, updatedMatchData, gameID.ID, q, store); err != nil {
					return fmt.Errorf("Failed to initialize the cricket score: ", err)
				}
			}
		}
		return err
	})
	return updatedMatchData, err
}

func UpdateFootballStatusCode(ctx context.Context, updatedMatchData models.Match, gameID int64, q *database.Queries, store *SQLStore) error {
	var ct *gin.Context

	if updatedMatchData.StatusCode == "in_progress" {
		var penaltyShootOut *int
		argAway := database.NewFootballScoreParams{
			MatchID:         int32(updatedMatchData.ID),
			TeamID:          int32(updatedMatchData.AwayTeamID),
			FirstHalf:       0,
			SecondHalf:      0,
			Goals:           0,
			PenaltyShootOut: penaltyShootOut,
		}

		awayScoreData, err := q.NewFootballScore(ctx, argAway)
		if err != nil {
			store.logger.Error("unable to add the football match score: ", err)
			return err
		}

		awayScore := map[string]interface{}{
			"id":               awayScoreData.ID,
			"public_id":        awayScoreData.PublicID,
			"match_id":         awayScoreData.MatchID,
			"team_id":          awayScoreData.TeamID,
			"first_half":       awayScoreData.FirstHalf,
			"second_half":      awayScoreData.SecondHalf,
			"penalty_shootout": awayScoreData.PenaltyShootOut,
		}

		if store.scoreBroadcaster != nil {
			err := store.scoreBroadcaster.BroadcastTournamentEvent(ct, "ADD_FOOTBALL_SCORE", awayScore)
			if err != nil {
				store.logger.Error("Failed to broadcast cricket event: ", err)
			}
		}

		argHome := database.NewFootballScoreParams{
			MatchID:         int32(updatedMatchData.ID),
			TeamID:          int32(updatedMatchData.HomeTeamID),
			FirstHalf:       0,
			SecondHalf:      0,
			Goals:           0,
			PenaltyShootOut: penaltyShootOut,
		}

		homeScoreData, err := q.NewFootballScore(ctx, argHome)
		if err != nil {
			store.logger.Error("unable to add the football match score: ", err)
			return err
		}

		homeScore := map[string]interface{}{
			"id":               homeScoreData.ID,
			"public_id":        homeScoreData.PublicID,
			"match_id":         homeScoreData.MatchID,
			"team_id":          homeScoreData.TeamID,
			"first_half":       homeScoreData.FirstHalf,
			"second_half":      homeScoreData.SecondHalf,
			"penalty_shootout": homeScoreData.PenaltyShootOut,
		}

		if store.scoreBroadcaster != nil {
			err := store.scoreBroadcaster.BroadcastTournamentEvent(ct, "ADD_FOOTBALL_SCORE", homeScore)
			if err != nil {
				store.logger.Error("Failed to broadcast cricket event: ", err)
			}
		}

		argStatisticsHome := database.CreateFootballStatisticsParams{
			MatchID:         int32(updatedMatchData.ID),
			TeamID:          int32(updatedMatchData.HomeTeamID),
			ShotsOnTarget:   0,
			TotalShots:      0,
			CornerKicks:     0,
			Fouls:           0,
			GoalkeeperSaves: 0,
			FreeKicks:       0,
			YellowCards:     0,
			RedCards:        0,
		}

		argStatisticsAway := database.CreateFootballStatisticsParams{
			MatchID:         int32(updatedMatchData.ID),
			TeamID:          int32(updatedMatchData.AwayTeamID),
			ShotsOnTarget:   0,
			TotalShots:      0,
			CornerKicks:     0,
			Fouls:           0,
			GoalkeeperSaves: 0,
			FreeKicks:       0,
			YellowCards:     0,
			RedCards:        0,
		}

		_, err = q.CreateFootballStatistics(ctx, argStatisticsHome)
		if err != nil {
			store.logger.Error("Failed to add the football statistics: ", err)
			return err
		}

		_, err = q.CreateFootballStatistics(ctx, argStatisticsAway)
		if err != nil {
			store.logger.Error("Failed to add the football statistics: ", err)
			return err
		}
	} else if updatedMatchData.StatusCode == "finished" {

		argAway := database.GetFootballScoreParams{
			MatchID: updatedMatchData.ID,
			TeamID:  int64(updatedMatchData.AwayTeamID),
		}

		awayScore, err := q.GetFootballScore(ctx, argAway)
		if err != nil {
			store.logger.Error("Failed to get away score: ", err)
			return err
		}

		argHome := database.GetFootballScoreParams{
			MatchID: updatedMatchData.ID,
			TeamID:  int64(updatedMatchData.HomeTeamID),
		}

		homeScore, err := q.GetFootballScore(ctx, argHome)
		if err != nil {
			store.logger.Error("Failed to get away score: ", err)
			return err
		}

		if awayScore.Goals > homeScore.Goals {
			_, err := q.UpdateMatchResult(ctx, int32(updatedMatchData.ID), int32(updatedMatchData.AwayTeamID))
			if err != nil {
				store.logger.Error("Failed to update match result: ", err)
				return err
			}
			_, err = q.UpdateFootballStanding(ctx, int64(updatedMatchData.TournamentID), int64(updatedMatchData.HomeTeamID))
			if err != nil {
				store.logger.Error("Failed to update football standing: ", err)
				return err
			}

			_, err = q.UpdateFootballStanding(ctx, int64(updatedMatchData.TournamentID), int64(updatedMatchData.AwayTeamID))
			if err != nil {
				store.logger.Error("Failed to update football standing: ", err)
				return err
			}
		} else if homeScore.Goals > awayScore.Goals {
			_, err := q.UpdateMatchResult(ctx, int32(updatedMatchData.ID), int32(updatedMatchData.HomeTeamID))
			if err != nil {
				store.logger.Error("Failed to update match result: ", err)
				return err
			}

			_, err = q.UpdateFootballStanding(ctx, int64(updatedMatchData.TournamentID), int64(updatedMatchData.HomeTeamID))
			if err != nil {
				store.logger.Error("Failed to update football standing: ", err)
				return err
			}

			_, err = q.UpdateFootballStanding(ctx, int64(updatedMatchData.TournamentID), int64(updatedMatchData.AwayTeamID))
			if err != nil {
				store.logger.Error("Failed to update football standing: ", err)
				return err
			}
		}
	}
	return nil
}

func UpdateCricketStatusCode(ctx context.Context, updatedMatchData models.Match, gameID int64, q *database.Queries, store *SQLStore) error {
	if updatedMatchData.StatusCode == "finished" {

		awayScore, err := q.GetCricketScore(ctx, int32(updatedMatchData.ID), int32(updatedMatchData.AwayTeamID))
		if err != nil {
			store.logger.Error("Failed to get away score: ", err)
			return err
		}

		homeScore, err := q.GetCricketScore(ctx, int32(updatedMatchData.ID), int32(updatedMatchData.HomeTeamID))
		if err != nil {
			store.logger.Error("Failed to get away score: ", err)
			return err
		}

		if awayScore.Score > homeScore.Score {
			_, err := q.UpdateMatchResult(ctx, int32(updatedMatchData.ID), int32(updatedMatchData.AwayTeamID))
			if err != nil {
				store.logger.Error("Failed to update match result: ", err)
				return err
			}
			_, err = q.UpdateCricketStanding(ctx, int32(updatedMatchData.TournamentID), int32(updatedMatchData.AwayTeamID))
			if err != nil {
				store.logger.Error("Failed to update tournament standing: ", err)
				return err
			}
			_, err = q.UpdateCricketStanding(ctx, int32(updatedMatchData.TournamentID), int32(updatedMatchData.HomeTeamID))
			if err != nil {
				store.logger.Error("Failed to update tournament standing: ", err)
				return err
			}
		} else if homeScore.Score > awayScore.Score {
			_, err := q.UpdateMatchResult(ctx, int32(updatedMatchData.ID), int32(updatedMatchData.HomeTeamID))
			if err != nil {
				store.logger.Error("Failed to update match result: ", err)
				return err
			}
			_, err = q.UpdateCricketStanding(ctx, int32(updatedMatchData.TournamentID), int32(updatedMatchData.AwayTeamID))
			if err != nil {
				store.logger.Error("Failed to update tournament standing: ", err)
				return err
			}
			_, err = q.UpdateCricketStanding(ctx, int32(updatedMatchData.TournamentID), int32(updatedMatchData.HomeTeamID))
			if err != nil {
				store.logger.Error("Failed to update tournament standing: ", err)
				return err
			}
		}
	}
	return nil
}
