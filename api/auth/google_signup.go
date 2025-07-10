package auth

import (
	"fmt"
	"khelogames/database/models"
	"khelogames/token"
	utils "khelogames/util"
	"net/http"
	"os"
	"strings"

	db "khelogames/database"

	"github.com/gin-gonic/gin"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

func getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

type getGoogleLoginRequest struct {
	Code string `json"code"`
}

func (s *AuthServer) HandleGoogleRedirect(ctx *gin.Context) {
	googleOauthConfig := getGoogleOauthConfig()
	url := googleOauthConfig.AuthCodeURL("randomstate")
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *AuthServer) CreateGoogleSignUpFunc(ctx *gin.Context) {
	var req struct {
		GoogleID  string `json:"google_id"`
		Email     string `json:"email"`
		FullName  string `json:"full_name"`
		AvatarURL string `json:"avatar_url"`
		IDToken   string `json:"id_token"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind the sign-up request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request data",
		})
		return
	}

	// Start database transaction
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin the transaction: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Internal server error",
		})
		return
	}

	defer tx.Rollback()

	// Validate the ID token
	googleOauthConfig := getGoogleOauthConfig()
	payload, err := idtoken.Validate(ctx, req.IDToken, googleOauthConfig.ClientID)
	if err != nil {
		s.logger.Error("Failed to verify idToken: ", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid Google token",
		})
		return
	}

	// Extract user information from the verified token
	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		s.logger.Error("Failed to get email from idToken")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get user email",
		})
		return
	}

	_, err = s.store.GetUsersByGmail(ctx, email)
	if err == nil {
		s.logger.Info("User already exists with email: ", req.Email)
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email already registered. Please sign in instead.",
		})
		return
	}

	// Get name from token (Google uses 'name' not 'full_name')
	name, ok := payload.Claims["name"].(string)
	if !ok || name == "" {
		s.logger.Error("Failed to get name from idToken")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get user name",
		})
		return
	}

	// Get Google ID from token
	googleID, ok := payload.Claims["sub"].(string)
	if !ok || googleID == "" {
		s.logger.Error("Failed to get google id from idToken")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get Google ID",
		})
		return
	}

	// Get avatar URL (optional)
	_ = payload.Claims["picture"].(string)

	// Verify the data matches what was sent from frontend
	if req.Email != email || req.GoogleID != googleID {
		s.logger.Error("Token data doesn't match request data")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Token validation failed",
		})
		return
	}

	// Generate username
	username := GenerateUsername(email)

	// Create the user in database
	userSignUp, err := s.store.CreateGoogleSignUp(ctx, name, username, email, googleID)
	if err != nil {
		s.logger.Error("Failed to create google signup: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create account",
		})
		return
	}

	tokens := CreateNewToken(ctx, userSignUp.Username, s, tx)
	session := tokens["session"].(models.Session)
	accessToken := tokens["accessToken"].(string)
	accessPayload := tokens["accessPayload"].(*token.Payload)
	refreshToken := tokens["refreshToken"].(string)
	refreshPayload := tokens["refreshPayload"].(*token.Payload)

	arg := db.CreateProfileParams{
		UserID:    int32(userSignUp.ID),
		Username:  userSignUp.Username,
		FullName:  req.FullName,
		Bio:       "",
		AvatarUrl: "",
	}

	_, err = s.store.CreateProfile(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create profile: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	s.logger.Info("Successfully created user_profile")

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create account",
		})
		return
	}

	s.logger.Info("Successfully created Google sign-up for: ", email)
	ctx.JSON(http.StatusCreated, gin.H{
		"Success": true,
		"User":    userSignUp,
		"Session": gin.H{
			"SessionID":             session.ID,
			"AccessToken":           accessToken,
			"AccessTokenExpiresAt":  accessPayload.ExpiredAt,
			"RefreshToken":          refreshToken,
			"RefreshTokenExpiresAt": refreshPayload.ExpiredAt,
		},
		"Message": "Account created successfully",
	})
}

func GenerateUsername(mail string) string {
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
	googleOauthConfig := getGoogleOauthConfig()
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

	existingUser, err := s.store.GetUsersByGmail(ctx, email)
	if err == nil && existingUser == nil {
		s.logger.Info("User does not exits with email: ", email)
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email does not registered. Please sign up instead.",
		})
		return
	}
	//create a token using user id
	tokens := CreateNewToken(ctx, existingUser.Username, s, tx)

	session := tokens["session"].(models.Session)
	accessToken := tokens["accessToken"].(string)
	accessPayload := tokens["accessPayload"].(*token.Payload)
	refreshToken := tokens["refreshToken"].(string)
	refreshPayload := tokens["refreshPayload"].(*token.Payload)

	s.logger.Info("Successfully Sign in using google ")

	ctx.JSON(http.StatusAccepted, gin.H{
		"SessionID":             session.ID,
		"AccessToken":           accessToken,
		"AccessTokenExpiresAt":  accessPayload.ExpiredAt,
		"RefreshToken":          refreshToken,
		"RefreshTokenExpiresAt": refreshPayload.ExpiredAt,
		"User":                  existingUser,
	})
}
