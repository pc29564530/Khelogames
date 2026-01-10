package tournaments

import (
	"fmt"
	"khelogames/api/orchestrator"
	"khelogames/api/shared"
	"khelogames/core/token"
	"khelogames/pkg"
	"khelogames/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *TournamentServer) GetTournamentMatch(ctx *gin.Context) {

	var req struct {
		TournamentPublicID string `uri:"tournament_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	sports := strings.TrimSpace(ctx.Param("sport"))
	s.logger.Debug(fmt.Sprintf("parse the tournament: %v and sports: %v", tournamentPublicID, sports))
	s.logger.Debug("Tournament match params: ", req.TournamentPublicID)

	matches, err := s.store.GetMatchByTournamentPublicID(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament match: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get match by tournament public id",
		})
		return
	}
	var scoreBraodcaster shared.ScoreBroadcaster
	checkSportServer := orchestrator.NewCheckSport(s.store, s.logger, scoreBraodcaster)
	matchDetailsWithScore := checkSportServer.CheckSport(sports, matches, tournamentPublicID)

	s.logger.Info("successfully  get the tournament match: ", matchDetailsWithScore)
	ctx.JSON(http.StatusAccepted, matchDetailsWithScore)
}

type createTournamentMatchRequest struct {
	TournamentPublicID string  `json:"tournament_public_id"`
	AwayTeamPublicID   string  `json:"away_team_public_id"`
	HomeTeamPublicID   string  `json:"home_team_public_id"`
	StartTimestamp     string  `json:"start_timestamp"`
	EndTimestamp       string  `json:"end_timestamp"`
	Type               string  `json:"type"`
	StatusCode         string  `json:"status_code"`
	Result             *int64  `json:"result"`
	Stage              string  `json:"stage"`
	KnockoutLevelID    *int32  `json:"knockout_level_id"`
	MatchFormat        *string `json:"match_format"`
	DayNumber          *int    `json:"day_number"`
	SubStatus          *string `json:"sub_status"`
	Latitude           string  `json:"latitude"`
	Longitude          string  `json:"longitude"`
	City               string  `json:"city"`
	State              string  `json:"state"`
	Country            string  `json:"country"`
}

func (s *TournamentServer) CreateTournamentMatch(ctx *gin.Context) {

	var req createTournamentMatchRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	s.logger.Debug("bind the request: ", req)

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse tournament_public_id: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}
	homeTeamPublicID, err := uuid.Parse(req.HomeTeamPublicID)
	if err != nil {
		s.logger.Error("Failed to parse home_team_public_id: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}
	awayTeamPublicID, err := uuid.Parse(req.AwayTeamPublicID)
	if err != nil {
		s.logger.Error("Failed to parse away_team_public_id: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	gameName := ctx.Param("sport")

	game, err := s.store.GetGamebyName(ctx, gameName)
	if err != nil {
		s.logger.Error("Failed to get game: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get game",
		})
		return
	}

	startTimeStamp, err := util.ConvertTimeStamp(req.StartTimestamp)
	if err != nil {
		s.logger.Error("unable to convert time to second: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid timestamp format",
		})
		return
	}
	var endTimeStamp int64
	if req.EndTimestamp != "" {
		endTimeStamp, err = util.ConvertTimeStamp(req.EndTimestamp)
		if err != nil {
			s.logger.Error("unable to convert time to second: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    "VALIDATION_ERROR",
				"message": "Invalid timestamp format",
			})
			return
		}
	}
	var matchFormat string
	if gameName == "cricket" {
		if req.MatchFormat != nil {
			matchFormat = *req.MatchFormat
		} else {
			s.logger.Error("Match format is required for cricket matches")
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    "NOT_FOUND",
				"message": "Match format is required",
			})
			return
		}
	}

	var latitude float64
	var longitude float64
	if req.Latitude != "" && req.Longitude != "" {
		latitude, err = strconv.ParseFloat(req.Latitude, 64)
		if err != nil {
			s.logger.Error("Failed to parse to float: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request format",
			})
			return
		}
		longitude, err = strconv.ParseFloat(req.Longitude, 64)
		if err != nil {
			s.logger.Error("Failed to parse to float: ", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request format",
			})
			return
		}
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get tournament",
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(tournament.ID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get user role",
		})
		return
	}
	if !isExists {
		s.logger.Error("No user role exist for tournament: ", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    "NOT_FOUND_ERROR",
			"message": "Failed to find user role",
		})
		return
	}

	match, err := s.txStore.CreateMatchTx(ctx,
		authPayload.UserID,
		latitude,
		longitude,
		req.City,
		req.State,
		req.Country,
		tournamentPublicID,
		awayTeamPublicID,
		homeTeamPublicID,
		startTimeStamp,
		endTimeStamp,
		req.Type,
		req.StatusCode,
		req.Result,
		req.Stage,
		req.KnockoutLevelID,
		&matchFormat,
		req.SubStatus,
		game.ID,
	)

	if err != nil {
		s.logger.Error("Failed to create new match: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get new match",
		})
		return
	}

	s.logger.Debug("Successfully created match: ", match)

	ctx.JSON(http.StatusAccepted, match)
}

type updateMatchSubStatusRequest struct {
	SubStatus string `json:"sub_status"`
}

func (s *TournamentServer) UpdateMatchSubStatusFunc(ctx *gin.Context) {
	// Bind URI
	var reqUri struct {
		MatchPublicID string `uri:"match_public_id"`
	}
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	// Bind JSON
	var req updateMatchSubStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	matchPublicID, err := uuid.Parse(reqUri.MatchPublicID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get tournament match",
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get user role",
		})
		return
	}
	if !isExists {
		s.logger.Error("Failed to find user role: ", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code": "NOT_FOUND_ERROR",
			"message": "Failed to find user role",
		})
	}

	updatedMatchData, err := s.store.UpdateMatchSubStatus(ctx, matchPublicID, req.SubStatus)
	if err != nil {
		s.logger.Error("Failed to update match status: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to update match sub status",
		})
		return
	}

	updateMatch := map[string]interface{}{
		"id":                updatedMatchData.ID,
		"public_id":         updatedMatchData.PublicID,
		"tournament_id":     updatedMatchData.TournamentID,
		"home_team_id":      updatedMatchData.HomeTeamID,
		"away_team_id":      updatedMatchData.AwayTeamID,
		"status_code":       updatedMatchData.StatusCode,
		"sub_status":        updatedMatchData.SubStatus,
		"match_format":      updatedMatchData.MatchFormat,
		"stage":             updatedMatchData.Stage,
		"day_number":        updatedMatchData.DayNumber,
		"type":              updatedMatchData.Type,
		"end_timestamp":     updatedMatchData.EndTimestamp,
		"start_timestamp":   updatedMatchData.StartTimestamp,
		"knockout_level_id": updatedMatchData.KnockoutLevelID,
		"result":            updatedMatchData.Result,
	}

	if s.scoreBroadcaster != nil {
		err := s.scoreBroadcaster.BroadcastTournamentEvent(ctx, "UPDATE_MATCH_SUB_STATUS", updateMatch)
		if err != nil {
			s.logger.Warn("Failed to broadcast tournament match event: ", err)
		}
	}

	s.logger.Info("successfully updated the match status")
	ctx.JSON(http.StatusAccepted, updateMatch)
}

type updateStatusRequest struct {
	StatusCode string `json:"status_code"`
}

func (s *TournamentServer) UpdateMatchStatusFunc(ctx *gin.Context) {
	// Bind URI
	var reqUri struct {
		MatchPublicID string `uri:"match_public_id"`
	}
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	// Bind JSON
	var req updateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	matchPublicID, err := uuid.Parse(reqUri.MatchPublicID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	game := ctx.Param("sport")
	gameID, err := s.store.GetGamebyName(ctx, game)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get game",
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament by match id: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get match",
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user tournament role: ", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    "NOT_FOUND_ERROR",
			"message": "Check failed",
		})
		return
	}
	if !isExists {
		s.logger.Error("Tournament user role does not exist: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get user role",
		})
		return
	}

	updatedMatchData, err := s.txStore.UpdateMatchStatusTx(ctx, matchPublicID, req.StatusCode, gameID)
	if err != nil {
		s.logger.Error("Failed to update match status: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to update match status",
		})
		return
	}

	updateMatch := map[string]interface{}{
		"id":                updatedMatchData.ID,
		"public_id":         updatedMatchData.PublicID,
		"tournament_id":     updatedMatchData.TournamentID,
		"home_team_id":      updatedMatchData.HomeTeamID,
		"away_team_id":      updatedMatchData.AwayTeamID,
		"status_code":       updatedMatchData.StatusCode,
		"sub_status":        updatedMatchData.SubStatus,
		"match_format":      updatedMatchData.MatchFormat,
		"stage":             updatedMatchData.Stage,
		"day_number":        updatedMatchData.DayNumber,
		"type":              updatedMatchData.Type,
		"end_timestamp":     updatedMatchData.EndTimestamp,
		"start_timestamp":   updatedMatchData.StartTimestamp,
		"knockout_level_id": updatedMatchData.KnockoutLevelID,
		"result":            updatedMatchData.Result,
	}

	if s.scoreBroadcaster != nil {
		err := s.scoreBroadcaster.BroadcastTournamentEvent(ctx, "UPDATE_MATCH_STATUS", updateMatch)
		if err != nil {
			s.logger.Warn("Failed to broadcast tournament match event: ", err)
		}
	}

	s.logger.Info("successfully updated the match status")
	ctx.JSON(http.StatusAccepted, updateMatch)
}

type updateMatchResultRequest struct {
	MatchPublicID string `json:"match_public_id"`
	Result        string `json:"result"`
}

func (s *TournamentServer) UpdateMatchResultFunc(ctx *gin.Context) {
	var req updateMatchResultRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code": "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	resultPublicID, err := uuid.Parse(req.Result)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	match, err := s.store.GetMatchModelByPublicId(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get match",
		})
		return
	}
	team, err := s.store.GetTeamByPublicID(ctx, resultPublicID)
	if err != nil {
		s.logger.Error("Failed to team ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATBASE_ERROR",
			"message": "Failed to get team",
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user tournament role: ", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    "NOT_FOUND_ERROR",
			"message": "Check failed",
		})
		return
	}
	if !isExists {
		s.logger.Error("Tournament user role does not exist: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to get user role",
		})
		return
	}

	response, err := s.store.UpdateMatchResult(ctx, int32(match.ID), int32(team.ID))
	if err != nil {
		s.logger.Error("Failed to update result: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code": "DATABASE_ERROR",
			"message": "Failed to get update match result",
		})
		return
	}

	s.logger.Info("Successfully update match result")
	ctx.JSON(http.StatusAccepted, response)
}
