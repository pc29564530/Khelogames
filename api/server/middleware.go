package server

import (
	"errors"
	"khelogames/core/token"
	"khelogames/logger"
	"khelogames/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	var logger *logger.Logger
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(pkg.AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, (err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, (err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != pkg.AuthorizationTypeBearer {
			logger.Error("unsupported authorization type %s", authorizationType)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, (err))
			return
		}

		ctx.Set(pkg.AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
