package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	db "khelogames/database"
	"khelogames/database/models"
	utils "khelogames/util"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo"},
	Endpoint:     google.Endpoint,
}

func (s *AuthServer) HandleGoogleLogin(ctx *gin.Context) {

	url := googleOauthConfig.AuthCodeURL("randomstate")
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *AuthServer) HandleGoogleCallback(ctx *gin.Context) {

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin the transcation: ", err)
		return
	}

	defer tx.Rollback()

	code := ctx.Query("code")

	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		s.logger.Error("Failed to exchange token", err)
		return
	}

	client := googleOauthConfig.Client(ctx, token)
	resp, err := client.Get(googleOauthConfig.Scopes[0])
	if err != nil {
		s.logger.Error("Failed to get user info: ", err)
		return
	}

	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		s.logger.Error("Failed to decode user info: ", err)
		return
	}

	//Check if user exists
	var emptyUser models.User
	existingUser, err := s.store.GetGoogleMailID(ctx, userInfo.Email)
	if err != nil {
		s.logger.Error("Failed to get mail id: ", err)
	}

	if existingUser != emptyUser {

	}
	if existingUser == emptyUser {
		username := generateUsername(userInfo.Email)
		role := "user"
		user, err := s.store.CreateGoogleUser(ctx, username, userInfo.Email, role)
		if err != nil {
			s.logger.Error("Failed to create google user: ", err)
		}

		resp := createUserResponse{
			Username:     user.Username,
			MobileNumber: *user.MobileNumber,
			Role:         user.Role,
			Email:        *user.Email,
		}

		s.logger.Debug("successfully created user: ", resp)

		authorizationCode(ctx, user.Username, *user.Email, user.Role, s, tx)

		argProfile := db.CreateProfileParams{
			Owner:     resp.Username,
			FullName:  "",
			Bio:       "",
			AvatarUrl: "",
		}

		_, err = s.store.CreateProfile(ctx, argProfile)
		if err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, (err))
			return
		}

		err = tx.Commit()
		if err != nil {
			s.logger.Error("Failed to commit transcation: ", err)
			return
		}

		s.logger.Info("Profile created successfully")
		s.logger.Info("Successfully created the user")
		return
	}

}

func generateUsername(mail string) string {
	// Extract the local part of the email address
	localPart := strings.Split(mail, "@")[0]

	// Remove dots and underscores from the local part
	localPart = strings.ReplaceAll(localPart, ".", "")
	localPart = strings.ReplaceAll(localPart, "_", "")

	// Generate a random number
	randomNumber := utils.RandomString(6)

	username := fmt.Sprintf("%s%s", localPart, randomNumber)

	return username
}

type createUserResponse struct {
	Username     string    `json:"username"`
	MobileNumber string    `json:"mobile_number"`
	CreatedAt    time.Time `json:"created_at"`
	Role         string    `json:"role"`
	Email        string    `json:"email"`
}

type loginUserGoogleResponse struct {
	SessionID             uuid.UUID          `json:"session_id"`
	AccessToken           string             `json:"access_token"`
	AccessTokenExpiresAt  time.Time          `json:"access_token_expires_at"`
	RefreshToken          string             `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time          `json:"refresh_token_expires_at"`
	User                  userGoogleResponse `json:"user"`
}

type userGoogleResponse struct {
	Username     string `json:"username"`
	MobileNumber string `json:"mobile_number"`
	Role         string `json:"role"`
	Email        string `json:"email"`
}

func authorizationCode(ctx *gin.Context, username string, email string, role string, s *AuthServer, tx *sql.Tx) {

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create access token", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("created a accesstoken: ", accessToken)

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		username,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create refresh token", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Debug("created a refresh token: ", refreshToken)

	session, err := s.store.CreateSessions(ctx, db.CreateSessionsParams{
		ID:           refreshPayload.ID,
		Username:     username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create session", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	rsp := loginUserGoogleResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: userGoogleResponse{
			Username: username,
			Role:     role,
			Email:    email,
		},
	}
	s.logger.Info("User logged in successfully")
	ctx.JSON(http.StatusAccepted, rsp)
	return
}
