package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"

	"khelogames/core/token"
	"khelogames/database/models"
	errorhandler "khelogames/error_handler"
	utils "khelogames/util"
)

func getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

type getGoogleLoginRequest struct {
	Code string `json:"code"`
}

func (s *AuthServer) HandleGoogleRedirect(ctx *gin.Context) {
	googleOauthConfig := getGoogleOauthConfig()
	url := googleOauthConfig.AuthCodeURL("randomstate")
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *AuthServer) CreateGoogleSignUpFunc(ctx *gin.Context) {
	var req struct {
		IDToken string `json:"id_token"`
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
				"message": "Failed to sign in with Google",
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

	name, _ := payload.Claims["name"].(string)
	if name == "" {
		name = strings.Split(email, "@")[0]
	}

	googleID, ok := payload.Claims["sub"].(string)
	if !ok || googleID == "" {
		s.logger.Error("Failed to get google id from idToken")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Failed to create account",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	avatarUrl, _ := payload.Claims["picture"].(string)

	username := GenerateUsername(email)

	var user *models.Users
	var session *models.Session
	var accessToken string
	var refreshToken string
	var accessPayload *token.Payload
	var refreshPayload *token.Payload
	userAgent := ctx.Request.UserAgent()
	clientIP := ctx.ClientIP()

	// Check if user already exists — sign in or sign up
	existingUser, err := s.store.GetUsersByGmail(ctx, email)
	if err != nil {
		s.logger.Error("Database error while fetching user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to process request",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	if existingUser != nil {
		// --- SIGN IN path: user already exists ---
		tokens, err := token.CreateNewToken(
			ctx,
			s.store,
			s.tokenMaker,
			int32(existingUser.ID),
			existingUser.PublicID,
			s.config.AccessTokenDuration,
			s.config.RefreshTokenDuration,
			userAgent,
			clientIP,
		)
		if err != nil {
			s.logger.Error("Token creation failed",
				"user_public_id", existingUser.PublicID,
				"request_id", ctx.GetString("request_id"),
				"client_ip", clientIP,
				"user_agent", userAgent,
				"error", err,
			)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "AUTH_SERVICE_UNAVAILABLE",
					"message": "Unable to sign in right now. Please try again later.",
				},
				"request_id": ctx.GetString("request_id"),
			})
			return
		}
		session = tokens["session"].(*models.Session)
		accessToken = tokens["accessToken"].(string)
		accessPayload = tokens["accessPayload"].(*token.Payload)
		refreshToken = tokens["refreshToken"].(string)
		refreshPayload = tokens["refreshPayload"].(*token.Payload)
		user = existingUser
		s.logger.Info("Successful Google sign in for user: ", existingUser.PublicID)

	} else {
		// --- SIGN UP path: new user ---
		_, userSignUp, tokens, err := s.txStore.CreateGoogleSignUpTx(
			ctx,
			s.config,
			name,
			username,
			email,
			googleID,
			avatarUrl,
		)
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
		session = tokens["session"].(*models.Session)
		accessToken = tokens["accessToken"].(string)
		accessPayload = tokens["accessPayload"].(*token.Payload)
		refreshToken = tokens["refreshToken"].(string)
		refreshPayload = tokens["refreshPayload"].(*token.Payload)
		user = userSignUp
		s.logger.Info("Successfully created Google sign-up for: ", email)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":               true,
		"sessionID":             session.ID,
		"accessToken":           accessToken,
		"accessTokenExpiresAt":  accessPayload.ExpiredAt,
		"refreshToken":          refreshToken,
		"refreshTokenExpiresAt": refreshPayload.ExpiredAt,
		"user":                  user,
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
