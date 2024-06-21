package api

import (
	"database/sql"
	"encoding/base64"
	db "khelogames/db/sqlc"
	"khelogames/util"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userResponse struct {
	Username     string `json:"username"`
	MobileNumber string `json:"mobile_number"`
	Role         string `json:"role"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) createLogin(ctx *gin.Context) {
	var req createLoginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			server.logger.Error("No row found: %v", err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			server.logger.Error("No row found: %v", err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		server.logger.Error("Failed to get the user: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		server.logger.Error("Failed to match password: %v", err)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		server.logger.Error("Failed to create access token: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		server.logger.Error("Failed to create refresh token: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSessions(ctx, db.CreateSessionsParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		server.logger.Error("Failed to create session: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: userResponse{
			Username:     user.Username,
			MobileNumber: user.MobileNumber,
			Role:         user.Role,
		},
	}
	server.logger.Info("Logged in successfully")
	ctx.JSON(http.StatusOK, rsp)
	return
}

func (server *Server) verifyMobileAndPassword(ctx *gin.Context, username string, password string, userData db.User) error {
	var err error
	if userData.Username != username {
		server.logger.Error("Failed to verify mobile and password: %v", err)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return err
	}
	pass, err := util.HashPassword(password)
	if err != nil {
		server.logger.Debug("Failed to convert password: %v", err)
		return err
	}
	server.logger.Info(pass)
	server.logger.Info(userData.HashedPassword)
	err = util.CheckPassword(pass, userData.HashedPassword)
	if err != nil {
		server.logger.Error("Failed to verify mobile and password: %v", err)
		return err
	}
	return nil
}

func (server *Server) generateSessionToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		server.logger.Error("Failed to generate session token: %v", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}
