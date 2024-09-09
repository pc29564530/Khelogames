package football

import (
	"encoding/json"
	"fmt"
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addFootballIncidentsRequest struct {
	MatchID      int64  `json:"match_id"`
	TeamID       int64  `json:"team_id"`
	Periods      string `json:"periods"`
	IncidentType string `json:"incident_type"`
	IncidentTime int64  `json:"incident_time"`
	PlayerID     int64  `json:"player_id"`
	Description  string `json:"description"`
}

func (s *FootballServer) AddFootballIncidents(ctx *gin.Context) {
	var req addFootballIncidentsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.CreateFootballIncidentsParams{
		MatchID:      req.MatchID,
		TeamID:       req.TeamID,
		Periods:      req.Periods,
		IncidentType: req.IncidentType,
		IncidentTime: req.IncidentTime,
		Description:  req.Description,
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("failed to start transcation: ", err)
		return
	}

	incidents, err := s.store.CreateFootballIncidents(ctx, arg)
	if err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create football incidents: ", err)
		return
	}

	incidentPlayerArg := db.AddFootballIncidentPlayerParams{
		IncidentID: incidents.ID,
		PlayerID:   req.PlayerID,
	}

	_, err = s.store.AddFootballIncidentPlayer(ctx, incidentPlayerArg)
	if err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create football incidents: ", err)
		return
	}

	if incidents.IncidentType == "goal" {
		if incidents.Periods == "first_half" {
			argGoalScore := db.UpdateFirstHalfScoreParams{
				FirstHalf: 1,
				MatchID:   incidents.MatchID,
				TeamID:    incidents.TeamID,
			}

			fmt.Println("First Half: ", argGoalScore)

			_, err := s.store.UpdateFirstHalfScore(ctx, argGoalScore)
			if err != nil {
				tx.Rollback()
				s.logger.Error("Failed to update football score: ", err)
				return
			}
		} else {
			argGoalScore := db.UpdateSecondHalfScoreParams{
				SecondHalf: 1,
				MatchID:    incidents.MatchID,
				TeamID:     incidents.TeamID,
			}

			fmt.Println("Second Half: ", argGoalScore)

			_, err := s.store.UpdateSecondHalfScore(ctx, argGoalScore)
			if err != nil {
				tx.Rollback()
				s.logger.Error("Failed to update football score: ", err)
				return
			}
		}
	}

	//commit the transcation if all operation are successfull
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		s.logger.Error("unable to commit the transcation: ", err)
		return
	}

	s.logger.Info("successfully update the add football incident ")

	ctx.JSON(http.StatusAccepted, incidents)
}

type addFootballIncidentsSubsRequest struct {
	MatchID      int64  `json:"match_id"`
	TeamID       int64  `json:"team_id"`
	Periods      string `json:"periods"`
	IncidentType string `json:"incident_type"`
	IncidentTime int64  `json:"incident_time"`
	Description  string `json:"description"`
	PlayerInID   int64  `json:"player_in_id"`
	PlayerOutID  int64  `json:"player_out_in"`
}

