package messenger

import (
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/pkg"
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	s.logger.Debug("message receiver username: ", err)

	receiverPublicID, err := uuid.Parse(req.ReceiverPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get message by receiver",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, messageContent)
}

func (s *MessageServer) GetUserByMessageSendFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	messageUser, err := s.store.GetUserByMessageSend(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get user by message send: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get user by message send",
		})
		return
	}

	s.logger.Debug("Get message by user: ", messageUser)

	ctx.JSON(http.StatusAccepted, messageUser)
}

func (s *MessageServer) DeleteScheduleMessageFunc(ctx *gin.Context) {

	response, err := s.store.ScheduledDeleteMessage(ctx)
	if err != nil {
		s.logger.Error("Failed to delete message: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to delete message",
		})
		return
	}
	ctx.JSON(http.StatusAccepted, response)
}
