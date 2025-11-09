package transactions

import (
	footballhelper "khelogames/api/sports/football_helper"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (store *SQLStore) AddFootballIncidentsTx(ctx *gin.Context, arg database.CreateFootballIncidentsParams, playerPublicID uuid.UUID) (map[string]interface{}, error) {
	var incidentData map[string]interface{}
	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		incidents, err := q.CreateFootballIncidents(ctx, arg)
		if err != nil {
			store.logger.Error("Failed to create football incidents: ", err)
			return err
		}

		store.logger.Info("successfully created the incident: ", incidents)

		var playerData *models.Player
		incidentID := incidents.ID
		if incidents.IncidentType != "period" {

			addedPlayer, err := q.AddFootballIncidentPlayer(ctx, incidentID, playerPublicID)
			if err != nil {
				store.logger.Error("Failed to add football incidents player: ", err)
				return err
			}

			playerData, err = q.GetPlayerByID(ctx, int64(addedPlayer.PlayerID))
			if err != nil {
				store.logger.Error("Failed to add get player: ", err)
				return err
			}

			statsUpdate := footballhelper.GetStatisticsUpdateFromIncident(incidents.IncidentType)

			statsArg := database.UpdateFootballStatisticsParams{
				MatchID:         incidents.MatchID,
				TeamID:          *incidents.TeamID,
				ShotsOnTarget:   statsUpdate.ShotsOnTarget,
				TotalShots:      statsUpdate.TotalShots,
				CornerKicks:     statsUpdate.CornerKicks,
				Fouls:           statsUpdate.Fouls,
				GoalkeeperSaves: statsUpdate.GoalkeeperSaves,
				FreeKicks:       statsUpdate.FreeKicks,
				YellowCards:     statsUpdate.YellowCards,
				RedCards:        statsUpdate.RedCards,
			}

			_, err = q.UpdateFootballStatistics(ctx, statsArg)
			if err != nil {
				store.logger.Error("Failed to update statistics: ", err)
				return err
			}

			//Handle goals, penalty_shootout and penalty
			switch incidents.IncidentType {
			case "goal", "penalty":
				if incidents.Periods == "first_half" {
					argGoalScore := database.UpdateFirstHalfScoreParams{
						FirstHalf: 1,
						MatchID:   incidents.MatchID,
						TeamID:    *incidents.TeamID,
					}

					firstHalfData, err := q.UpdateFirstHalfScore(ctx, argGoalScore)
					if err != nil {
						store.logger.Error("Failed to update football score: ", err)
						return err
					}

					firstHalf := map[string]interface{}{
						"id":               firstHalfData.ID,
						"public_id":        firstHalfData.PublicID,
						"match_id":         firstHalfData.MatchID,
						"team_id":          firstHalfData.TeamID,
						"first_half":       firstHalfData.FirstHalf,
						"second_half":      firstHalfData.SecondHalf,
						"goals":            firstHalfData.Goals,
						"penalty_shootout": firstHalfData.PenaltyShootOut,
					}

					if store.scoreBroadcaster != nil {
						err := store.scoreBroadcaster.BroadcastFootballEvent(ctx, "UPDATE_FOOTBALL_SCORE", firstHalf)
						if err != nil {
							store.logger.Error("Failed to broadcast cricket event: ", err)
						}
					}
				} else if incidents.Periods == "second_half" {
					argGoalScore := database.UpdateSecondHalfScoreParams{
						SecondHalf: 1,
						MatchID:    incidents.MatchID,
						TeamID:     *incidents.TeamID,
					}

					secondHalfData, err := q.UpdateSecondHalfScore(ctx, argGoalScore)
					if err != nil {
						store.logger.Error("Failed to update football score: ", err)
						return err
					}

					secondHalf := map[string]interface{}{
						"id":               secondHalfData.ID,
						"public_id":        secondHalfData.PublicID,
						"match_id":         secondHalfData.MatchID,
						"team_id":          secondHalfData.TeamID,
						"first_half":       secondHalfData.FirstHalf,
						"second_half":      secondHalfData.SecondHalf,
						"goals":            secondHalfData.Goals,
						"penalty_shootout": secondHalfData.PenaltyShootOut,
					}

					if store.scoreBroadcaster != nil {
						err := store.scoreBroadcaster.BroadcastFootballEvent(ctx, "UPDATE_FOOTBALL_SCORE", secondHalf)
						if err != nil {
							store.logger.Error("Failed to broadcast cricket event: ", err)
						}
					}
				}
			case "penalty_shootout":
				if incidents.PenaltyShootoutScored {
					penaltyData, err := q.UpdatePenaltyShootoutScore(ctx, incidents.MatchID, *incidents.TeamID)
					if err != nil {
						store.logger.Error("Failed to update penalty shootout score: ", err)
						return err
					}
					penaltyShootout := map[string]interface{}{
						"id":               penaltyData.ID,
						"public_id":        penaltyData.PublicID,
						"match_id":         penaltyData.MatchID,
						"team_id":          penaltyData.TeamID,
						"first_half":       penaltyData.FirstHalf,
						"second_half":      penaltyData.SecondHalf,
						"goals":            penaltyData.Goals,
						"penalty_shootout": penaltyData.PenaltyShootOut,
					}
					if store.scoreBroadcaster != nil {
						err := store.scoreBroadcaster.BroadcastFootballEvent(ctx, "UPDATE_FOOTBALL_SCORE", penaltyShootout)
						if err != nil {
							store.logger.Error("Failed to broadcast cricket event: ", err)
						}
					}
				}
			default:
				store.logger.Errorf("Failed to found the incident type ")
			}

			currentMatch, err := q.GetMatchByPublicId(ctx, arg.MatchPublicID, 1)
			if err != nil {
				store.logger.Error("Failed to get match by public id: ", err)
			}

			homeScore := currentMatch["homeScore"].(map[string]interface{})
			awayScore := currentMatch["awayScore"].(map[string]interface{})

			incidentData = map[string]interface{}{
				"id":                      incidents.ID,
				"public_id":               incidents.PublicID,
				"match_id":                incidents.MatchID,
				"team_id":                 incidents.TeamID,
				"periods":                 incidents.Periods,
				"incident_type":           incidents.IncidentType,
				"incident_time":           incidents.IncidentTime,
				"description":             incidents.Description,
				"penalty_shootout_scored": incidents.PenaltyShootoutScored,
				"tournament_id":           incidents.TournamentID,
				"created_at":              incidents.CreatedAt,
				"player": map[string]interface{}{
					"id":         playerData.ID,
					"public_id":  playerData.PublicID,
					"user_id":    playerData.UserID,
					"name":       playerData.Name,
					"slug":       playerData.Slug,
					"short_name": playerData.Slug,
					"positions":  playerData.Positions,
					"country":    playerData.Country,
					"media_url":  playerData.MediaUrl,
				},
				"awayScore": map[string]interface{}{
					"goals": awayScore["goals"],
				},
				"homeScore": map[string]interface{}{
					"goals": homeScore["goals"],
				},
			}

		} else {
			incidentData = map[string]interface{}{
				"id":                      incidents.ID,
				"public_id":               incidents.PublicID,
				"match_id":                incidents.MatchID,
				"team_id":                 incidents.TeamID,
				"periods":                 incidents.Periods,
				"incident_type":           incidents.IncidentType,
				"incident_time":           incidents.IncidentTime,
				"description":             incidents.Description,
				"penalty_shootout_scored": incidents.PenaltyShootoutScored,
				"tournament_id":           incidents.TournamentID,
				"created_at":              incidents.CreatedAt,
			}
		}
		return err
	})
	return incidentData, err
}

