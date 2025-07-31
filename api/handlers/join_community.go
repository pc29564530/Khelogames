package handlers

import (
	"khelogames/pkg"
	"khelogames/token"
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

	// Add user to join_community table
	communityUser, err := s.store.AddJoinCommunity(ctx, communityPublicID, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to join community: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join community"})
		return
	}
	s.logger.Debug("Successfully joined community: ", communityUser)

	// Increment member count
	err = s.store.IncrementCommunityMemberCount(ctx, communityPublicID)
	if err != nil {
		s.logger.Error("Failed to increment member count: ", err)
	}

	// Return result
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully joined", "member": communityUser})
}

// get the community joined by the users
func (s *HandlersServer) GetCommunityByUserFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	communityList, err := s.store.GetCommunityByUser(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get community by user: ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	s.logger.Debug("community by user: ", communityList)

	ctx.JSON(http.StatusOK, communityList)
}
