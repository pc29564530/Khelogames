package messenger

import (
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createCommunityMessageRequest struct {
	CommuntiyPublicID string `json:"community_public_id"`
	Name              string `json:"name"`
	Content           string `json:"content"`
	MediaUrl          string `json:"media_url"`
	MediaType         string `json:"media_type"`
}

func (s *MessageServer) CreateCommunityMessageFunc(ctx *gin.Context) {

	var req createCommunityMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Error("Failed to bind JSON: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Debug("Successfully bind: ", req)

	communityPublicID, err := uuid.Parse(req.CommuntiyPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
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
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Info("Successfully created community message")
	ctx.JSON(http.StatusAccepted, response)
}

func (s *MessageServer) GetCommunityMessageFunc(ctx *gin.Context) {
	var req struct {
		CommunityPublicID string `json:"community_public_id"`
	}
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind URI: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	s.logger.Info("Received request to get community message")

	communityPublicID, err := uuid.Parse(req.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	response, err := s.store.GetCommuntiyMessage(ctx, communityPublicID) //spelling mistake
	if err != nil {
		s.logger.Error("Failed to get community message: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved community message")
	ctx.JSON(http.StatusOK, response)
}

func (s *MessageServer) GetCommunityByMessageFunc(ctx *gin.Context) {
	s.logger.Info("Received request to get community by message")

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	response, err := s.store.GetCommunityByMessage(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get community by message: ", err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	s.logger.Info("Successfully retrieved community by message")
	ctx.JSON(http.StatusOK, response)
}
