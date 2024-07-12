package handlers

import (
	"database/sql"
	db "khelogames/db/sqlc"
	"khelogames/pkg"
	"khelogames/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createCommunitiesRequest struct {
	CommunityName string `json:"communityName"`
	Description   string `json:"description"`
	CommunityType string `json:"communityType"`
}

// Create communities function
func (s *HandlersServer) CreateCommunitesFunc(ctx *gin.Context) {
	var req createCommunitiesRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateCommunityParams{
		Owner:           authPayload.Username,
		CommunitiesName: req.CommunityName,
		Description:     req.Description,
		CommunityType:   req.CommunityType,
	}
	s.logger.Debug("params arg: %v", arg)

	communities, err := s.store.CreateCommunity(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create community: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("created community: %v", communities)
	ctx.JSON(http.StatusOK, communities)
	return
}

type getCommunityRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
type getCommunityResponse struct {
	CommunitiesName string    `json:"communityName"`
	Description     string    `json:"description"`
	CommunityType   string    `json:"communityType"`
	CreatedAt       time.Time `json:"created_at"`
}

// get Community by id.
func (s *HandlersServer) GetCommunityFunc(ctx *gin.Context) {
	var req getCommunityRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)

	community, err := s.store.GetCommunity(ctx, req.ID)
	if err != nil {
		s.logger.Error("Failed to get community: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	resp := getCommunityResponse{
		CommunitiesName: community.CommunitiesName,
		Description:     community.Description,
		CommunityType:   community.CommunityType,
		CreatedAt:       community.CreatedAt,
	}
	s.logger.Debug("get community response: %v", resp)

	ctx.JSON(http.StatusOK, resp)
	return
}

// Get all communities by owner.
func (s *HandlersServer) GetAllCommunitiesFunc(ctx *gin.Context) {

	user, err := s.store.GetAllCommunities(ctx)
	if err != nil {
		s.logger.Error("Failed to  get communities: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	s.logger.Debug("get all community: %v", user)
	ctx.JSON(http.StatusOK, user)
	return
}

// Get all users that have joined a particular communities
type getCommunitiesMemberRequest struct {
	CommunitiesName string `uri:"communities_name"`
}

func (s *HandlersServer) GetCommunitiesMemberFunc(ctx *gin.Context) {
	var req getCommunitiesMemberRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)

	usersList, err := s.store.GetCommunitiesMember(ctx, req.CommunitiesName)
	if err != nil {
		s.logger.Error("Failed to get community member: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get community member: %v", usersList)
	ctx.JSON(http.StatusOK, usersList)
}

type getCommunityByCommunityNameRequest struct {
	CommunitiesName string `uri:"communities_name"`
}

func (s *HandlersServer) GetCommunityByCommunityNameFunc(ctx *gin.Context) {
	var req getCommunityByCommunityNameRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row error %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)

	usersList, err := s.store.GetCommunityByCommunityName(ctx, req.CommunitiesName)
	if err != nil {
		s.logger.Error("Failed to get community by community name: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get community by commmunity name: %v", usersList)

	ctx.JSON(http.StatusOK, usersList)
}
