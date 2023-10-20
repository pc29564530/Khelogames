package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"
)

type addJoinCommunityRequest struct {
	CommunityName string `uri:"community_name"`
}

func (server *Server) addJoinCommunity(ctx *gin.Context) {
	var req addJoinCommunityRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println("CommunityName: ", req.CommunityName)

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.AddJoinCommunityParams{
		CommunityName: req.CommunityName,
		Username:      authPayload.Username,
	}
	fmt.Println(arg)
	communityUser, err := server.store.AddJoinCommunity(ctx, arg)
	fmt.Println("community User: ", communityUser)
	fmt.Println("error: ", err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println("CommunityUser: ", communityUser)

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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	communityUserList, err := server.store.GetUserByCommunity(ctx, req.CommunityName)
	if err != nil {
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
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, communityList)
}
