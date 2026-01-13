package handlers

import (
	"khelogames/core/token"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *HandlersServer) CreateMatchMediaFunc(ctx *gin.Context) {
	var reqUri struct {
		MatchPublicID string `uri:"match_public_id" binding:"required"`
	}
	var reqJSON struct {
		Title       string `json:"title" binding:"required,min=1,max=200"`
		Description string `json:"description" binding:"max=1000"`
		MediaURL    string `json:"media_url" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(reqUri.MatchPublicID)
	if err != nil {
		s.logger.Error("Failed to parse to uuid: %w", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	gameName := ctx.Param("sport")

	_, err = s.store.GetGamebyName(ctx, gameName)
	if err != nil {
		s.logger.Error("Failed to get game: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid game name",
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
				"message": "Failed to get tournament match",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, match.TournamentID, authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get tournament by user role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get tournament user role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	if !isExists {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You are not allowed to upload media for this match",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	response, err := s.store.CreateMatchMedia(ctx, authPayload.UserID, matchPublicID, reqJSON.MediaURL, reqJSON.Title, reqJSON.Description)
	if err != nil {
		s.logger.Error("Failed to create match media: %w", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Could not create match media",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *HandlersServer) GetMatchMediaFunc(ctx *gin.Context) {
	var reqUri struct {
		MatchPublicID string `uri:"match_public_id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(reqUri.MatchPublicID)
	if err != nil {
		s.logger.Error("Failed to parse to uuid: %w", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	response, err := s.store.GetMatchMedia(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match media: %w", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Could not get match media",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}
