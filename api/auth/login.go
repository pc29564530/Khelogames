package auth

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/token"
	"khelogames/util"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
}

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

func NewLoginServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config) *LoginServer {
	return &LoginServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config}
}

func (s *LoginServer) CreateLoginFunc(ctx *gin.Context) {
	var req createLoginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("No row found: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	fmt.Println("linen no 64: ", req.Username)
	fmt.Println("line no 65: ", s)
	fmt.Println("Line no 66: ", s.store)
	user, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("No row found: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		fmt.Errorf("Failed to get the user: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	fmt.Println("Line no 68 Login: ", user)

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		fmt.Errorf("Failed to match password: %v", err)
		ctx.JSON(http.StatusUnauthorized, (err))
		return
	}

	fmt.Println("Line no 8: N7", s.tokenMaker)

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		fmt.Errorf("Failed to create access token: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Username,
		s.config.RefreshTokenDuration,
	)

	if err != nil {
		fmt.Errorf("Failed to create refresh token: %v", err)
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
		fmt.Errorf("Failed to create session: %v", err)
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
	fmt.Println("Logged in successfully")
	ctx.JSON(http.StatusOK, rsp)
	return
}

func (s *LoginServer) VerifyMobileAndPasswordFunc(ctx *gin.Context, username string, password string, userData db.User) error {
	var err error
	if userData.Username != username {
		fmt.Errorf("Failed to verify mobile and password: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return err
	}
	pass, err := util.HashPassword(password)
	if err != nil {
		s.logger.Debug("Failed to convert password: %v", err)
		return err
	}
	fmt.Println(pass)
	fmt.Println(userData.HashedPassword)
	err = util.CheckPassword(pass, userData.HashedPassword)
	if err != nil {
		fmt.Errorf("Failed to verify mobile and password: %v", err)
		return err
	}
	return nil
}

func (s *LoginServer) GenerateSessionTokenFunc() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		fmt.Errorf("Failed to generate session token: %v", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}
