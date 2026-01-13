package messenger

import (
	"khelogames/core/token"
	db "khelogames/database"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type getMessageByReceiverRequest struct {
	ReceiverPublicID string `uri:"receiver_public_id" binding:"required"`
}

func (s *MessageServer) GetMessageByReceiverFunc(ctx *gin.Context) {
	var req getMessageByReceiverRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	receiverPublicID, err := uuid.Parse(req.ReceiverPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"receiver_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
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
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get message by receiver",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    messageContent,
	})
}

func (s *MessageServer) GetUserByMessageSendFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	messageUser, err := s.store.GetUserByMessageSend(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get user by message send: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user by message send",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Debug("Get message by user: ", messageUser)

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    messageUser,
	})
}

func (s *MessageServer) DeleteScheduleMessageFunc(ctx *gin.Context) {
	response, err := s.store.ScheduledDeleteMessage(ctx)
	if err != nil {
		s.logger.Error("Failed to delete message: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to delete message",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}
