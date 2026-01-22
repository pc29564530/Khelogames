package tournaments

import (
	"fmt"
	"khelogames/api/orchestrator"
	"khelogames/api/shared"
	"khelogames/core/token"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"khelogames/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *TournamentServer) GetTournamentMatch(ctx *gin.Context) {
	s.logger.Info("Received request to get tournament matches")

	var req struct {
		TournamentPublicID string `uri:"tournament_public_id"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"tournament_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	sports := strings.TrimSpace(ctx.Param("sport"))

	fieldErrors := make(map[string]string)
	if sports == "" {
		fieldErrors["sport"] = "Sport parameter is required"
	}

	if len(fieldErrors) > 0 {
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	s.logger.Debug(fmt.Sprintf("parse the tournament: %v and sports: %v", tournamentPublicID, sports))

	matches, err := s.store.GetMatchByTournamentPublicID(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament match: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get match by tournament public id",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	var scoreBraodcaster shared.ScoreBroadcaster
	checkSportServer := orchestrator.NewCheckSport(s.store, s.logger, scoreBraodcaster)
	matchDetailsWithScore := checkSportServer.CheckSport(sports, matches, tournamentPublicID)

	s.logger.Info("Successfully retrieved tournament match: ", matchDetailsWithScore)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    matchDetailsWithScore,
	})
}

type createTournamentMatchRequest struct {
	TournamentPublicID string  `json:"tournament_public_id" binding:"required"`
	AwayTeamPublicID   string  `json:"away_team_public_id" binding:"required"`
	HomeTeamPublicID   string  `json:"home_team_public_id" binding:"required"`
	StartTimestamp     string  `json:"start_timestamp" binding:"required"`
	EndTimestamp       string  `json:"end_timestamp" binding:"omitempty"`
	Type               string  `json:"type" binding:"required,min=2,max=50"`
	StatusCode         string  `json:"status_code" binding:"required,oneof=not_started"`
	Result             *int64  `json:"result" binding:"omitempty"`
	Stage              string  `json:"stage" binding:"required,oneof=group knockout league"`
	KnockoutLevelID    *int32  `json:"knockout_level_id" binding:"omitempty,min=1"`
	MatchFormat        *string `json:"match_format" binding:"required,min=2,max=50"`
	DayNumber          *int    `json:"day_number" binding:"omitempty,min=1"`
	SubStatus          *string `json:"sub_status" binding:"omitempty,min=2,max=50"`
	Latitude           string  `json:"latitude" binding:"required"`
	Longitude          string  `json:"longitude" binding:"required"`
	City               string  `json:"city" binding:"omitempty,min=2,max=100"`
	State              string  `json:"state" binding:"omitempty,min=2,max=100"`
	Country            string  `json:"country" binding:"omitempty,min=2,max=100"`
}

func (s *TournamentServer) CreateTournamentMatch(ctx *gin.Context) {
	s.logger.Info("Received request to create tournament match")

	var req createTournamentMatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	s.logger.Debug("bind the request: ", req)

	fieldErrors := make(map[string]string)

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to parse tournament_public_id: ", err)
		fieldErrors["tournament_public_id"] = "Invalid UUID format"
	}

	homeTeamPublicID, err := uuid.Parse(req.HomeTeamPublicID)
	if err != nil {
		s.logger.Error("Failed to parse home_team_public_id: ", err)
		fieldErrors["home_team_public_id"] = "Invalid UUID format"
	}

	awayTeamPublicID, err := uuid.Parse(req.AwayTeamPublicID)
	if err != nil {
		s.logger.Error("Failed to parse away_team_public_id: ", err)
		fieldErrors["away_team_public_id"] = "Invalid UUID format"
	}

	if len(fieldErrors) > 0 {
		fmt.Println("Field Error: ", fieldErrors)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	gameName := ctx.Param("sport")

	if gameName == "" {
		fieldErrors := map[string]string{"sport": "Sport parameter is required"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	game, err := s.store.GetGamebyName(ctx, gameName)
	if err != nil {
		s.logger.Error("Failed to get game: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get game",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	fieldErrors = make(map[string]string)

	startTimeStamp, err := util.ConvertTimeStamp(req.StartTimestamp)
	if err != nil {
		s.logger.Error("unable to convert time to second: ", err)
		fieldErrors["start_timestamp"] = "Invalid timestamp format"
	}

	var endTimeStamp int64
	if req.EndTimestamp != "" {
		endTimeStamp, err = util.ConvertTimeStamp(req.EndTimestamp)
		if err != nil {
			s.logger.Error("unable to convert time to second: ", err)
			fieldErrors["end_timestamp"] = "Invalid timestamp format"
		}
	}

	var matchFormat string
	if gameName == "cricket" {
		if req.MatchFormat != nil {
			matchFormat = *req.MatchFormat
		} else {
			s.logger.Error("Match format is required for cricket matches")
			fieldErrors["match_format"] = "Match format is required for cricket matches"
		}
	}

	var latitude float64
	var longitude float64
	if req.Latitude != "" {
		latitude, err = strconv.ParseFloat(req.Latitude, 64)
		if err != nil {
			s.logger.Error("Failed to parse latitude: ", err)
			fieldErrors["latitude"] = "Invalid format"
		}
	}

	if req.Longitude != "" {
		longitude, err = strconv.ParseFloat(req.Longitude, 64)
		if err != nil {
			s.logger.Error("Failed to parse longitude: ", err)
			fieldErrors["longitude"] = "Invalid format"
		}
	}

	if len(fieldErrors) > 0 {
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	tournament, err := s.store.GetTournament(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get tournament",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(tournament.ID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	if !isExists {
		s.logger.Error("No user role exist for tournament")
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "User does not have permission to create matches for this tournament",
			},
			"request_id": ctx.GetString("request_id"),
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
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create match",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("Successfully created match: ", match)

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    match,
	})
}

type updateMatchSubStatusRequest struct {
	SubStatus string `json:"sub_status" binding: "required"`
}

func (s *TournamentServer) UpdateMatchSubStatusFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update match sub-status")

	var reqUri struct {
		MatchPublicID string `uri:"match_public_id"`
	}
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	var req updateMatchSubStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(reqUri.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get tournament match",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	if !isExists {
		s.logger.Error("User does not have permission for this tournament")
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "User does not have permission to update this match",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	updatedMatchData, err := s.store.UpdateMatchSubStatus(ctx, matchPublicID, req.SubStatus)
	if err != nil {
		s.logger.Error("Failed to update match sub status: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update match sub status",
			},
			"request_id": ctx.GetString("request_id"),
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

	s.logger.Info("Successfully updated match sub-status")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updateMatch,
	})
}

type updateStatusRequest struct {
	StatusCode string `json:"status_code"`
}

func (s *TournamentServer) UpdateMatchStatusFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update match status")

	var reqUri struct {
		MatchPublicID string `uri:"match_public_id"`
	}
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	var req updateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(reqUri.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	game := ctx.Param("sport")

	if game == "" {
		fieldErrors := map[string]string{"sport": "Sport parameter is required"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	gameID, err := s.store.GetGamebyName(ctx, game)
	if err != nil {
		s.logger.Error("Failed to get game: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get game",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament by match id: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get match",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user tournament role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	if !isExists {
		s.logger.Error("User does not have permission for this tournament")
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "User does not have permission to update this match",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	updatedMatchData, err := s.txStore.UpdateMatchStatusTx(ctx, matchPublicID, req.StatusCode, gameID)
	if err != nil {
		s.logger.Error("Failed to update match status: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update match status",
			},
			"request_id": ctx.GetString("request_id"),
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

	s.logger.Info("Successfully updated match status")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updateMatch,
	})
}

type updateMatchResultRequest struct {
	MatchPublicID string `json:"match_public_id" binding: "required`
	Result        string `json:"result" binding: "required"`
}

func (s *TournamentServer) UpdateMatchResultFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update match result")

	var req updateMatchResultRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	fieldErrors := make(map[string]string)

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format for match_public_id: ", err)
		fieldErrors["match_public_id"] = "Invalid UUID format"
	}

	resultPublicID, err := uuid.Parse(req.Result)
	if err != nil {
		s.logger.Error("Invalid UUID format for result: ", err)
		fieldErrors["result"] = "Invalid UUID format"
	}

	if len(fieldErrors) > 0 {
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	match, err := s.store.GetMatchModelByPublicId(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get match",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	team, err := s.store.GetTeamByPublicID(ctx, resultPublicID)
	if err != nil {
		s.logger.Error("Failed to get team: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get team",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	isExists, err := s.store.GetTournamentUserRole(ctx, int32(match.TournamentID), authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get user tournament role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	if !isExists {
		s.logger.Error("User does not have permission for this tournament")
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "User does not have permission to update match result",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	response, err := s.store.UpdateMatchResult(ctx, int32(match.ID), int32(team.ID))
	if err != nil {
		s.logger.Error("Failed to update result: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update match result",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("Successfully updated match result")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
