package api

import (
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getMessageByReceiverRequest struct {
	ReceiverUsername string `uri:"receiver_username"`
}

func (server *Server) getMessageByReceiver(ctx *gin.Context) {
	var req getMessageByReceiverRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authToken := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.GetMessageByReceiverParams{
		SenderUsername:   authToken.Username,
		ReceiverUsername: req.ReceiverUsername,
	}

	messageContent, err := server.store.GetMessageByReceiver(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.broadcast <- []byte("getMessageByReceiver")

	ctx.JSON(http.StatusAccepted, messageContent)
	return
}

func (server *Server) getUserByMessageSend(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	messageUserName, err := server.store.GetUserByMessageSend(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, messageUserName)
	return
}
