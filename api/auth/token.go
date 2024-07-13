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

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *AuthServer) RenewAccessTokenFunc(ctx *gin.Context) {
	config, err := util.LoadConfig(".")
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusBadRequest, (err))
		return
	}
	s.logger.Debug(fmt.Sprintf("access token request: %v", req))
	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		s.logger.Error("Failed to verify token: %v", err)
		ctx.JSON(http.StatusUnauthorized, (err))
		return
	}

	s.logger.Debug(fmt.Sprintf("verify refresh token: %v", refreshPayload))

	session, err := s.store.GetSessions(ctx, refreshPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("successfully get sessions")
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	if session.Username != refreshPayload.Username {
		s.logger.Error("Failed to get correct session user")
		return
	}

	if session.RefreshToken != req.RefreshToken {
		s.logger.Error("mismatched session token")
		return
	}

	if time.Now().After(session.ExpiresAt) {
		s.logger.Error("Expired session")
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		refreshPayload.Username,
		config.AccessTokenDuration,
	)
	if err != nil {
		s.logger.Error("Failed to create token: %v", err)
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
