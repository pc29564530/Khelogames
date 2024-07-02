package auth

import (
	"database/sql"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/token"
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

type TokenServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
}

func NewTokenServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker) *TokenServer {
	return &TokenServer{store: store, logger: logger, tokenMaker: tokenMaker}
}

func (s *TokenServer) RenewAccessTokenFunc(ctx *gin.Context) {
	config, err := util.LoadConfig(".")
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusBadRequest, (err))
		return
	}

	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		fmt.Errorf("Failed to verify token: %v", err)
		ctx.JSON(http.StatusUnauthorized, (err))
		return
	}

	session, err := s.store.GetSessions(ctx, refreshPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	if session.Username != refreshPayload.Username {
		fmt.Errorf("Failed to get correct session user")
		return
	}

	if session.RefreshToken != req.RefreshToken {
		fmt.Errorf("mismatched session token")
		return
	}

	if time.Now().After(session.ExpiresAt) {
		fmt.Errorf("Expired session")
		return
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		refreshPayload.Username,
		config.AccessTokenDuration,
	)
	if err != nil {
		fmt.Errorf("Failed to create token: %v", err)
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
