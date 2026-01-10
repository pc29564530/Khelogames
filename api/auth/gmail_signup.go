package auth

import (
	"khelogames/core/token"
	"khelogames/database/models"
	utils "khelogames/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *AuthServer) CreateEmailSignUpFunc(ctx *gin.Context) {
	var req struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind the sign-up request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	// Check if user already exists with this email
	existingUser, err := s.store.GetUsersByGmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		s.logger.Info("User already exists with email: ", req.Email)
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"code":    "EMAIL_ALREADY_REGISTERED",
			"message": "Email already registered. Please sign in instead.",
		})
		return
	}

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to convert to hash: ", err)
		ctx.JSON(http.StatusInternalServerError, "Failed to convert to hash")
		return
	}

	// Generate username
	username := GenerateUsername(req.Email)

	_, userSignUp, tokens, err := s.txStore.CreateEmailSignUpTx(ctx, s.config, s.store,
		req.FullName,
		username,
		req.Email,
		hashPassword)

	if err != nil {
		s.logger.Errorf("Failed to create new account: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "SIGNUP_ERROR",
			"message": "Failed to create account. Please try again.",
		})
		return
	}

	session := tokens["session"].(*models.Session)
	accessToken := tokens["accessToken"].(string)
	accessPayload := tokens["accessPayload"].(*token.Payload)
	refreshToken := tokens["refreshToken"].(string)
	refreshPayload := tokens["refreshPayload"].(*token.Payload)

	ctx.JSON(http.StatusCreated, gin.H{
		"Success": true,
		"User":    userSignUp,
		"Session": gin.H{
			"SessionID":             session.ID,
			"AccessToken":           accessToken,
			"AccessTokenExpiresAt":  accessPayload.ExpiredAt,
			"RefreshToken":          refreshToken,
			"RefreshTokenExpiresAt": refreshPayload.ExpiredAt,
		},
		"Message": "Account created successfully! Please check your email to verify your account.",
	})
}
