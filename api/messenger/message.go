package messenger

import (
	"fmt"
	db "khelogames/database"
	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getMessageByReceiverRequest struct {
	ReceiverPublicID string `uri:"receiver_public_id"`
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

	receiverPublicID, err := uuid.Parse(req.ReceiverPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	authToken := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.GetMessageByReceiverParams{
		SenderID:   authToken.PublicID,
		ReceiverID: receiverPublicID,
	}

	s.logger.Debug("message by receiver arg: ", arg)

	messageContent, err := s.store.GetMessageByReceiver(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get message by receiver", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.logger.Debug("get message by receiver: ", messageContent)

	broadcastMessage := fmt.Sprintf("User: %s retrieved messages from %s", authToken.PublicID, req.ReceiverPublicID)
	s.broadcast <- []byte(broadcastMessage)

	ctx.JSON(http.StatusAccepted, messageContent)
	return
}

func (s *MessageServer) GetUserByMessageSendFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	messageUserName, err := s.store.GetUserByMessageSend(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get user by message send: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Debug("get username by message sent: ", messageUserName)

	ctx.JSON(http.StatusAccepted, messageUserName)
	return
}

func (s *MessageServer) DeleteScheduleMessageFunc(ctx *gin.Context) {

	response, err := s.store.ScheduledDeleteMessage(ctx)
	if err != nil {
		s.logger.Error("Failed to delete message: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