func (s *FootballServer) AddFootballIncidentsSubs(ctx *gin.Context) {
	var req addFootballIncidentsSubsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.CreateFootballIncidentsParams{
		MatchID:      req.MatchID,
		TeamID:       req.TeamID,
		Periods:      req.Periods,
		IncidentType: req.IncidentType,
		IncidentTime: req.IncidentTime,
		Description:  req.Description,
	}

	incidents, err := s.store.CreateFootballIncidents(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create football incidents: ", err)
		return
	}

	incidentPlayerArg := db.ADDFootballSubsPlayerParams{
		IncidentID:  incidents.ID,
		PlayerInID:  req.PlayerInID,
		PlayerOutID: req.PlayerOutID,
	}

	_, err = s.store.ADDFootballSubsPlayer(ctx, incidentPlayerArg)
	if err != nil {
		s.logger.Error("Failed to create football incidents: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, incidents)
}

type getFootballIncidentsRequest struct {
	MatchID int64 `json:"match_id" form:"match_id"`
}

func (s *FootballServer) GetFootballIncidents(ctx *gin.Context) {
	var req getFootballIncidentsRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	s.logger.Debug("Successfully bind the req: ", req)

	response, err := s.store.GetFootballIncidentWithPlayer(ctx, req.MatchID)
	if err != nil {
		s.logger.Error("Failed to get football incidents: ", err)
		return
	}

	var incidents []map[string]interface{}

	for _, incident := range response {

		if incident.IncidentType == "substitutions" {

			var data map[string]interface{}
			tt := (incident.Players).([]byte)
			err := json.Unmarshal(tt, &data)
			if err != nil {
				s.logger.Error("unable to unmarshal incident player: ", err)
			}

			playerInData := data["player_in"].(map[string]interface{})
			playerOutData := data["player_out"].(map[string]interface{})
			incidentDataMap := map[string]interface{}{
				"id":            incident.ID,
				"match_id":      incident.MatchID,
				"team_id":       incident.TeamID,
				"incident_type": incident.IncidentType,
				"incident_time": incident.IncidentTime,
				"description":   incident.Description,
				"player_in": map[string]interface{}{
					"id":         playerInData["id"],
					"name":       playerInData["name"],
					"slug":       playerInData["slug"],
					"short_name": playerInData["short_name"],
					"positions":  playerInData["positions"],
					"country":    playerInData["country"],
					"media_url":  playerInData["media_url"],
				},
				"player_out": map[string]interface{}{
					"id":         playerOutData["id"],
					"name":       playerOutData["name"],
					"slug":       playerOutData["slug"],
					"short_name": playerOutData["short_name"],
					"positions":  playerOutData["positions"],
					"country":    playerOutData["country"],
					"media_url":  playerOutData["media_url"],
				},
			}

			incidents = append(incidents, incidentDataMap)

		} else {

			var data map[string]interface{}
			tt := (incident.Players).([]byte)
			err := json.Unmarshal(tt, &data)
			if err != nil {
				s.logger.Error("unable to unmarshal incident player: ", err)
			}

			playerData := data["player"].(map[string]interface{})
			incidentDataMap := map[string]interface{}{
				"id":            incident.ID,
				"match_id":      incident.MatchID,
				"team_id":       incident.TeamID,
				"incident_type": incident.IncidentType,
				"incident_time": incident.IncidentTime,
				"description":   incident.Description,
				"player": map[string]interface{}{
					"id":         playerData["id"],
					"name":       playerData["name"],
					"slug":       playerData["slug"],
					"short_name": playerData["short_name"],
					"positions":  playerData["positions"],
					"country":    playerData["country"],
					"media_url":  playerData["media_url"],
				},
			}
			incidents = append(incidents, incidentDataMap)
		}
	}

	match, err := s.store.GetMatchByMatchID(ctx, req.MatchID)
	if err != nil {
		s.logger.Error("Failed to get match data: ", err)
		return
	}

	var matchIncidents []map[string]interface{}

	matchDetail := map[string]interface{}{
		"id":              match.ID,
		"tournament_id":   match.TournamentID,
		"home_team_id":    match.HomeTeamID,
		"away_team_id":    match.AwayTeamID,
		"status_code":     match.StatusCode,
		"start_timestamp": match.StartTimestamp,
		"end_timestamp":   match.EndTimestamp,
		"type":            match.Type,
	}
	matchIncidents = append(matchIncidents, map[string]interface{}{
		"match": matchDetail,
	})

	matchIncidents = append(matchIncidents, map[string]interface{}{
		"incidents": incidents,
	})

	s.logger.Info("Successfully get match incidents")
	ctx.JSON(http.StatusAccepted, matchIncidents)
}
