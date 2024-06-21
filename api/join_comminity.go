package api

import (
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addJoinCommunityRequest struct {
	CommunityName string `json:"community_name"`
}

func (server *Server) addJoinCommunity(ctx *gin.Context) {
	var req addJoinCommunityRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.AddJoinCommunityParams{
		CommunityName: req.CommunityName,
		Username:      authPayload.Username,
	}

	communityUser, err := server.store.AddJoinCommunity(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, communityUser)
	return
}

type getUserByCommunityRequest struct {
	CommunityName string `uri:"community_name"`
}

func (server *Server) getUserByCommunity(ctx *gin.Context) {
	var req getUserByCommunityRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		server.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	communityUserList, err := server.store.GetUserByCommunity(ctx, req.CommunityName)
	if err != nil {
		server.logger.Error("Failed to get user by community: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, communityUserList)
	return
}

//get the community joined by the users

func (server *Server) getCommunityByUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	communityList, err := server.store.GetCommunityByUser(ctx, authPayload.Username)
	if err != nil {
		server.logger.Error("Failed to get community by user: %v", err)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, communityList)
}
