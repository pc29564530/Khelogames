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
		return
	}
	err = ctx.ShouldBindJSON(&reqJSON)
	if err != nil {
		s.logger.Error("Failed to bind json: %w", err)
		return
	}

	matchPublicID, err := uuid.Parse(reqUri.MatchPublicID)
	if err != nil {
		s.logger.Error("Failed to parse to uuid: %w", err)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	response, err := s.store.CreateMatchMedia(ctx, authPayload.UserID, matchPublicID, reqJSON.MediaURL, reqJSON.Title, reqJSON.Description)
	if err != nil {
		s.logger.Error("Failed to create match media: %w", err)
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
		return
	}

	matchPublicID, err := uuid.Parse(reqUri.MatchPublicID)
	if err != nil {
		s.logger.Error("Failed to parse to uuid: %w", err)
		return
	}

	response, err := s.store.GetMatchMedia(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match media: %w", err)
		return
	}
	ctx.JSON(http.StatusAccepted, response)
}
