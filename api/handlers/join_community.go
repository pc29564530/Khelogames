package handlers

import (
	db "khelogames/database"

	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addJoinCommunityRequest struct {
	CommunityName string `json:"community_name"`
}

func (s *HandlersServer) AddJoinCommunityFunc(ctx *gin.Context) {
	var req addJoinCommunityRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.AddJoinCommunityParams{
		CommunityName: req.CommunityName,
		Username:      authPayload.Username,
	}
	s.logger.Debug("params arg: ", arg)

	communityUser, err := s.store.AddJoinCommunity(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("successfully join community: ", communityUser)

	ctx.JSON(http.StatusOK, communityUser)
	return
}

type getUserByCommunityRequest struct {
	CommunityName string `uri:"community_name"`
}

func (s *HandlersServer) GetUserByCommunityFunc(ctx *gin.Context) {
	var req getUserByCommunityRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind : ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)

	communityUserList, err := s.store.GetUserByCommunity(ctx, req.CommunityName)
	if err != nil {
		s.logger.Error("Failed to get user by community: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("user by community: ", communityUserList)

	ctx.JSON(http.StatusOK, communityUserList)
	return
}

// get the community joined by the users
func (s *HandlersServer) GetCommunityByUserFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	communityList, err := s.store.GetCommunityByUser(ctx, authPayload.Username)
	if err != nil {
		s.logger.Error("Failed to get community by user: ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	s.logger.Debug("community by user: ", communityList)

	ctx.JSON(http.StatusOK, communityList)
}

type inActiveUserRequest struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
}

// soft delete of user of community
func (s *HandlersServer) InActiveUserFromCommunityFunc(ctx *gin.Context) {
	var req inActiveUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}

	arg := db.InActiveUserFromCommunityParams{
		Username: req.Username,
		ID:       req.ID,
	}

	response, err := s.store.InActiveUserFromCommunity(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to unfollow the user from the community: ", err)
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
