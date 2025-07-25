package auth

import (
	db "khelogames/database"
	"khelogames/database/models"
	"khelogames/token"
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
			"message": "Invalid request data",
		})
		return
	}

	// Check if user already exists with this email
	existingUser, err := s.store.GetUsersByGmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		s.logger.Info("User already exists with email: ", req.Email)
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
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

	// Start database transaction
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin the transaction: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}
	defer tx.Rollback()

	// Generate username
	username := GenerateUsername(req.Email)

	// Create the user in database
	userSignUp, err := s.store.CreateEmailSignUp(ctx, req.FullName, username, req.Email, hashPassword)
	if err != nil {
		s.logger.Error("Failed to create email signup: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create account. Please try again.",
		})
		return
	}
	//create a token using user id
	tokens := CreateNewToken(ctx, userSignUp.PublicID, int32(userSignUp.ID), s, tx)

	session := tokens["session"].(models.Session)
	accessToken := tokens["accessToken"].(string)
	accessPayload := tokens["accessPayload"].(*token.Payload)
	refreshToken := tokens["refreshToken"].(string)
	refreshPayload := tokens["refreshPayload"].(*token.Payload)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create account. Please try again.",
		})
		return
	}

	s.logger.Info("Successfully created email sign-up for: ", req.Email)

	arg := db.CreateProfileParams{
		UserID:    int32(userSignUp.ID),
		Bio:       "",
		AvatarUrl: "",
	}

	_, err = s.store.CreateProfile(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create profile: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

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
