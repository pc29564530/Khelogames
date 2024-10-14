package auth

import (
	"database/sql"
	"fmt"
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

func (s *AuthServer) HandleGoogleRedirect(ctx *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("randomstate")
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *AuthServer) CreateGoogleSignUp(ctx *gin.Context) {

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

	s.logger.Info("Successfully Sign up using google ", email)
	ctx.JSON(http.StatusAccepted, email)
	return

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

func (s *AuthServer) CreateGoogleSignIn(ctx *gin.Context) {

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
	gmail, ok := idToken.Claims["email"].(string)
	if !ok {
		s.logger.Error("Failed to get email from idToken")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	user, err := s.store.GetUserByGmail(ctx, gmail)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("User does not exists: ", err)
			return
		}
		s.logger.Error("Failed to get the user by gmail")
		return
	}

	tokens := CreateNewToken(ctx, user.Username, s, tx)
	session := tokens["session"].(map[string]interface{})
	accessToken := tokens["accessToken"].(string)
	accessPayload := tokens["accessPayload"].(map[string]interface{})
	refreshToken := tokens["refreshToken"].(string)
	refreshPayload := tokens["refreshPayload"].(map[string]interface{})

	rsp := loginUserGoogleResponse{
		SessionID:             session["id"].(uuid.UUID),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload["expired_at"].(time.Time),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload["expired_at"].(time.Time),
		User: userGoogleResponse{
			Username:     user.Username,
			MobileNumber: *user.MobileNumber,
			Role:         user.Role,
			Gmail:        *user.Gmail,
		},
	}

	s.logger.Info("Successfully Sign in using google ")

	ctx.JSON(http.StatusAccepted, rsp)
	return
}
