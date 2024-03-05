package api

import (
	"fmt"
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
	fmt.Println("get message", req)
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println("RecieverUsername: ", req.ReceiverUsername)

	authToken := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	fmt.Println(authToken)
	arg := db.GetMessageByReceiverParams{
		SenderUsername:   authToken.Username,
		ReceiverUsername: req.ReceiverUsername,
	}
	fmt.Println(arg)

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
