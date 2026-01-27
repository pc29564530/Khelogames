package handlers

import (
	"khelogames/core/token"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *HandlersServer) AddMatchUserRoleFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID   string `json:"match_public_id"`
		ProfilePublicID string `json:"profile_public_id"`
		Role            string `json:"role"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	s.logger.Debug("bind the request: ", req)

	profilePublicID, err := uuid.Parse(req.ProfilePublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"profile_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	// Get match ID from public ID
	match, err := s.store.GetTournamentMatchByMatchID(ctx, matchPublicID)
	if err != nil {
		s.logger.Error("Failed to get match: ", err)
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

	// Get profile ID from public ID
	profile, err := s.store.GetProfileByPublicID(ctx, profilePublicID)
	if err != nil {
		s.logger.Error("Failed to get profile: ", err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Profile not found",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	userRole, err := s.store.AddMatchUserRole(ctx, int32(match.ID), profile.UserID, req.Role, authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to add user role: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to add match user role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    userRole,
	})
	return
}

func (s *HandlersServer) GetMatchUserRoleFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID string `uri:"match_public_id"`
	}

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	matchPublicID, err := uuid.Parse(req.MatchPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"match_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	userRole, err := s.store.GetMatchUserRole(ctx, matchPublicID, authPayload.UserID)
	if err != nil {
		s.logger.Error("Failed to get match user role: %w", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get match user role",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    userRole,
	})
}
