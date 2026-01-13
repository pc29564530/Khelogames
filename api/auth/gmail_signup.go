package auth

import (
	"khelogames/core/token"
	"khelogames/database/models"
	utils "khelogames/util"
	"net/http"

	errorhandler "khelogames/error_handler"

	"github.com/gin-gonic/gin"
)

func (s *AuthServer) CreateEmailSignUpFunc(ctx *gin.Context) {
	var req struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	// Check if user already exists with this email
	existingUser, err := s.store.GetUsersByGmail(ctx, req.Email)
	if err != nil {
		s.logger.Error("Database error while checking existing user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get existing users.",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	if existingUser != nil {
		s.logger.Info("User already exists with email: ", req.Email)
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "EMAIL_ALREADY_REGISTERED",
				"message": "Email already registered. Please sign in instead.",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Hash the password
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create account. Please try again.",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Generate username
	username := GenerateUsername(req.Email)

	// Create user account
	_, userSignUp, tokens, err := s.txStore.CreateEmailSignUpTx(ctx, s.config, s.store,
		req.FullName,
		username,
		req.Email,
		hashPassword)

	if err != nil {
		s.logger.Error("Failed to create new account: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "SIGNUP_ERROR",
				"message": "Failed to create account. Please try again.",
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

	s.logger.Info("Successfully created account for user: ", userSignUp.PublicID)

	ctx.JSON(http.StatusCreated, gin.H{
		"success":               true,
		"user":                  userSignUp,
		"sessionID":             session.ID,
		"accessToken":           accessToken,
		"accessTokenExpiresAt":  accessPayload.ExpiredAt,
		"refreshToken":          refreshToken,
		"refreshTokenExpiresAt": refreshPayload.ExpiredAt,
		"message":               "Account created successfully! Please check your email to verify your account.",
	})
}
