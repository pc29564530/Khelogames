package api

import (
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

	authToken := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateNewMessageParams{
		Content:          req.Content,
		IsSeen:           req.IsSeen,
		SenderUsername:   authToken.Username,
		ReceiverUsername: req.ReceiverUsername,
	}

	messageContent, err := server.store.CreateNewMessage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, messageContent)
	return
}

type getMessageByReceiverRequest struct {
	ReceiverUsername string `json:"receiver_username"`
}

func (server *Server) getMessageByReceiver(ctx *gin.Context) {
	var req getMessageByReceiverRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authToken := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateNewMessageParams{
		SenderUsername:   authToken.Username,
		ReceiverUsername: req.ReceiverUsername,
	}

	messageContent, err := server.store.CreateNewMessage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, messageContent)
	return
}
