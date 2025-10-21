package football

import (
	"encoding/json"
	db "khelogames/database"
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
	IncidentTime         int64   `json:"incident_time"`
	Description          string  `json:"description"`
	PenaltShootoutScored bool    `json:"penalty_shootout_scored"`
}

func (s *FootballServer) AddFootballIncidents(ctx *gin.Context) {
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

	teamPublicID, err := uuid.Parse(*req.TeamPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	playerPublicID, err := uuid.Parse(req.PlayerPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
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

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("failed to start transcation: ", err)
		return
	}

	defer tx.Rollback()

	incidents, err := s.store.CreateFootballIncidents(ctx, arg)
	if err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create football incidents: ", err)
		return
	}

	s.logger.Info("successfully created the incident: ", incidents)

	var playerData map[string]interface{}
	var incidentData map[string]interface{}

	if incidents.IncidentType != "period" {

		data, err := s.store.AddFootballIncidentPlayer(ctx, incidents.PublicID, playerPublicID)
		if err != nil {
			tx.Rollback()
			s.logger.Error("Failed to create football incidents: ", err)
			return
		}

		playerData = *data

		statsUpdate := GetStatisticsUpdateFromIncident(incidents.IncidentType)

		statsArg := db.UpdateFootballStatisticsParams{
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

		_, err = s.store.UpdateFootballStatistics(ctx, statsArg)
		if err != nil {
			tx.Rollback()
			s.logger.Error("Failed to update statistics: ", err)
			return
		}

		//Handle goals, penalty_shootout and penalty
		switch incidents.IncidentType {
		case "goal", "penalty":
			if incidents.Periods == "first_half" {
				argGoalScore := db.UpdateFirstHalfScoreParams{
					FirstHalf: 1,
					MatchID:   incidents.MatchID,
					TeamID:    *incidents.TeamID,
				}

				_, err := s.store.UpdateFirstHalfScore(ctx, argGoalScore)
				if err != nil {
					tx.Rollback()
					s.logger.Error("Failed to update football score: ", err)
					return
				}
			} else if incidents.Periods == "second_half" {
				argGoalScore := db.UpdateSecondHalfScoreParams{
					SecondHalf: 1,
					MatchID:    incidents.MatchID,
					TeamID:     *incidents.TeamID,
				}

				_, err := s.store.UpdateSecondHalfScore(ctx, argGoalScore)
				if err != nil {
					tx.Rollback()
					s.logger.Error("Failed to update football score: ", err)
					return
				}
			}
		case "penalty_shootout":
			if incidents.PenaltyShootoutScored {
				_, err := s.store.UpdatePenaltyShootoutScore(ctx, incidents.MatchID, *incidents.TeamID)
				if err != nil {
					tx.Rollback()
					s.logger.Error("Failed to update penalty shootout score: ", err)
					return
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

	//commit the transcation if all operation are successfull
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		s.logger.Error("unable to commit the transcation: ", err)
		return
	}

	s.logger.Info("successfully update the add football incident ")

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
	IncidentTime      int64  `json:"incident_time"`
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

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transcation: ", err)
		return
	}

	defer tx.Rollback()

	arg := db.CreateFootballIncidentsParams{
		MatchPublicID: matchPublicID,
		TeamPublicID:  &teamPublicID,
		Periods:       req.Periods,
		IncidentType:  req.IncidentType,
		IncidentTime:  req.IncidentTime,
		Description:   req.Description,
	}

	incidents, err := s.store.CreateFootballIncidents(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create football incidents: ", err)
		return
	}

	data, err := s.store.ADDFootballSubsPlayer(ctx, incidents.PublicID, playerInPublicID, playerOutPublicID)
	if err != nil {
		s.logger.Error("Failed to create football incidents: ", err)
		return
	}

	subsData := *data

	incidentData := map[string]interface{}{
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

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, incidents)
	if s.scoreBroadcaster != nil {
		err := s.scoreBroadcaster.BroadcastFootballEvent(ctx, "ADD_FOOTBALL_SUB_INCIDENT", incidentData)
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

	for _, incident := range response {

		var awayScore map[string]interface{}
		var homeScore map[string]interface{}

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
					"public_id":  playerInData["public_id"],
					"user_id":    playerInData["user_id"],
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
			if incident.IncidentType == "penalty_shootout" {
				homefootballScore, err := s.store.GetFootballShootoutScoreByTeam(ctx, incident.PublicID, matchPublicID, int32(match.HomeTeamID))
				if err != nil {
					s.logger.Error("unable to fetch the home score: ", err)
				}
				homeScore = map[string]interface{}{
					"goals": homefootballScore[0],
				}
				awayfootballScore, err := s.store.GetFootballShootoutScoreByTeam(ctx, incident.PublicID, matchPublicID, int32(match.AwayTeamID))
				if err != nil {
					s.logger.Error("unable to fetch the home score: ", err)
				}
				awayScore = map[string]interface{}{
					"goals": awayfootballScore[0],
				}

				incidentDataMap["home_score"] = homeScore
				incidentDataMap["away_score"] = awayScore

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
					"name":       playerData["name"],
					"slug":       playerData["slug"],
					"short_name": playerData["short_name"],
					"positions":  playerData["positions"],
					"country":    playerData["country"],
					"media_url":  playerData["media_url"],
				},
			}
			if incident.IncidentType == "goal" {
				homefootballScore, err := s.store.GetFootballScoreByIncidentTime(ctx, int32(incident.ID), incident.MatchID, int32(match.HomeTeamID))
				if err != nil {
					s.logger.Error("unable to fetch the home score: ", err)
				}
				homeScore = map[string]interface{}{
					"goals": homefootballScore[0],
				}
				awayfootballScore, err := s.store.GetFootballScoreByIncidentTime(ctx, int32(incident.ID), incident.MatchID, int32(match.HomeTeamID))
				if err != nil {
					s.logger.Error("unable to fetch the home score: ", err)
				}
				awayScore = map[string]interface{}{
					"goals": awayfootballScore[0],
				}
				incidentDataMap["home_score"] = homeScore
				incidentDataMap["away_score"] = awayScore

			}
			incidents = append(incidents, incidentDataMap)
		}
	}

	var matchIncidents []map[string]interface{}

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
	matchIncidents = append(matchIncidents, map[string]interface{}{
		"match": matchDetail,
	})

	matchIncidents = append(matchIncidents, map[string]interface{}{
		"incidents": incidents,
	})

	s.logger.Info("Successfully get match incidents")
	ctx.JSON(http.StatusAccepted, matchIncidents)
}
