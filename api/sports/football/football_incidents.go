package football

import (
	"context"
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
		Description:  req.Description,
	}

	incidents, err := s.store.CreateFootballIncidents(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create football incidents: ", err)
		return
	}

	incidentPlayerArg := db.AddFootballIncidentPlayerParams{
		IncidentID: incidents.ID,
		PlayerID:   req.PlayerID,
	}

	_, err = s.store.AddFootballIncidentPlayer(ctx, incidentPlayerArg)
	if err != nil {
		s.logger.Error("Failed to create football incidents: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, incidents)
}

type addFootballIncidentsSubsRequest struct {
	MatchID      int64  `json:"match_id"`
	TeamID       int64  `json:"team_id"`
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
		incidentData := getSwitchFunc(incident, s)
		matchIncidents = append(matchIncidents, incidentData)
	}

	s.logger.Info("Successfully get match incidents")
	ctx.JSON(http.StatusAccepted, matchIncidents)
}

func substitutionIncidentData(incident db.FootballIncident, s *FootballServer) map[string]interface{} {
	ctx := context.Background()
	team, err := s.store.GetTeam(ctx, incident.TeamID)
	if err != nil {
		s.logger.Error("Failed to get teams: ", err)
	}

	subsPlayer, err := s.store.GetFootballIncidentSubsPlayer(ctx, incident.ID)
	if err != nil {
		s.logger.Error("Failed to get player: ", err)
	}

	substitutionIn, err := s.store.GetPlayer(ctx, subsPlayer.PlayerInID)
	if err != nil {
		s.logger.Error("Failed to get player: ", err)
	}

	substitutionOut, err := s.store.GetPlayer(ctx, subsPlayer.PlayerOutID)
	if err != nil {
		s.logger.Error("Failed to get player: ", err)
	}

	matchDetail := map[string]interface{}{
		"id": incident.ID,
		"team": map[string]interface{}{
			"id":        team.ID,
			"name":      team.Name,
			"slug":      team.Slug,
			"shortName": team.Shortname,
			"gender":    team.Gender,
			"national":  team.National,
			"country":   team.Country,
			"type":      team.Type,
		},
		"incident_type": incident.IncidentType,
		"incident_time": incident.IncidentTime,
		"substitution_in_player": map[string]interface{}{
			"id":         substitutionIn.ID,
			"username":   substitutionIn.Username,
			"name":       substitutionIn.PlayerName,
			"slug":       substitutionIn.Slug,
			"short_name": substitutionIn.ShortName,
			"country":    substitutionIn.Country,
			"position":   substitutionIn.Positions,
			"media_url":  substitutionIn.MediaUrl,
		},
		"substitution_out_player": map[string]interface{}{
			"id":         substitutionOut.ID,
			"username":   substitutionOut.Username,
			"name":       substitutionOut.PlayerName,
			"slug":       substitutionOut.Slug,
			"short_name": substitutionOut.ShortName,
			"country":    substitutionOut.Country,
			"position":   substitutionOut.Positions,
			"media_url":  substitutionOut.MediaUrl,
		},
		"created_at": incident.CreatedAt,
	}
	return matchDetail
}

func incidentDataFunc(incident db.FootballIncident, s *FootballServer) map[string]interface{} {
	ctx := context.Background()
	team, err := s.store.GetTeam(ctx, incident.TeamID)
	if err != nil {
		s.logger.Error("Failed to get teams: ", err)
	}

	incidentPlayer, err := s.store.GetFootballIncidentPlayer(ctx, incident.ID)
	if err != nil {
		s.logger.Error("Failed to get player: ", err)
	}

	player, err := s.store.GetPlayer(ctx, incidentPlayer.PlayerID)

	matchDetail := map[string]interface{}{
		"id": incident.ID,
		"team": map[string]interface{}{
			"id":        team.ID,
			"name":      team.Name,
			"slug":      team.Slug,
			"shortName": team.Shortname,
			"gender":    team.Gender,
			"national":  team.National,
			"country":   team.Country,
			"type":      team.Type,
		},
		"incident_type": incident.IncidentType,
		"incident_time": incident.IncidentTime,
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
		"created_at": incident.CreatedAt,
	}
	return matchDetail
}

func getSwitchFunc(incident db.FootballIncident, server *FootballServer) map[string]interface{} {
	switch incident.IncidentType {
	case "substitutions":
		subsIncidentData := substitutionIncidentData(incident, server)
		return subsIncidentData
	default:
		incidentData := incidentDataFunc(incident, server)
		return incidentData
	}
}
