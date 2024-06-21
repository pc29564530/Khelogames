package api

import (
	"database/sql"
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

func (server *Server) renewAccessToken(ctx *gin.Context) {
	config, err := util.LoadConfig(".")
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		server.logger.Error("Failed to verify token: %v", err)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSessions(ctx, refreshPayload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			server.logger.Error("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		server.logger.Error("Failed to get correct session user")
		return
	}

	if session.RefreshToken != req.RefreshToken {
		server.logger.Error("mismatched session token")
		return
	}

	if time.Now().After(session.ExpiresAt) {
		server.logger.Error("Expired session")
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		config.AccessTokenDuration,
	)
	if err != nil {
		server.logger.Error("Failed to create token: %v", err)
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
