package token

import (
	"database/sql"
	"fmt"
	errorhandler "khelogames/error_handler"
	"khelogames/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *TokenServer) RenewAccessTokenFunc(ctx *gin.Context) {
	config, err := util.LoadConfig(".")
	if err != nil {
		s.logger.Error("Failed to load config: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "LOAD_CONFIG_ERROR",
			"message": "Failed to load config",
		})
		return
	}

	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Error("Failed to bind request: %v", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	s.logger.Debug(fmt.Sprintf("access token request: %v", req))

	// Verify the refresh token
	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		s.logger.Error("Failed to verify refresh token: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTHORIZATION_ERROR",
			"message": "Failed to verify token",
		})
		return
	}

	s.logger.Debug(fmt.Sprintf("verify refresh token: %v", refreshPayload))

	// Check if refresh token is expired (this should be redundant after VerifyToken, but good for safety)
	if time.Now().After(refreshPayload.ExpiredAt) {
		s.logger.Error("Refresh token expired")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTHORIZATION_ERROR",
			"message": "Failed to check refresh token is expired",
		})
		return
	}

	// Get session by refresh token
	session, err := s.store.GetSessionByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("Session not found: %v", err)
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    "DATABASE_ERROR",
				"message": "Failed to get sessions",
			})
			return
		}
		s.logger.Error("Failed to get session: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Verify session belongs to the same user
	if session.UserID != refreshPayload.UserID {
		s.logger.Error("Session user mismatch")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUHTORIZARION_ERROR",
			"message": "Failed to verify user",
		})
		return
	}

	// Verify the refresh token matches
	if session.RefreshToken != req.RefreshToken {
		s.logger.Error("Session token mismatch")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTHORIZATION_ERROR",
			"message": "Session token mismatch",
		})
		return
	}

	// Check session expiry (this is the database session expiry, not the JWT expiry)
	if time.Now().After(session.ExpiresAt) {
		s.logger.Error("Database session expired")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"code":    "AUTHORIZATION_ERROR",
			"message": "Session token expired",
		})
		return
	}

	// Create new access token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		refreshPayload.PublicID,
		refreshPayload.UserID,
		config.AccessTokenDuration,
	)
	if err != nil {
		s.logger.Error("Failed to create token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "DATABASE_ERROR",
			"message": "Failed to create token",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"AccessToken":          accessToken,
		"AccessTokenExpiresAt": accessPayload.ExpiredAt,
	})
}

// func (s *TokenServer) CreateNewToken(ctx *gin.Context, publicID uuid.UUID, userID int32, tx *sql.Tx) map[string]interface{} {
// 	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
// 		publicID,
// 		userID,
// 		s.config.AccessTokenDuration,
// 	)
// 	if err != nil {
// 		s.logger.Error("Failed to create access token: ", err)
// 		return nil
// 	}
// 	s.logger.Debug("created a accesstoken: ", accessToken)

// 	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
// 		publicID,
// 		userID,
// 		s.config.RefreshTokenDuration,
// 	)
// 	if err != nil {
// 		s.logger.Error("Failed to create refresh token: ", err)
// 		return nil
// 	}

// 	s.logger.Debug("created a refresh token: ", refreshToken)

// 	session, err := s.store.CreateSessions(ctx, db.CreateSessionsParams{
// 		UserID:       int32(userID),
// 		RefreshToken: refreshToken,
// 		UserAgent:    ctx.Request.UserAgent(),
// 		ClientIp:     ctx.ClientIP(),
// 	})

// 	if err != nil {
// 		s.logger.Error("Failed to create session: ", err)
// 		return nil
// 	}

// 	return map[string]interface{}{
// 		"accessToken":    accessToken,
// 		"accessPayload":  accessPayload,
// 		"refreshToken":   refreshToken,
// 		"refreshPayload": refreshPayload,
// 		"session":        session,
// 	}
// }