func (store *SQLStore) AddFootballIncidentsSubsTx(ctx *gin.Context, matchPublicID, teamPublicID uuid.UUID, periods, incidentType string, incidentTime int, description string, playerInPublicID, playerOutPublicID uuid.UUID) (*map[string]interface{}, error) {
	var incidentData map[string]interface{}
	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		arg := database.CreateFootballIncidentsParams{
			MatchPublicID: matchPublicID,
			TeamPublicID:  &teamPublicID,
			Periods:       periods,
			IncidentType:  incidentType,
			IncidentTime:  incidentTime,
			Description:   description,
		}

		incidents, err := store.CreateFootballIncidents(ctx, arg)
		if err != nil {
			store.logger.Error("Failed to create football incidents: ", err)
			return err
		}

		data, err := store.ADDFootballSubsPlayer(ctx, incidents.PublicID, playerInPublicID, playerOutPublicID)
		if err != nil {
			store.logger.Error("Failed to create football incidents: ", err)
			return err
		}

		subsData := *data

		incidentData = map[string]interface{}{
			"id":                      incidents.ID,
			"public_id":               incidents.PublicID,
			"match_id":                incidents.MatchID,
			"team_id":                 incidents.TeamID,
			"periods":                 incidents.Periods,
			"incident_type":           incidents.IncidentType,
			"incident_time":           incidents.IncidentTime,
			"description":             incidents.Description,
			"penalty_shootout_scored": incidents.PenaltyShootoutScored,
			"tournament_id":           incidents.TournamentID,
			"created_at":              incidents.CreatedAt,
			"player_in":               subsData["player_id"],
			"player_out":              subsData["player_out"],
		}
		return err
	})
	return &incidentData, err
}
