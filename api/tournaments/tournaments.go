package tournaments

import (
	"khelogames/core/token"
	db "khelogames/database"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getTournamentPublicIDRequest struct {
	TournamentPublicID string `uri:"tournament_public_id"`
}

type addTournamentRequest struct {
	Name           string `json:"name" binding:"required,min=3,max=100"`
	Status         string `json:"status" binding:"required,oneof=draft not_started live completed cancelled"`
	Level          string `json:"level" binding:"required,oneof=local state national international"`
	StartTimestamp string `json:"start_timestamp" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`

	GameID        *int64 `json:"game_id" binding:"required,min=1"`
	GroupCount    *int32 `json:"group_count" binding:"omitempty,min=1,max=64"`
	MaxGroupTeams *int32 `json:"max_group_teams" binding:"omitempty,min=2,max=64"`

	Stage       string `json:"stage" binding:"required,oneof=league group knockout"`
	HasKnockout bool   `json:"has_knockout"`

	City    string `json:"city" binding:"required"`
	State   string `json:"state" binding:"required"`
	Country string `json:"country" binding:"required"`
}

//TODO: Normalization the words for location and other input

func (s *TournamentServer) AddTournamentFunc(ctx *gin.Context) {
	var req addTournamentRequest
	fieldErrors := make(map[string]string)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors = errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	if req.Stage == "group" && req.GroupCount == nil {
		fieldErrors["group_count"] = "Required when stage is group"
	}

	if len(fieldErrors) > 0 {
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	startTimestamp, err := util.ConvertTimeStamp(req.StartTimestamp)
	if err != nil {
		s.logger.Warn("invalid_start_timestamp",
			"request_id", ctx.GetString("request_id"),
		)

		errorhandler.ValidationErrorResponse(ctx, map[string]string{
			"start_timestamp": "Invalid timestamp",
		})
		return
	}

	slug := util.GenerateSlug(req.Name)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	tournament, err := s.txStore.AddNewTournamentTx(
		ctx,
		authPayload,
		req.Name,
		slug,
		req.Status,
		req.Level,
		startTimestamp,
		req.GameID,
		req.GroupCount,
		req.MaxGroupTeams,
		req.Stage,
		req.HasKnockout,
		req.City,
		req.State,
		req.Country,
	)

	if err != nil {
		s.logger.Error("failed_to_create_tournament",
			"request_id", ctx.GetString("request_id"),
			"error", err,
		)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create tournament",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    tournament,
	})
}

func (s *TournamentServer) GetTournamentTeamCountFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get team count for tournament")
	var req getTournamentPublicIDRequest
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

	response, err := s.store.GetTournamentTeamsCount(ctx, tournamentPublicID)
	if err != nil {
		s.logger.Error("Failed to get team count: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get team count",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Info("Successfully retrieved team count: %v", response)

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *TournamentServer) GetTournamentsFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get tournaments")

	response, err := s.store.GetTournaments(ctx)
	if err != nil {
		s.logger.Error("Failed to get tournaments: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get tournaments",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Info("Successfully retrieved tournaments: %v", response)

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *TournamentServer) GetTournamentFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get a tournament")
	var req getTournamentPublicIDRequest
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

	response, err := s.store.GetTournament(ctx, tournamentPublicID)
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
	s.logger.Info("Successfully retrieved tournament: %v", response)

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *TournamentServer) UpdateTournamentDateFunc(ctx *gin.Context) {
	s.logger.Info("Received request to update tournament dates")
	var req getTournamentPublicIDRequest
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

	startOnStr := ctx.Query("start_on")
	startTimeStamp, err := util.ConvertTimeStamp(startOnStr)
	if err != nil {
		s.logger.Error("Unable to convert timestamp: ", err)
		fieldErrors := map[string]string{"start_on": "Invalid timestamp format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	arg := db.UpdateTournamentDateParams{
		TournamentPublicID: tournamentPublicID,
		StartTimestamp:     startTimeStamp,
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

	if tournament.UserID != authPayload.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You are not allowed to make this change",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	response, err := s.store.UpdateTournamentDate(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update tournament date: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update tournament date",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Info("Successfully updated tournament date: ", response)

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *TournamentServer) GetTournamentByLevelFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get tournaments by level")

	sports := ctx.Param("sport")
	level := ctx.Query("category")
	s.logger.Debug("Category: %v", level)

	fieldErrors := make(map[string]string)

	if sports == "" {
		fieldErrors["sport"] = "Sport parameter is required"
	}

	if level == "" {
		fieldErrors["category"] = "Category query parameter is required"
	}

	if len(fieldErrors) > 0 {
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	game, err := s.store.GetGamebyName(ctx, sports)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
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

	response, err := s.store.GetTournamentsByLevel(ctx, game.ID, level)
	if err != nil {
		s.logger.Error("Failed to get tournaments by level: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get tournaments by level",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Info("Successfully retrieved tournaments by level: ", response)

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *TournamentServer) UpdateTournamentStatusFunc(ctx *gin.Context) {
	var req getTournamentPublicIDRequest
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

	statusCode := ctx.Query("status_code")

	if statusCode == "" {
		fieldErrors := map[string]string{"status_code": "Status code query parameter is required"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	validStatuses := []string{"draft", "upcoming", "ongoing", "completed", "cancelled"}
	isValid := false
	for _, status := range validStatuses {
		if statusCode == status {
			isValid = true
			break
		}
	}

	if !isValid {
		fieldErrors := map[string]string{"status_code": "Invalid status code. Must be one of: draft, upcoming, ongoing, completed, cancelled"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	arg := db.UpdateTournamentStatusParams{
		TournamentPublicID: tournamentPublicID,
		Status:             statusCode,
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

	if tournament.UserID != authPayload.UserID {
		s.logger.Error("Failed to match tournament user ID with current user")
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You are not allowed to update this tournament",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	updatedTournament, err := s.store.UpdateTournamentStatus(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update tournament status: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update tournament status",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("Successfully updated tournament status")

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    updatedTournament,
	})
}

type getTournamentByGameIdRequest struct {
	GameID int64 `uri:"game_id"`
}

func (s *TournamentServer) GetTournamentsBySportFunc(ctx *gin.Context) {
	var req getTournamentByGameIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	gameName := ctx.Param("sport")

	game, err := s.store.GetGamebyName(ctx, gameName)
	if err != nil {
		s.logger.Error("Failed to get game by name: ", err)
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

	tournaments, err := s.store.GetTournamentsBySport(ctx, req.GameID)
	if err != nil {
		s.logger.Error("Failed to get tournaments by sport: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get tournaments by sport",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	results := map[string]interface{}{
		"id":          game.ID,
		"name":        game.Name,
		"min_players": game.MinPlayers,
		"tournament":  tournaments,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
	})
}

func (s *TournamentServer) AddTournamentUserRolesFunc(ctx *gin.Context) {
	var req struct {
		TournamentPublicID string `json:"tournament_public_id"`
		Role               string `json:"role"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	tournamentPublicID, err := uuid.Parse(req.TournamentPublicID)
	if err != nil {
		fieldErrors := map[string]string{"tournament_public_id": "Invalid UUID format"}
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

	newTournamentUserRoles, err := s.store.AddTournamentUserRoles(ctx, int32(tournament.ID), authPayload.UserID, req.Role)
	if err != nil {
		s.logger.Error("Failed to add tournament user role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to add tournament user role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    newTournamentUserRoles,
	})
}

func (s *TournamentServer) GetTournamentByLocationFunc(ctx *gin.Context) {
	city := ctx.Query("city")
	state := ctx.Query("state")
	country := ctx.Query("country")

	gameName := ctx.Param("sport")

	fieldErrors := make(map[string]string)

	if gameName == "" {
		fieldErrors["sport"] = "Sport parameter is required"
	}

	if city == "" {
		fieldErrors["city"] = "City query parameter is required"
	} else if len(city) < 2 {
		fieldErrors["city"] = "City must be at least 2 characters"
	}

	if state == "" {
		fieldErrors["state"] = "State query parameter is required"
	} else if len(state) < 2 {
		fieldErrors["state"] = "State must be at least 2 characters"
	}

	if country == "" {
		fieldErrors["country"] = "Country query parameter is required"
	} else if len(country) < 2 {
		fieldErrors["country"] = "Country must be at least 2 characters"
	}

	if len(fieldErrors) > 0 {
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

	tournament, err := s.store.GetTournamentByLocation(ctx, game.ID, city, state, country)
	if err != nil {
		s.logger.Error("Failed to get tournaments by location: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get tournaments by location",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	results := map[string]interface{}{
		"id":          game.ID,
		"name":        game.Name,
		"min_players": game.MinPlayers,
		"tournament":  tournament,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
	})
}
