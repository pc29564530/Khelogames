package auth

import (
	"database/sql"
	"encoding/base64"
	"fmt"
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

func (s *AuthServer) CreateLoginFunc(ctx *gin.Context) {
	var req createLoginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row found: ", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Debug(fmt.Sprintf("successfully get the login request: ", req))

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transcation: ", err)
		return
	}

	defer tx.Rollback()

	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row found: ", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to get the user: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug(fmt.Sprintf("successfully get the user: ", user))

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		s.logger.Error("Failed to match password: ", err)
		ctx.JSON(http.StatusUnauthorized, (err))
		return
	}

	s.logger.Debug(fmt.Sprintf("successfully check password"))

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		s.logger.Error("Failed to create access token: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug(fmt.Sprintf("successfully create a access token %s and access payload %s", accessToken, accessPayload))
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.RefreshTokenDuration,
	)

	if err != nil {
		s.logger.Error("Failed to create refresh token: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	session, err := s.store.CreateSessions(ctx, db.CreateSessionsParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		s.logger.Error("Failed to create session: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
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
	s.logger.Info("Logged in successfully")
	ctx.JSON(http.StatusOK, rsp)
	return
}

func (s *AuthServer) VerifyMobileAndPasswordFunc(ctx *gin.Context, username string, password string, userData db.User) error {
	var err error
	if userData.Username != username {
		s.logger.Error("Failed to verify mobile and password: ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return err
	}
	s.logger.Debug(fmt.Sprintf("successfully matches the username: ", username))
	pass, err := util.HashPassword(password)
	if err != nil {
		s.logger.Debug("Failed to convert password: ", err)
		return err
	}
	s.logger.Debug(fmt.Sprintf("input password of user: ", pass))
	s.logger.Debug(fmt.Sprintf("store password of the user: ", userData.HashedPassword))
	err = util.CheckPassword(pass, userData.HashedPassword)
	if err != nil {
		s.logger.Error("Failed to verify mobile and password: ", err)
		return err
	}
	s.logger.Info("successfully verify the password")
	return nil
}

func (s *AuthServer) GenerateSessionTokenFunc() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		s.logger.Error("Failed to generate session token: ", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}
