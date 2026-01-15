package auth

import (
	"fmt"

	"khelogames/core/token"
	"khelogames/database/models"
	errorhandler "khelogames/error_handler"

	utils "khelogames/util"
	"net/http"
	"os"
	"strings"

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

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	// Validate the ID token
	googleOauthConfig := getGoogleOauthConfig()
	payload, err := idtoken.Validate(ctx, req.IDToken, googleOauthConfig.ClientID)
	if err != nil {
		s.logger.Error("Failed to verify idToken: ", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "AUTHENTICATION_ERROR",
				"message": "Failed to sign up",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Extract user information from the verified token
	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		s.logger.Error("Failed to get email from idToken")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Failed to get user email",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Check if user already exists
	_, err = s.store.GetUsersByGmail(ctx, email)
	if err == nil {
		s.logger.Info("User already exists with email: ", req.Email)
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "EMAIL_ALREADY_REGISTERED",
				"message": "Email already registered. Please sign in instead.",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Get name from token
	name, ok := payload.Claims["name"].(string)
	if !ok || name == "" {
		s.logger.Error("Failed to get name from idToken")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":  "VALIDATION_ERROR",
				"error": "Failed to get name",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Get Google ID from token
	googleID, ok := payload.Claims["sub"].(string)
	if !ok || googleID == "" {
		s.logger.Error("Failed to get google id from idToken")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":  "VALIDATION_ERROR",
				"error": "Failed to create account",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Verify the data matches what was sent from frontend
	if req.Email != email || req.GoogleID != googleID {

		s.logger.Error("Token data doesn't match request data")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Failed to create account",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	// Generate username
	username := GenerateUsername(email)

	_, userSignUp, tokens, err := s.txStore.CreateGoogleSignUpTx(
		ctx,
		s.config,
		name,
		username,
		email,
		googleID,
		req.AvatarURL)
	if err != nil {
		s.logger.Error("Failed google signup transaction: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create account",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	session := tokens["session"].(models.Session)
	accessToken := tokens["accessToken"].(string)
	accessPayload := tokens["accessPayload"].(*token.Payload)
	refreshToken := tokens["refreshToken"].(string)
	refreshPayload := tokens["refreshPayload"].(*token.Payload)

	s.logger.Info("Successfully created user_profile")

	s.logger.Info("Successfully created Google sign-up for: ", email)
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true, // Changed from "Success" to "success" for consistency
		"user":    userSignUp,
		"session": gin.H{
			"session_id":               session.ID,
			"access_token":             accessToken,
			"access_token_expires_at":  accessPayload.ExpiredAt,
			"refresh_token":            refreshToken,
			"refresh_token_expires_at": refreshPayload.ExpiredAt,
		},
		"message": "Account created successfully",
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
