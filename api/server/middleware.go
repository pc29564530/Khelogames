package server

import (
	"khelogames/core/token"
	"khelogames/logger"
	"khelogames/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func authMiddleware(tokenMaker token.Maker, log *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(pkg.AuthorizationHeaderKey)
		if authorizationHeader == "" {
			abortAuth(ctx, log, "AUTH_HEADER_MISSING", "Authorization header is required")
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			abortAuth(ctx, log, "INVALID_AUTH_HEADER", "Invalid authorization header format")
			return
		}

		if strings.ToLower(fields[0]) != pkg.AuthorizationTypeBearer {
			abortAuth(ctx, log, "INVALID_AUTH_TYPE", "Bearer token required")
			return
		}

		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			abortAuth(ctx, log, "TOKEN_INVALID", "Invalid or expired token")
			return
		}

		ctx.Set(pkg.AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

func abortAuth(ctx *gin.Context, log *logger.Logger, code, message string) {
	log.Error(message)
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
