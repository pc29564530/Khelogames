package handlers

import (
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JoinCommunityServer struct {
	store  *db.Store
	logger *logger.Logger
}

type addJoinCommunityRequest struct {
	CommunityName string `json:"community_name"`
}

func NewJoinCommunityServer(store *db.Store, logger *logger.Logger) *JoinCommunityServer {
	return &JoinCommunityServer{store: store, logger: logger}
}

func (s *JoinCommunityServer) AddJoinCommunityFunc(ctx *gin.Context) {
	var req addJoinCommunityRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.AddJoinCommunityParams{
		CommunityName: req.CommunityName,
		Username:      authPayload.Username,
	}

	communityUser, err := s.store.AddJoinCommunity(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, communityUser)
	return
}

type getUserByCommunityRequest struct {
	CommunityName string `uri:"community_name"`
}

func (s *JoinCommunityServer) GetUserByCommunityFunc(ctx *gin.Context) {
	var req getUserByCommunityRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind : %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	communityUserList, err := s.store.GetUserByCommunity(ctx, req.CommunityName)
	if err != nil {
		fmt.Errorf("Failed to get user by community: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, communityUserList)
	return
}

// get the community joined by the users
func (s *JoinCommunityServer) GetCommunityByUserFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	communityList, err := s.store.GetCommunityByUser(ctx, authPayload.Username)
	if err != nil {
		fmt.Errorf("Failed to get community by user: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}

	ctx.JSON(http.StatusOK, communityList)
}
