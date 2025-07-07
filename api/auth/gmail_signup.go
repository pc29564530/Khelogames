package auth

import (
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
	existingUser, err := s.store.GetModifyUserByGmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		s.logger.Info("User already exists with email: ", req.Email)
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email already registered. Please sign in instead.",
		})
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
	userSignUp, err := s.store.CreateEmailSignUp(ctx, username, req.Email, req.FullName, req.Password)
	if err != nil {
		s.logger.Error("Failed to create email signup: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create account. Please try again.",
		})
		return
	}

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

	ctx.JSON(http.StatusCreated, gin.H{
		"Success": true,
		"User":    userSignUp,
		"Message": "Account created successfully! Please check your email to verify your account.",
	})
}
