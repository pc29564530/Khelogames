package auth

import (
	"khelogames/core/token"
	"khelogames/database/models"
	"khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

func (s *AuthServer) CreateEmailSignInFunc(ctx *gin.Context) {
	userAgent := ctx.Request.UserAgent()
	clientIP := ctx.ClientIP()

	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind the login request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	// Get user by email
	existingUser, err := s.store.GetUsersByGmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("Database error while fetching user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "An error occurred. Please try again later.",
		})
		return
	}

	// Security best practice: Use same error message for both invalid email and password
	// This prevents attackers from enumerating valid email addresses
	if existingUser == nil {
		s.logger.Info("Sign in attempt with non-existent email: ", req.Email)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "NOT_FOUND",
			"message": "Invalid email or password",
		})
		return
	}

	// Verify password
	err = util.CheckPassword(req.Password, *existingUser.HashPassword)
	if err != nil {
		s.logger.Error("Failed password attempt for email: ", req.Email)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTHENTICATION_ERROR",
			"message": "Invalid email or password",
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
		s.logger.Error("Failed to create new token: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "An error occurred. Please try again later.",
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
		s.logger.Error("Failed to bind the login request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	// Validate Google ID token
	idToken, err := idtoken.Validate(ctx, req.Code, googleOauthConfig.ClientID)
	if err != nil {
		s.logger.Error("Failed to verify idToken: ", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTHENTICATION_ERROR",
			"message": "Invalid Google token",
		})
		return
	}

	// Extract user info from the verified token
	email, ok := idToken.Claims["email"].(string)
	if !ok {
		s.logger.Error("Failed to get email from idToken")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Failed to get user info from Google",
		})
		return
	}

	// Check if user exists
	existingUser, err := s.store.GetUsersByGmail(ctx, email)
	if err != nil {
		s.logger.Error("Database error while fetching user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "An error occurred. Please try again later.",
		})
		return
	}

	if existingUser == nil {
		s.logger.Info("User does not exist with email: ", email)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"code":    "NOT_FOUND",
			"message": "Email not registered. Please sign up instead.",
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
		s.logger.Error("Failed to create new token: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "An error occurred. Please try again later.",
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
