package football

import (
	"encoding/json"
	"fmt"
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type addFootballIncidentsRequest struct {
	MatchPublicID        string  `json:"match_public_id"`
	TeamPublicID         *string `json:"team_public_id"`
	TournamentPublicID   string  `json:"tournament_public_id"`
	PlayerPublicID       string  `json:"player_public_id"`
	Periods              string  `json:"periods"`
	IncidentType         string  `json:"incident_type"`
	IncidentTime         int     `json:"incident_time"`
	Description          string  `json:"description"`
	PenaltShootoutScored bool    `json:"penalty_shootout_scored"`
}

func (s *FootballServer) AddFootballIncidentsFunc(ctx *gin.Context) {
	var req addFootballIncidentsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var teamPublicID uuid.UUID

	if req.TeamPublicID != nil && *req.TeamPublicID != "" {
		parsed, err := uuid.Parse(*req.TeamPublicID)
		if err != nil {
			s.logger.Error("Invalid UUID format", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TeamPublicID UUID"})
			return
		}
		teamPublicID = parsed
	}

	var playerPublicID uuid.UUID

	if req.IncidentType != "periods" {
		playerPublicID, err = uuid.Parse(req.PlayerPublicID)
		if err != nil {
			s.logger.Error("Invalid UUID format", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		}
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Tournamet not found"})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Check  failed"})
		return
	}
	if !isExists {
		ctx.JSON(403, gin.H{"error": "You do not own this match"})
		return
	}

	fmt.Println("Player PUblic ID: ", playerPublicID)

	arg := db.CreateFootballIncidentsParams{
		TournamentPublicID:    tournamentPublicID,
		MatchPublicID:         matchPublicID,
		TeamPublicID:          &teamPublicID,
		Periods:               req.Periods,
		IncidentType:          req.IncidentType,
		IncidentTime:          req.IncidentTime,
		Description:           req.Description,
		PenaltyShootoutScored: req.PenaltShootoutScored,
	}
	s.logger.Debugf("Creating incident with params: %+v", arg)

	incidentData, err := s.txStore.AddFootballIncidentsTx(ctx, arg, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to create footbal incidents: ", err)
		return
	}

	s.logger.Info("successfully update the add football incident:  , incidentData")

	ctx.JSON(http.StatusAccepted, incidentData)
	if s.scoreBroadcaster != nil {
		err := s.scoreBroadcaster.BroadcastFootballEvent(ctx, "ADD_FOOTBALL_INCIDENT", incidentData)
		if err != nil {
			s.logger.Error("Failed to broadcast cricket event: ", err)
		}
	}
}

type addFootballIncidentsSubsRequest struct {
	MatchPublicID     string `json:"match_public_id"`
	TeamPublicID      string `json:"team_public_id"`
	Periods           string `json:"periods"`
	IncidentType      string `json:"incident_type"`
	IncidentTime      int    `json:"incident_time"`
	Description       string `json:"description"`
	PlayerInPublicID  string `json:"player_in_public_id"`
	PlayerOutPublicID string `json:"player_out_public_in"`
}

func (s *FootballServer) AddFootballIncidentsSubs(ctx *gin.Context) {
	var req addFootballIncidentsSubsRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	playerInPublicID, err := uuid.Parse(req.PlayerInPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	playerOutPublicID, err := uuid.Parse(req.PlayerOutPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Tournamet not found"})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Check  failed"})
		return
	}
	if !isExists {
		ctx.JSON(403, gin.H{"error": "You do not own this match"})
		return
	}

	incidentData, err := s.txStore.AddFootballIncidentsSubsTx(ctx,
		matchPublicID,
		teamPublicID,
		req.Periods,
		req.IncidentType,
		req.IncidentTime,
		req.Description,
		playerInPublicID,
		playerOutPublicID,
	)
	if err != nil {
		s.logger.Error("Failed to create football incidents transaction: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, incidentData)
	if s.scoreBroadcaster != nil {
		err := s.scoreBroadcaster.BroadcastFootballEvent(ctx, "ADD_FOOTBALL_SUB_INCIDENT", *incidentData)
		if err != nil {
			s.logger.Error("Failed to broadcast cricket event: ", err)
		}
	}
}

type getFootballIncidentsRequest struct {
	MatchPublicID string `uri:"match_public_id" form:"match_public_id"`
}

func (s *FootballServer) GetFootballIncidentsFunc(ctx *gin.Context) {
	var req getFootballIncidentsRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	s.logger.Debug("Successfully bind the req: ", req)

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	response, err := s.store.GetFootballIncidentWithPlayer(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get football incidents: ", err)
		return
	}

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match data: ", err)
		return
	}

	var incidents []map[string]interface{}

	// Initialize score counters outside the loop
	homeGoals := 0
	awayGoals := 0
	homeTeamID := match.HomeTeamID
	awayTeamID := match.AwayTeamID

	for _, incident := range response {

		if incident.IncidentType == "substitutions" {
			var data map[string]interface{}
			tt := (incident.Players).([]byte)
			err := json.Unmarshal(tt, &data)
			if err != nil {
				s.logger.Error("unable to unmarshal incident player: ", err)
				continue
			}

			playerInData := data["player_in"].(map[string]interface{})
			playerOutData := data["player_out"].(map[string]interface{})
			incidentDataMap := map[string]interface{}{
				"id":            incident.ID,
				"public_id":     incident.PublicID,
				"match_id":      incident.MatchID,
				"team_id":       incident.TeamID,
				"periods":       incident.Periods,
				"incident_type": incident.IncidentType,
				"incident_time": incident.IncidentTime,
				"description":   incident.Description,
				"player_in": map[string]interface{}{
					"id":         playerInData["id"],
					"public_id":  playerInData["public_id"],
					"user_id":    playerInData["user_id"],
					"name":       playerInData["name"],
					"slug":       playerInData["slug"],
					"short_name": playerInData["short_name"],
					"positions":  playerInData["positions"],
					"country":    playerInData["country"],
					"media_url":  playerInData["media_url"],
				},
				"player_out": map[string]interface{}{
					"id":         playerOutData["id"],
					"public_id":  playerOutData["public_id"],
					"user_id":    playerOutData["user_id"],
					"name":       playerOutData["name"],
					"slug":       playerOutData["slug"],
					"short_name": playerOutData["short_name"],
					"positions":  playerOutData["positions"],
					"country":    playerOutData["country"],
					"media_url":  playerOutData["media_url"],
				},
			}
			incidents = append(incidents, incidentDataMap)

		} else if incident.IncidentType == "penalty_shootout" {
			var data map[string]interface{}
			tt := (incident.Players).([]byte)
			err := json.Unmarshal(tt, &data)
			if err != nil {
				s.logger.Error("unable to unmarshal incident player: ", err)
				continue
			}

			playerData := data["player"].(map[string]interface{})
			incidentDataMap := map[string]interface{}{
				"id":                      incident.ID,
				"public_id":               incident.PublicID,
				"match_id":                incident.MatchID,
				"team_id":                 incident.TeamID,
				"periods":                 incident.Periods,
				"incident_type":           incident.IncidentType,
				"description":             incident.Description,
				"penalty_shootout_scored": incident.PenaltyShootoutScored,
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

			homefootballScore, err := s.store.GetFootballShootoutScoreByTeam(ctx, incident.PublicID, matchPublicID, int32(homeTeamID))
			if err != nil {
				s.logger.Error("unable to fetch the home score: ", err)
			} else {
				incidentDataMap["home_score"] = map[string]interface{}{
					"goals": homefootballScore[0],
				}
			}

			awayfootballScore, err := s.store.GetFootballShootoutScoreByTeam(ctx, incident.PublicID, matchPublicID, int32(awayTeamID))
			if err != nil {
				s.logger.Error("unable to fetch the away score: ", err)
			} else {
				incidentDataMap["away_score"] = map[string]interface{}{
					"goals": awayfootballScore[0],
				}
			}

			incidents = append(incidents, incidentDataMap)

		} else if incident.IncidentType == "period" {
			incidentDataMap := map[string]interface{}{
				"id":            incident.ID,
				"public_id":     incident.PublicID,
				"match_id":      incident.MatchID,
				"periods":       incident.Periods,
				"incident_type": incident.IncidentType,
				"incident_time": incident.IncidentTime,
				"description":   incident.Description,
			}
			incidents = append(incidents, incidentDataMap)

		} else if incident.IncidentType == "goal" {
			var data map[string]interface{}
			tt := (incident.Players).([]byte)
			err := json.Unmarshal(tt, &data)
			if err != nil {
				s.logger.Error("unable to unmarshal incident player: ", err)
				continue
			}

			// Update score counters BEFORE creating the incident map
			if incident.TeamID != nil {
				if homeTeamID == *incident.TeamID {
					homeGoals++
				} else if awayTeamID == *incident.TeamID {
					awayGoals++
				}
			}

			playerData := data["player"].(map[string]interface{})
			incidentDataMap := map[string]interface{}{
				"id":            incident.ID,
				"public_id":     incident.PublicID,
				"match_id":      incident.MatchID,
				"team_id":       incident.TeamID,
				"periods":       incident.Periods,
				"incident_type": incident.IncidentType,
				"incident_time": incident.IncidentTime,
				"description":   incident.Description,
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
				"home_score": map[string]interface{}{
					"goals": homeGoals,
				},
				"away_score": map[string]interface{}{
					"goals": awayGoals,
				},
			}

			incidents = append(incidents, incidentDataMap)

		} else {
			// Handle other incident types (cards, etc.)
			var data map[string]interface{}
			tt := (incident.Players).([]byte)
			err := json.Unmarshal(tt, &data)
			if err != nil {
				s.logger.Error("unable to unmarshal incident player: ", err)
				continue
			}

			playerData := data["player"].(map[string]interface{})
			incidentDataMap := map[string]interface{}{
				"id":            incident.ID,
				"public_id":     incident.PublicID,
				"match_id":      incident.MatchID,
				"team_id":       incident.TeamID,
				"periods":       incident.Periods,
				"incident_type": incident.IncidentType,
				"incident_time": incident.IncidentTime,
				"description":   incident.Description,
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
			incidents = append(incidents, incidentDataMap)
		}
	}

	matchDetail := map[string]interface{}{
		"id":              match.ID,
		"public_id":       match.PublicID,
		"tournament_id":   match.TournamentID,
		"home_team_id":    match.HomeTeamID,
		"away_team_id":    match.AwayTeamID,
		"status_code":     match.StatusCode,
		"start_timestamp": match.StartTimestamp,
		"end_timestamp":   match.EndTimestamp,
		"type":            match.Type,
	}

	matchIncidents := []map[string]interface{}{
		{"match": matchDetail},
		{"incidents": incidents},
	}

	s.logger.Info("Successfully get match incidents")
	ctx.JSON(http.StatusOK, matchIncidents)
}
