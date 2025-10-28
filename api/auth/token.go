package auth

import (
	"database/sql"
	"fmt"
	"khelogames/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *AuthServer) RenewAccessTokenFunc(ctx *gin.Context) {
	config, err := util.LoadConfig(".")
	if err != nil {
		s.logger.Error("Failed to load config: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	s.logger.Debug(fmt.Sprintf("access token request: %v", req))

	// Verify the refresh token
	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		s.logger.Error("Failed to verify refresh token: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	s.logger.Info("Refresh Payload: ", refreshPayload)

	s.logger.Debug(fmt.Sprintf("verify refresh token: %v", refreshPayload))

	// Check if refresh token is expired (this should be redundant after VerifyToken, but good for safety)
	if time.Now().After(refreshPayload.ExpiredAt) {
		s.logger.Error("Refresh token expired")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	// Get session by refresh token
	session, err := s.store.GetSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("Session not found: %v", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			return
		}
		s.logger.Error("Failed to get session: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Verify session belongs to the same user
	if session.UserID != refreshPayload.UserID {
		s.logger.Error("Session user mismatch")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	// Verify the refresh token matches
	if session.RefreshToken != req.RefreshToken {
		s.logger.Error("Session token mismatch")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Check session expiry (this is the database session expiry, not the JWT expiry)
	if time.Now().After(session.ExpiresAt) {
		s.logger.Error("Database session expired")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
		return
	}

	// Create new access token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		refreshPayload.PublicID,
		refreshPayload.UserID,
		config.AccessTokenDuration,
	)
	if err != nil {
		s.logger.Error("Failed to create token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"AccessToken":          accessToken,
		"AccessTokenExpiresAt": accessPayload.ExpiredAt,
	})
}
