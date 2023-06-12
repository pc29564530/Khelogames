package api

import (
	"github.com/gin-gonic/gin"
	"khelogames/token"
	"net/http"
)

func (server *Server) getAllFriends(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	allFriends, err := server.store.GetAllFriends(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, allFriends)
	return
}
