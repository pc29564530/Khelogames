package auth

import (
	"khelogames/core/token"
	"khelogames/database"
	"khelogames/database/models"
	"khelogames/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

func (s *AuthServer) CreateEmailSignInFunc(ctx *gin.Context) {
	userAgent := ctx.Request.UserAgent()
	clientIP := ctx.ClientIP()
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind the login request : ", err)
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin the transaction: ", err)
		return
	}

	defer tx.Rollback()

	existingUser, err := s.store.GetUsersByGmail(ctx, req.Email)
	if err != nil && existingUser == nil {
		s.logger.Info("User already does not exists with email: ", req.Email)
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email not registered. Please sign up instead.",
		})
		return
	}

	err = util.CheckPassword(req.Password, *existingUser.HashPassword)
	s.logger.Info("Existing User Password: ", *existingUser.HashPassword)
	s.logger.Info("New Sign In : ", req.Password)
	if err != nil {
		s.logger.Info("Email and password does not match: ")
		ctx.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email and password does not match:",
		})
		return
	}

	//create a token using user id
	tokens, err := token.CreateNewToken(ctx,
		s.store,
		s.tokenMaker,
		int32(existingUser.ID),
		existingUser.PublicID,
		s.config.RefreshTokenDuration,
		s.config.AccessTokenDuration,
		userAgent,
		clientIP)
	if err != nil {
		s.logger.Errorf("Failed to create new token: ", err)
		return
	}

	session := tokens["session"].(*models.Session)
	accessToken := tokens["accessToken"].(string)
	accessPayload := tokens["accessPayload"].(*token.Payload)
	refreshToken := tokens["refreshToken"].(string)
	refreshPayload := tokens["refreshPayload"].(*token.Payload)

	s.logger.Info("Successfully Sign in using google ")

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transaction: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"Success":               true,
		"SessionID":             session.ID,
		"AccessToken":           accessToken,
		"AccessTokenExpiresAt":  accessPayload.ExpiredAt,
		"RefreshToken":          refreshToken,
		"RefreshTokenExpiresAt": refreshPayload.ExpiredAt,
		"User":                  existingUser,
	})
}

func (s *AuthServer) CreateGoogleSignIn(ctx *gin.Context) {
	var ct *gin.Context
	var maker token.Maker
	var refreshDuration time.Duration
	var accessDuration time.Duration
	var st *database.Store
	userAgent := ct.Request.UserAgent()
	clientIP := ct.ClientIP()
	googleOauthConfig := getGoogleOauthConfig()
	var req getGoogleLoginRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind the login request : ", err)
		return
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin the transaction: ", err)
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
	tokens, err := token.CreateNewToken(ctx, st, maker, int32(existingUser.ID), existingUser.PublicID, refreshDuration, accessDuration, userAgent, clientIP)
	if err != nil {
		s.logger.Errorf("Failed to create new token: ", err)
		return
	}
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
