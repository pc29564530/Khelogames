package transactions

import (
	footballhelper "khelogames/api/sports/football_helper"
	"khelogames/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (store *SQLStore) AddFootballIncidentsTx(ctx *gin.Context, arg database.CreateFootballIncidentsParams, playerPublicID uuid.UUID) (map[string]interface{}, error) {
	var incidentData map[string]interface{}
	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		incidents, err := store.CreateFootballIncidents(ctx, arg)
		if err != nil {
			store.logger.Error("Failed to create football incidents: ", err)
			return err
		}

		store.logger.Info("successfully created the incident: ", incidents)

		var playerData map[string]interface{}

		if incidents.IncidentType != "period" {

			data, err := store.AddFootballIncidentPlayer(ctx, incidents.PublicID, playerPublicID)
			if err != nil {
				store.logger.Error("Failed to create football incidents: ", err)
				return err
			}

			playerData = *data

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

			_, err = store.UpdateFootballStatistics(ctx, statsArg)
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

					firstHalfData, err := store.UpdateFirstHalfScore(ctx, argGoalScore)
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

					secondHalfData, err := store.UpdateSecondHalfScore(ctx, argGoalScore)
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
					penaltyData, err := store.UpdatePenaltyShootoutScore(ctx, incidents.MatchID, *incidents.TeamID)
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
						"penalty_shootout": penaltyData.PenaltyShootOut,
					}
					if store.scoreBroadcaster != nil {
						err := store.scoreBroadcaster.BroadcastFootballEvent(ctx, "UPDATE_FOOTBALL_SCORE", penaltyShootout)
						if err != nil {
							store.logger.Error("Failed to broadcast cricket event: ", err)
						}
					}
				}
			}

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
					"id":         playerData["id"],
					"public_id":  playerData["public_id"],
					"user_id":    playerData["user_id"],
					"name":       playerData["name"],
					"slug":       playerData["slug"],
					"short_name": playerData["short_name"],
					"positions":  playerData["positions"],
					"country":    playerData["country"],
					"media_url":  playerData["media_url"],
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

func (store *SQLStore) AddFootballIncidentsSubsTx(ctx *gin.Context, matchPublicID, teamPublicID uuid.UUID, periods, incidentType string, incidentTime int64, description string, playerInPublicID, playerOutPublicID uuid.UUID) (*map[string]interface{}, error) {
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
