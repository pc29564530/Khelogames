package auth

import (
	"database/sql"
	"fmt"
	db "khelogames/database"
	"khelogames/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
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

	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		s.logger.Error("Failed to verify refresh token: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	s.logger.Debug(fmt.Sprintf("verify refresh token: %v", refreshPayload))

	session, err := s.store.GetSessions(ctx, refreshPayload.ID)
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

	if session.PublicID != refreshPayload.PublicID {
		s.logger.Error("Session user mismatch")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	if session.RefreshToken != req.RefreshToken {
		s.logger.Error("Session token mismatch")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if time.Now().After(session.ExpiresAt) {
		s.logger.Error("Session expired")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		refreshPayload.PublicID,
		config.AccessTokenDuration,
	)
	if err != nil {
		s.logger.Error("Failed to create token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func CreateNewToken(ctx *gin.Context, publicID uuid.UUID, userID int32, s *AuthServer, tx *sql.Tx) map[string]interface{} {
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		publicID,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		s.logger.Error("Failed to create access token: ", err)
		return nil
	}
	s.logger.Debug("created a accesstoken: ", accessToken)

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		publicID,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		s.logger.Error("Failed to create refresh token: ", err)
		return nil
	}

	s.logger.Debug("created a refresh token: ", refreshToken)

	session, err := s.store.CreateSessions(ctx, db.CreateSessionsParams{
		UserID:       int32(userID),
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
	})

	if err != nil {
		s.logger.Error("Failed to create session: ", err)
		return nil
	}

	return map[string]interface{}{
		"accessToken":    accessToken,
		"accessPayload":  accessPayload,
		"refreshToken":   refreshToken,
		"refreshPayload": refreshPayload,
		"session":        session,
	}
}
