package football

import (
	"encoding/json"
	db "khelogames/database"
	errorhandler "khelogames/error_handler"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	s.logger.Info("Received request to add football incident")
	var req addFootballIncidentsRequest
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		s.logger.Error("Failed to bind request: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid tournament UUID format: ", err)
		fieldErrors := map[string]string{"tournament_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid match UUID format: ", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	var teamPublicID uuid.UUID

	if req.TeamPublicID != nil && *req.TeamPublicID != "" {
		parsed, err := uuid.Parse(*req.TeamPublicID)
		if err != nil {
			s.logger.Error("Invalid team UUID format: ", err)
			fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
			errorhandler.ValidationErrorResponse(ctx, fieldErrors)
			return
		}
		teamPublicID = parsed
	}

	var playerPublicID uuid.UUID

	if req.IncidentType != "periods" {
		playerPublicID, err = uuid.Parse(req.PlayerPublicID)
		if err != nil {
			s.logger.Error("Invalid player UUID format: ", err)
			fieldErrors := map[string]string{"player_public_id": "Invalid UUID format"}
			errorhandler.ValidationErrorResponse(ctx, fieldErrors)
			return
		}
	}

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

	txResult, err := s.txStore.AddFootballIncidentsTx(ctx, arg, playerPublicID)
	if err != nil {
		s.logger.Error("Failed to create football incidents: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Could not create football incident",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("Successfully added football incident")

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    txResult.IncidentData,
	})

	// Broadcast after successful transaction commit
	if s.scoreBroadcaster != nil {
		// Broadcast the incident event
		if err := s.scoreBroadcaster.BroadcastFootballEvent(ctx, "ADD_FOOTBALL_INCIDENT", txResult.IncidentData); err != nil {
			s.logger.Warn("Failed to broadcast football incident event: ", err)
		}

		// Broadcast score update if a score-changing incident occurred
		if txResult.ScoreData != nil {
			if err := s.scoreBroadcaster.BroadcastFootballEvent(ctx, "UPDATE_FOOTBALL_SCORE", txResult.ScoreData); err != nil {
				s.logger.Warn("Failed to broadcast football score event: ", err)
			}
		}
	}
}

type addFootballIncidentsSubsRequest struct {
	MatchPublicID      string `json:"match_public_id"`
	TeamPublicID       string `json:"team_public_id"`
	TournamentPublicID string `json:"tournament_public_id"`
	Periods            string `json:"periods"`
	IncidentType       string `json:"incident_type"`
	IncidentTime       int    `json:"incident_time"`
	Description        string `json:"description"`
	PlayerInPublicID   string `json:"player_in_public_id"`
	PlayerOutPublicID  string `json:"player_out_public_id"`
}

func (s *FootballServer) AddFootballIncidentsSubs(ctx *gin.Context) {
	s.logger.Info("Received request to add football substitution incident")
	var req addFootballIncidentsSubsRequest
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		s.logger.Error("Failed to bind request: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid match UUID format: ", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	teamPublicID, err := uuid.Parse(req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid team UUID format: ", err)
		fieldErrors := map[string]string{"team_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid tournament UUID format: ", err)
		fieldErrors := map[string]string{"tournament_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	playerInPublicID, err := uuid.Parse(req.PlayerInPublicID)
	if err != nil {
		s.logger.Error("Invalid player in UUID format: ", err)
		fieldErrors := map[string]string{"player_in_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	playerOutPublicID, err := uuid.Parse(req.PlayerOutPublicID)
	if err != nil {
		s.logger.Error("Invalid player out UUID format: ", err)
		fieldErrors := map[string]string{"player_out_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	incidentData, err := s.txStore.AddFootballIncidentsSubsTx(ctx,
		matchPublicID,
		teamPublicID,
		tournamentPublicID,
		req.Periods,
		req.IncidentType,
		req.IncidentTime,
		req.Description,
		playerInPublicID,
		playerOutPublicID,
	)
	if err != nil {
		s.logger.Error("Failed to create football substitution incident: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Could not create football incident substitution",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("Successfully added football substitution incident")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    incidentData,
	})
	if s.scoreBroadcaster != nil {
		err := s.scoreBroadcaster.BroadcastFootballEvent(ctx, "ADD_FOOTBALL_INCIDENT", *incidentData)
		if err != nil {
			s.logger.Warn("Failed to broadcast football event: ", err)
		}
	}
}

type getFootballIncidentsRequest struct {
	MatchPublicID string `uri:"match_public_id" form:"match_public_id"`
}

func (s *FootballServer) GetFootballIncidentsFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get football incidents")
	var req getFootballIncidentsRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind request: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("Successfully bind the req: ", req)

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid match UUID format: ", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	response, err := s.store.GetFootballIncidentWithPlayer(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get football incidents: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get football incidents",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match data: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get match details",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Initialize score counters
	homeGoals := 0
	awayGoals := 0
	homeTeamID := match.HomeTeamID
	awayTeamID := match.AwayTeamID

	var incidents []map[string]interface{}
	for _, incident := range response {
		if incident.IncidentType == "substitution" {
			var data map[string]interface{}
			tt, ok := (incident.Players).([]byte)
			if !ok || tt == nil {
				s.logger.Error("incident.Players is nil or not []byte for substitution incident: ", incident.ID)
				continue
			}
			err := json.Unmarshal(tt, &data)
			if err != nil {
				s.logger.Error("unable to unmarshal incident player: ", err)
				continue
			}

			//

			playerInData, _ := data["player_in"].(map[string]interface{})
			playerOutData, _ := data["player_out"].(map[string]interface{})
			if playerInData == nil || playerOutData == nil {
				s.logger.Error("missing player_in or player_out data for substitution incident: ", incident.ID)
				continue
			}
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
			tt, ok := (incident.Players).([]byte)
			if !ok || tt == nil {
				s.logger.Error("incident.Players is nil or not []byte for penalty_shootout incident: ", incident.ID)
				continue
			}
			err := json.Unmarshal(tt, &data)
			if err != nil {
				s.logger.Error("unable to unmarshal incident player: ", err)
				continue
			}

			playerData, _ := data["player"].(map[string]interface{})
			if playerData == nil {
				s.logger.Error("missing player data for penalty_shootout incident: ", incident.ID)
				continue
			}
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
				incidentDataMap["homeScore"] = map[string]interface{}{
					"goals": homefootballScore[0],
				}
			}

			awayfootballScore, err := s.store.GetFootballShootoutScoreByTeam(ctx, incident.PublicID, matchPublicID, int32(awayTeamID))
			if err != nil {
				s.logger.Error("unable to fetch the away score: ", err)
			} else {
				incidentDataMap["awayScore"] = map[string]interface{}{
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

		} else if incident.IncidentType == "goal" || incident.IncidentType == "penalty" {
			var data map[string]interface{}
			tt, ok := (incident.Players).([]byte)
			if !ok || tt == nil {
				s.logger.Error("incident.Players is nil or not []byte for goal/penalty incident: ", incident.ID)
				continue
			}
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

			playerData, _ := data["player"].(map[string]interface{})
			if playerData == nil {
				s.logger.Error("missing player data for goal/penalty incident: ", incident.ID)
				continue
			}
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
				"homeScore": map[string]interface{}{
					"goals": homeGoals,
				},
				"awayScore": map[string]interface{}{
					"goals": awayGoals,
				},
			}

			incidents = append(incidents, incidentDataMap)

		} else {
			// Handle other incident types (cards, fouls, etc.)
			var data map[string]interface{}
			tt, ok := (incident.Players).([]byte)
			if !ok || tt == nil {
				s.logger.Error("incident.Players is nil or not []byte for incident: ", incident.ID)
				continue
			}
			err := json.Unmarshal(tt, &data)
			if err != nil {
				s.logger.Error("unable to unmarshal incident player: ", err)
				continue
			}

			playerData, _ := data["player"].(map[string]interface{})
			if playerData == nil {
				s.logger.Error("missing player data for incident: ", incident.ID)
				continue
			}
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

	s.logger.Info("Successfully retrieved match incidents")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    matchIncidents,
	})
}
