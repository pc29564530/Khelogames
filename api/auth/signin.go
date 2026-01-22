package auth

import (
	"khelogames/core/token"
	"khelogames/database/models"
	"khelogames/util"
	"net/http"

	errorhandler "khelogames/error_handler"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
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

func (s *AuthServer) CreateGoogleSignIn(ctx *gin.Context) {
	userAgent := ctx.Request.UserAgent()
	clientIP := ctx.ClientIP()
	googleOauthConfig := getGoogleOauthConfig()

	var req getGoogleLoginRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	// Validate Google ID token
	idToken, err := idtoken.Validate(ctx, req.Code, googleOauthConfig.ClientID)
	if err != nil {
		s.logger.Error("Failed to verify idToken: ", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "AUTHENTICATION_ERROR",
				"message": "Invalid google token",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Extract user info from the verified token
	email, ok := idToken.Claims["email"].(string)
	if !ok {
		s.logger.Error("Failed to get email from idToken")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "AUTH_SERVICE_UNAVAILABLE",
				"message": "Unable to sign in right now. Please try again later.",
			},
			"request_id": ctx.GetString("request_id"),
		})
	}

	// Check if user exists
	existingUser, err := s.store.GetUsersByGmail(ctx, email)
	if err != nil {
		s.logger.Error("Database error while fetching user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get users by gmail",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	if existingUser == nil {
		s.logger.Info("User does not exist with email: ", email)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "Email not registered. Please sign up instead.",
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

	s.logger.Info("Successful Google sign in for user: ", existingUser.PublicID)

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
