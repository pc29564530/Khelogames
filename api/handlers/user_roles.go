package handlers

import (
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"

	"github.com/gin-gonic/gin"
)

func (s *Server) AddMatchUserRoleFunc(ctx *gin.Context) {
	var req struct {
		MatchPublicID   string `json:"match_public_id"`
		ProfilePublicID string `json:"profile_public_id"`
		Role            string `json:"role"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fieldsErrors := errorhandler.ExtractValidationErrors
		errorhandler.ValidationErrorResponse(err, fieldErrors)
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

	authToken := ctx.MustGet(pkg.AuthorizationPayloadKey).(*tokenMaker.payload)

	userRole, err := s.store.AddMatchUserRole(ctx, matchPublicID, profilePublicID, role, authPayload.UserID)
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

func (s *Server) GetMatchUserRole(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.payload)
	userRole, err := s.store.GetMatchUserRole(ctx, authPayload.UserID)
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
