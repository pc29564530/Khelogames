package api

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "khelogames/db/sqlc"
	"khelogames/util"
	"math/rand"
	"net/http"
	"time"
)

type createLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userResponse struct {
	Username     string    `json:"username"`
	MobileNumber string    `json:"mobile_number"`
	createdAt    time.Time `json:"created_at"`
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
	fmt.Println(err)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println(user)

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	//err = verifyMobileAndPassword(ctx, req.Username, req.Password, userData)
	////fmt.Printf("ramram")
	////if err != nil {
	////	ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	////}
	////}
	//if err != nil {
	//	ctx.JSON(http.StatusNotFound, errorResponse(err))
	//}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
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
			createdAt:    user.CreatedAt,
		},
	}

	ctx.JSON(http.StatusOK, rsp)

	fmt.Printf("Loggen in successfully")
	return
}

func verifyMobileAndPassword(ctx *gin.Context, username string, password string, userData db.User) error {
	var err error
	if userData.Username != username {
		fmt.Errorf("username does not exist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return err
	}
	pass, err := util.HashPassword(password)
	if err != nil {
		fmt.Errorf("Not able to convert the password to hashed string:%w", err)
		return err
	}
	fmt.Println(pass)
	fmt.Println(userData.HashedPassword)
	err = util.CheckPassword(pass, userData.HashedPassword)
	if err != nil {
		fmt.Errorf("password does not match: %w", err)
		return err
	}
	return nil
}

func generateSessionToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}
