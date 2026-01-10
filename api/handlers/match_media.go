package handlers

import (
	"khelogames/core/token"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *HandlersServer) CreateMatchMediaFunc(ctx *gin.Context) {
	var reqUri struct {
		MatchPublicID string `uri:"match_public_id"`
	}
	var reqJSON struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		MediaURL    string `json:"media_url"`
	}
	err := ctx.ShouldBindUri(&reqUri)
	if err != nil {
		s.logger.Error("Failed to bind uri: %w", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}
	err = ctx.ShouldBindJSON(&reqJSON)
	if err != nil {
		s.logger.Error("Failed to bind json: %w", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	matchPublicID, err := uuid.Parse(reqUri.MatchPublicID)
	if err != nil {
		s.logger.Error("Failed to parse to uuid: %w", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	gameName := ctx.Param("sport")

	_, err = s.store.GetGamebyName(ctx, gameName)
	if err != nil {
		s.logger.Error("Failed to get game: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid game name",
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get tournament by match id: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get tournament match",
		})
		return
	}

	isExists, err := s.store.GetTournamentUserRole(ctx, match.TournamentID, authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get tournament by user role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get tournament user role",
		})
		return
	}

	if !isExists {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    "FORBIDDEN_ERROR",
			"message": "You are not allowed to upload media for this match",
		})
		return
	}

	response, err := s.store.CreateMatchMedia(ctx, authPayload.UserID, matchPublicID, reqJSON.MediaURL, reqJSON.Title, reqJSON.Description)
	if err != nil {
		s.logger.Error("Failed to create match media: %w", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Could not create match media",
		})
		return
	}
	ctx.JSON(http.StatusAccepted, response)
}

func (s *HandlersServer) GetMatchMediaFunc(ctx *gin.Context) {
	var reqUri struct {
		MatchPublicID string `uri:"match_public_id"`
	}
	err := ctx.ShouldBindUri(&reqUri)
	if err != nil {
		s.logger.Error("Failed to bind uri: %w", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	matchPublicID, err := uuid.Parse(reqUri.MatchPublicID)
	if err != nil {
		s.logger.Error("Failed to parse to uuid: %w", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	response, err := s.store.GetMatchMedia(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match media: %w", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Could not get match media",
		})
		return
	}
	ctx.JSON(http.StatusAccepted, response)
}
