package transcation_setup

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/google/uuid"
)

func (store *SQLStore) AddFootballIncidentsTx(ctx context.Context, argIncidents database.CreateFootballIncidentsParams, playerPublicID uuid.UUID) (models.FootballIncident, error) {
	var footballIncidents models.FootballIncident
	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		incidents, err := q.CreateFootballIncidents(ctx, argIncidents)
		if err != nil {

			store.logger.Error("Failed to create football incidents: ", err)
			return err
		}

		store.logger.Info("successfully created the incident: ", incidents)

		if incidents.IncidentType != "period" {

			_, err = q.AddFootballIncidentPlayer(ctx, incidents.PublicID, playerPublicID)
			if err != nil {

				store.logger.Error("Failed to create football incidents: ", err)
				return err
			}

			statsUpdate := GetStatisticsUpdateFromIncident(incidents.IncidentType)

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

				_, err := q.UpdateFirstHalfScore(ctx, argGoalScore)
				if err != nil {

					store.logger.Error("Failed to update football score: ", err)
					return err
				}
			} else if incidents.Periods == "second_half" {
				argGoalScore := database.UpdateSecondHalfScoreParams{
					SecondHalf: 1,
					MatchID:    incidents.MatchID,
					TeamID:     *incidents.TeamID,
				}

				_, err := q.UpdateSecondHalfScore(ctx, argGoalScore)
				if err != nil {

					store.logger.Error("Failed to update football score: ", err)
					return err
				}
			}
		case "penalty_shootout":
			if incidents.PenaltyShootoutScored {
				_, err := q.UpdatePenaltyShootoutScore(ctx, incidents.MatchID, *incidents.TeamID)
				if err != nil {

					store.logger.Error("Failed to update penalty shootout score: ", err)
					return err
				}
			}
		}
		return err
	})
	return footballIncidents, err
}

func (store *SQLStore) AddFootballIncidentsSubs(ctx context.Context, matchPublicID, teamPublicID uuid.UUID, periods, incidentType string, incidentTime int64, description string, playerInPublicID, playerOutPublicID uuid.UUID) (*models.FootballIncident, error) {
	var footballIncidents *models.FootballIncident
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

		_, err = store.ADDFootballSubsPlayer(ctx, incidents.PublicID, playerInPublicID, playerOutPublicID)
		if err != nil {
			store.logger.Error("Failed to create football incidents: ", err)
			return err
		}
		return err
	})
	return footballIncidents, err
}
