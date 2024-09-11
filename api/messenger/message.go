package messenger

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getMessageByReceiverRequest struct {
	ReceiverUsername string `uri:"receiver_username"`
}

func (s *MessageServer) GetMessageByReceiverFunc(ctx *gin.Context) {
	var req getMessageByReceiverRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind URI", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.logger.Debug("message receiver username: ", err)

	authToken := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.GetMessageByReceiverParams{
		SenderUsername:   authToken.Username,
		ReceiverUsername: req.ReceiverUsername,
	}

	s.logger.Debug("message by receiver arg: ", arg)

	messageContent, err := s.store.GetMessageByReceiver(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get message by receiver", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.logger.Debug("get message by receiver: ", messageContent)

	broadcastMessage := fmt.Sprintf("User: %s retrieved messages from %s", authToken.Username, req.ReceiverUsername)
	s.broadcast <- []byte(broadcastMessage)

	ctx.JSON(http.StatusAccepted, messageContent)
	return
}

func (s *MessageServer) GetUserByMessageSendFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	messageUserName, err := s.store.GetUserByMessageSend(ctx, authPayload.Username)
	if err != nil {
		s.logger.Error("Failed to get user by message send: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Debug("get username by message sent: ", messageUserName)

	ctx.JSON(http.StatusAccepted, messageUserName)
	return
}
