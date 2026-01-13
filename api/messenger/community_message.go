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

type createCommunityMessageRequest struct {
	CommuntiyPublicID string `json:"community_public_id" binding:"required"`
	Name              string `json:"name" binding:"required,min=1,max=100"`
	Content           string `json:"content" binding:"required,min=1"`
	MediaUrl          string `json:"media_url"`
	MediaType         string `json:"media_type"`
}

func (s *MessageServer) CreateCommunityMessageFunc(ctx *gin.Context) {
	var req createCommunityMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	communityPublicID, err := uuid.Parse(req.CommuntiyPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"community_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.CreateCommunityMessageParams{
		CommuntiyPublicID: communityPublicID,
		SenderPublicID:    authPayload.PublicID,
		Name:              req.Name,
		Content:           req.Content,
		MediaUrl:          req.MediaUrl,
		MediaType:         req.MediaType,
	}

	s.logger.Debug("Create community message params: ", arg)

	response, err := s.store.CreateCommunityMessage(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create community message: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create community message",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Info("Successfully created community message")
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *MessageServer) GetCommunityMessageFunc(ctx *gin.Context) {
	var req struct {
		CommunityPublicID string `uri:"community_public_id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	s.logger.Info("Received request to get community message")

	communityPublicID, err := uuid.Parse(req.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"community_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	response, err := s.store.GetCommuntiyMessage(ctx, communityPublicID) //spelling mistake
	if err != nil {
		s.logger.Error("Failed to get community message: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get community message",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Info("Successfully retrieved community message")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

func (s *MessageServer) GetCommunityByMessageFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get community by message")

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	response, err := s.store.GetCommunityByMessage(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get community by message: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get community by message",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Info("Successfully retrieved community by message")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
