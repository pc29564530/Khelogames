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
			s.logger.Error("No row error: ", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateCommunityParams{
		Owner:           authPayload.Username,
		CommunitiesName: req.CommunityName,
		Description:     req.Description,
		CommunityType:   req.CommunityType,
	}
	s.logger.Debug("params arg: ", arg)

	communities, err := s.store.CreateCommunity(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create community: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("created community: ", communities)
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
			s.logger.Error("No row error: ", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)

	community, err := s.store.GetCommunity(ctx, req.ID)
	if err != nil {
		s.logger.Error("Failed to get community: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	resp := getCommunityResponse{
		CommunitiesName: community.CommunitiesName,
		Description:     community.Description,
		CommunityType:   community.CommunityType,
		CreatedAt:       community.CreatedAt,
	}
	s.logger.Debug("get community response: ", resp)

	ctx.JSON(http.StatusOK, resp)
	return
}

// Get all communities by owner.
func (s *HandlersServer) GetAllCommunitiesFunc(ctx *gin.Context) {

	user, err := s.store.GetAllCommunities(ctx)
	if err != nil {
		s.logger.Error("Failed to  get communities: ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	s.logger.Debug("get all community: ", user)
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
			s.logger.Error("No row error: ", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)

	usersList, err := s.store.GetCommunitiesMember(ctx, req.CommunitiesName)
	if err != nil {
		s.logger.Error("Failed to get community member: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get community member: ", usersList)
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
			s.logger.Error("No row error ", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)

	usersList, err := s.store.GetCommunityByCommunityName(ctx, req.CommunitiesName)
	if err != nil {
		s.logger.Error("Failed to get community by community name: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get community by commmunity name: ", usersList)

	ctx.JSON(http.StatusOK, usersList)
}

type updateCommunityByCommunityNameRequest struct {
	CommunityName string `json:"community_name"`
	ID            int64  `json:"id"`
}

func (s *HandlersServer) UpdateCommunityByCommunityNameFunc(ctx *gin.Context) {
	var req updateCommunityByCommunityNameRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row error ", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)

	arg := db.UpdateCommunityNameParams{
		CommunitiesName: req.CommunityName,
		ID:              req.ID,
	}

	response, err := s.store.UpdateCommunityName(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update the community name: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

type updateCommunityByDescriptionRequest struct {
	Description string `json:"description"`
	ID          int64  `json:"id"`
}

func (s *HandlersServer) UpdateCommunityByDescriptionFunc(ctx *gin.Context) {
	var req updateCommunityByDescriptionRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row error ", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)

	arg := db.UpdateCommunityDescriptionParams{
		Description: req.Description,
		ID:          req.ID,
	}

	response, err := s.store.UpdateCommunityDescription(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to update the  description: ", err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
