package messenger

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageSever struct {
	store     *db.Store
	logger    *logger.Logger
	broadcast chan []byte
}

func NewMessageServer(store *db.Store, logger *logger.Logger, broadcast chan []byte) *MessageSever {
	return &MessageSever{store: store, logger: logger, broadcast: broadcast}
}

type getMessageByReceiverRequest struct {
	ReceiverUsername string `uri:"receiver_username"`
}

func (s *MessageSever) GetMessageByReceiverFunc(ctx *gin.Context) {
	var req getMessageByReceiverRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind URI", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	authToken := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.GetMessageByReceiverParams{
		SenderUsername:   authToken.Username,
		ReceiverUsername: req.ReceiverUsername,
	}
	messageContent, err := s.store.GetMessageByReceiver(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get message by receiver", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	broadcastMessage := fmt.Sprintf("User: %s retrieved messages from %s", authToken.Username, req.ReceiverUsername)
	s.broadcast <- []byte(broadcastMessage)

	ctx.JSON(http.StatusAccepted, messageContent)
	return
}

func (s *MessageSever) GetUserByMessageSendFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	messageUserName, err := s.store.GetUserByMessageSend(ctx, authPayload.Username)
	if err != nil {
		fmt.Errorf("Failed to get user by message send: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusAccepted, messageUserName)
	return
}
