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
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)

	communityPublicID, err := uuid.Parse(req.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	communityUser, err := s.store.AddJoinCommunity(ctx, communityPublicID, authPayload.PublicID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("successfully join community: ", communityUser)

	ctx.JSON(http.StatusOK, communityUser)
	return
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
