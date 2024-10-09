package auth

import (
	"database/sql"
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
	"google.golang.org/api/idtoken"
)

var googleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo"},
	Endpoint:     google.Endpoint,
}

type getGoogleLoginRequest struct {
	Code string `json"code"`
}

func (s *AuthServer) HandleGoogleCallback(ctx *gin.Context) {

	var req getGoogleLoginRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind the login request : ", err)
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin the transcation: ", err)
		return
	}

	defer tx.Rollback()

	idToken, err := idtoken.Validate(ctx, req.Code, googleOauthConfig.ClientID)
	if err != nil {
		s.logger.Error("Failed to verify idToken: ", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid idToken"})
		return
	}

	// Extract user info from the verified token
	email, ok := idToken.Claims["email"].(string)
	if !ok {
		s.logger.Error("Failed to get email from idToken")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	//Check if user exists
	var emptyUser models.User
	existingUser, err := s.store.GetGoogleMailID(ctx, email)
	if err != nil {
		s.logger.Error("Failed to get mail id: ", err)
	}

	if existingUser != emptyUser {
		authorizationCode(ctx, existingUser.Username, *existingUser.Gmail, existingUser.Role, s, tx, false)
	} else {
		username := generateUsername(email)
		role := "user"
		user, err := s.store.CreateGoogleUser(ctx, username, email, role)
		if err != nil {
			s.logger.Error("Failed to create google user: ", err)
		}

		resp := createUserResponse{
			Username: user.Username,
			Role:     user.Role,
			Gmail:    *user.Gmail,
		}

		s.logger.Debug("successfully created user: ", resp)

		authorizationCode(ctx, user.Username, *user.Gmail, user.Role, s, tx, true)

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
	Username string `json:"username"`
	Role     string `json:"role"`
	Gmail    string `json:"gmail"`
}

type loginUserGoogleResponse struct {
	IsNewUser             bool               `json:"isNewUser"`
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
	Gmail        string `json:"gmail"`
}

func authorizationCode(ctx *gin.Context, username string, gmail string, role string, s *AuthServer, tx *sql.Tx, isNewUser bool) {

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
		IsNewUser:             isNewUser,
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: userGoogleResponse{
			Username: username,
			Role:     role,
			Gmail:    gmail,
		},
	}
	s.logger.Info("User logged in successfully")
	ctx.JSON(http.StatusAccepted, rsp)
	return
}
