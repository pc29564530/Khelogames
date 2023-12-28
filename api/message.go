package api

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createNewMessageRequest struct {
	Content          string `json:"content"`
	IsSeen           bool   `json:"is_seen"`
	ReceiverUsername string `json:"receiver_username"`
}

func (server *Server) createNewMessage(ctx *gin.Context) {
	var req createNewMessageRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println("New Message: ", req)
	authToken := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateNewMessageParams{
		Content:          req.Content,
		IsSeen:           req.IsSeen,
		SenderUsername:   authToken.Username,
		ReceiverUsername: req.ReceiverUsername,
	}

	fmt.Println(arg)

	messageContent, err := server.store.CreateNewMessage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.broadcast <- []byte("createNewMessage")
	ctx.JSON(http.StatusAccepted, messageContent)
	return
}

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
