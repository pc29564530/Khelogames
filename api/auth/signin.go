package auth

import (
	"khelogames/core/token"
	"khelogames/database/models"
	"khelogames/util"
	"net/http"

	errorhandler "khelogames/error_handler"

	"github.com/gin-gonic/gin"
)

func (s *AuthServer) CreateEmailSignInFunc(ctx *gin.Context) {
	userAgent := ctx.Request.UserAgent()
	clientIP := ctx.ClientIP()

	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	// Get user by email
	existingUser, err := s.store.GetUsersByGmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("Database error while fetching user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user by gmail",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Security best practice: Use same error message for both invalid email and password
	// This prevents attackers from enumerating valid email addresses
	if existingUser == nil {
		s.logger.Info("Sign in attempt with non-existent email: ", req.Email)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Invalid email or password",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Verify password
	err = util.CheckPassword(req.Password, *existingUser.HashPassword)
	if err != nil {
		s.logger.Error("Failed password attempt for email: ", req.Email)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "AUTHENTICATION_ERROR",
				"message": "Invalid email or password",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Create tokens
	tokens, err := token.CreateNewToken(
		ctx,
		s.store,
		s.tokenMaker,
		int32(existingUser.ID),
		existingUser.PublicID,
		s.config.AccessTokenDuration,
		s.config.RefreshTokenDuration,
		userAgent,
		clientIP,
	)
	if err != nil {
		s.logger.Error("Token creation failed",
			"user_public_id", existingUser.PublicID,
			"request_id", ctx.GetString("request_id"),
			"client_ip", clientIP,
			"user_agent", userAgent,
			"error", err,
		)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "AUTH_SERVICE_UNAVAILABLE",
				"message": "Unable to sign in right now. Please try again later.",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	session := tokens["session"].(*models.Session)
	accessToken := tokens["accessToken"].(string)
	accessPayload := tokens["accessPayload"].(*token.Payload)
	refreshToken := tokens["refreshToken"].(string)
	refreshPayload := tokens["refreshPayload"].(*token.Payload)

	s.logger.Info("Successful sign in for user: ", existingUser.PublicID)

	ctx.JSON(http.StatusOK, gin.H{
		"success":               true,
		"sessionID":             session.ID,
		"accessToken":           accessToken,
		"accessTokenExpiresAt":  accessPayload.ExpiredAt,
		"refreshToken":          refreshToken,
		"refreshTokenExpiresAt": refreshPayload.ExpiredAt,
		"user":                  existingUser,
	})
}
