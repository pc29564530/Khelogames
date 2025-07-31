package handlers

import (
	"database/sql"
	db "khelogames/database"
	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	utils "khelogames/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createCommunitiesRequest struct {
	CommunityName string `json:"communityName"`
	Description   string `json:"description"`
	CommunityType string `json:"communityType"`
	AvatarUrl     string `json:"avatar_url"`
	CoverImageUrl string `json:"cover_image_url"`
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
		UserPublicID:  authPayload.PublicID,
		Name:          req.CommunityName,
		Slug:          utils.GenerateSlug(req.CommunityName),
		Description:   req.Description,
		CommunityType: req.CommunityType,
		AvatarUrl:     req.AvatarUrl,
		CoverImageUrl: req.CoverImageUrl,
	}
	s.logger.Debug("params arg: ", arg)

	communities, err := s.store.CreateCommunity(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create community: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("created community : ", communities)
	ctx.JSON(http.StatusOK, communities)
	return
}

type getCommunityRequest struct {
	PublicID string `uri:"public_id" binding:"required,min=1"`
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

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	community, err := s.store.GetCommunity(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get community: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, community)
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
	CommunityPublicID string `uri:"community_public_id"`
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

	communityPublicID, err := uuid.Parse(req.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	usersList, err := s.store.GetCommunitiesMember(ctx, communityPublicID)
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

	usersList, err := s.store.GetCommunityByCommunityName(ctx, req.CommunitiesName)
	if err != nil {
		s.logger.Error("Failed to get community by community name: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get community by commmunity name: ", usersList)

	ctx.JSON(http.StatusOK, usersList)
}

type updateCommunityName struct {
	CommunityName string `json:"community_name"`
}

func (s *HandlersServer) UpdateCommunityByCommunityNameFunc(ctx *gin.Context) {
	var reqUri communityPublicIDReq
	err := ctx.ShouldBindUri(&reqUri)
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

	publicID, err := uuid.Parse(reqUri.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var reqJSON updateCommunityName
	err = ctx.ShouldBindJSON(&reqJSON)
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

	response, err := s.store.UpdateCommunityName(ctx, publicID, reqJSON.CommunityName)
	if err != nil {
		s.logger.Error("Failed to update the community name: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

type communityPublicIDReq struct {
	CommunityPublicID string `uri:"community_public_id"`
}

type updateCommunityDescription struct {
	CommunityDescription string `json:"community_description"`
}

func (s *HandlersServer) UpdateCommunityByDescriptionFunc(ctx *gin.Context) {
	var reqUri communityPublicIDReq
	err := ctx.ShouldBindUri(&reqUri)
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

	publicID, err := uuid.Parse(reqUri.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var reqJSON updateCommunityDescription
	err = ctx.ShouldBindJSON(&reqJSON)
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

	response, err := s.store.UpdateCommunityDescription(ctx, publicID, reqJSON.CommunityDescription)
	if err != nil {
		s.logger.Error("Failed to update the  description: ", err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
