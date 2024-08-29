package football

import (
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addFootballIncidentsRequest struct {
	MatchID      int64  `json:"match_id"`
	TeamID       int64  `json:"team_id"`
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
		IncidentType: req.IncidentType,
		IncidentTime: req.IncidentTime,
		PlayerID:     req.PlayerID,
		Description:  req.Description,
	}

	incidents, err := s.store.CreateFootballIncidents(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create football incidents: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, incidents)
}

type getFootballIncidentsRequest struct {
	MatchID int64 `json:"match_id"`
}

func (s *FootballServer) GetFootballIncidents(ctx *gin.Context) {
	var req getFootballIncidentsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	incidents, err := s.store.GetFootballIncidents(ctx, req.MatchID)
	if err != nil {
		s.logger.Error("Failed to get football incidents: ", err)
		return
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

	for _, incident := range incidents {
		team, err := s.store.GetTeam(ctx, incident.TeamID)
		if err != nil {
			s.logger.Error("Failed to get teams: ", err)
			return
		}
		player, err := s.store.GetPlayer(ctx, incident.PlayerID)
		if err != nil {
			s.logger.Error("Failed to get football incidents: ", err)
			return
		}

		matchIncident := map[string]interface{}{
			"team": map[string]interface{}{
				"id":        team.ID,
				"name":      team.Name,
				"slug":      team.Slug,
				"shortName": team.Shortname,
				"gender":    team.Gender,
				"national":  team.National,
				"country":   team.Country,
				"type":      team.Type},
			"player": map[string]interface{}{
				"id":         player.ID,
				"username":   player.Username,
				"name":       player.PlayerName,
				"slug":       player.Slug,
				"short_name": player.ShortName,
				"country":    player.Country,
				"position":   player.Positions,
				"media_url":  player.MediaUrl,
			},
			"id":            incident.ID,
			"incident_type": incident.IncidentType,
			"incident_time": incident.IncidentTime,
			"description":   incident.Description,
			"created_at":    incident.CreatedAt,
		}

		matchIncidents = append(matchIncidents, matchIncident)
	}
	s.logger.Info("Successfully get match incidents")
	ctx.JSON(http.StatusAccepted, matchIncidents)
}
