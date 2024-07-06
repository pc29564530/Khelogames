package handlers

import (
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

	s.logger.Debug("bind the request: %v", req)
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.AddJoinCommunityParams{
		CommunityName: req.CommunityName,
		Username:      authPayload.Username,
	}
	s.logger.Debug("params arg: %v", arg)

	communityUser, err := s.store.AddJoinCommunity(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("successfully join community: %v", communityUser)

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
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)

	communityUserList, err := s.store.GetUserByCommunity(ctx, req.CommunityName)
	if err != nil {
		s.logger.Error("Failed to get user by community: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("user by community: %v", communityUserList)

	ctx.JSON(http.StatusOK, communityUserList)
	return
}

// get the community joined by the users
func (s *JoinCommunityServer) GetCommunityByUserFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	communityList, err := s.store.GetCommunityByUser(ctx, authPayload.Username)
	if err != nil {
		s.logger.Error("Failed to get community by user: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	s.logger.Debug("community by user: %v", communityList)

	ctx.JSON(http.StatusOK, communityList)
}
