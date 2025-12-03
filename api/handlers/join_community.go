package handlers

import (
	"khelogames/core/token"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *HandlersServer) AddJoinCommunityFunc(ctx *gin.Context) {
	var req struct {
		CommunityPublicID string `uri:"community_public_id"`
	}

	// Bind the URI parameter
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind URI: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	s.logger.Debug("Bind the request: ", req)

	// Parse UUID
	communityPublicID, err := uuid.Parse(req.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Get auth payload
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	community, err := s.store.GetCommunity(ctx, communityPublicID)
	if err != nil {
		s.logger.Error("Failed to get community: ", err)
		return
	}

	if community.UserID != authPayload.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed"})
		return
	}

	communityUser, err := s.txStore.AddJoinCommunityTx(ctx.Request.Context(), communityPublicID, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to join community: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join community"})
		return
	}

	s.logger.Debug("Successfully joined community: ", communityUser)

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully joined", "member": communityUser})
}

func (s *HandlersServer) GetCommunityByUserFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	communityList, err := s.store.GetCommunityByUser(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get community by user: ", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get community list"})
		return
	}
	s.logger.Debug("community by user: ", communityList)

	ctx.JSON(http.StatusOK, communityList)
}
