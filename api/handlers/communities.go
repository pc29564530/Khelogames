package handlers

import (
	"khelogames/core/token"
	db "khelogames/database"
	"khelogames/pkg"
	"net/http"

	utils "khelogames/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createCommunitiesRequest struct {
	CommunityName string `json:"communityName" binding:"required,min=3,max=100"`
	Description   string `json:"description" binding:"required,min=10,max=500"`
	AvatarUrl     string `json:"avatar_url" binding:"omitempty,url"`
	CoverImageUrl string `json:"cover_image_url" binding:"omitempty,url"`
}

// Create communities function
func (s *HandlersServer) CreateCommunitesFunc(ctx *gin.Context) {
	var req createCommunitiesRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind communities ", err)

		// Provide more specific validation error messages
		errorMessage := "Invalid request format"
		if err.Error() != "" {
			// Parse validation errors to provide specific feedback
			errorMessage = err.Error()
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": errorMessage,
		})
		return
	}
	s.logger.Debug("bind the request: ", req)
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	// Sanitize user text inputs to prevent XSS
	sanitizedName := utils.SanitizeString(req.CommunityName)
	sanitizedDesc := utils.SanitizeString(req.Description)

	arg := db.CreateCommunityParams{
		UserPublicID:  authPayload.PublicID,
		Name:          sanitizedName,
		Slug:          utils.GenerateSlug(sanitizedName),
		Description:   sanitizedDesc,
		AvatarUrl:     req.AvatarUrl,
		CoverImageUrl: req.CoverImageUrl,
	}
	s.logger.Debug("params arg: ", arg)

	communities, err := s.store.CreateCommunity(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create community: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to create community",
		})
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
		s.logger.Error("Failed to bind community: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}
	s.logger.Debug("bind the request: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	community, err := s.store.GetCommunity(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get community: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get community",
		})
		return
	}

	ctx.JSON(http.StatusOK, community)
	return
}

// Get all communities by owner.
func (s *HandlersServer) GetAllCommunitiesFunc(ctx *gin.Context) {

	communities, err := s.store.GetAllCommunities(ctx)
	if err != nil {
		s.logger.Error("Failed to  get communities: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get communities",
		})
		return
	}
	s.logger.Debug("get all community: ", communities)
	ctx.JSON(http.StatusOK, communities)
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
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	communityPublicID, err := uuid.Parse(req.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	usersList, err := s.store.GetCommunitiesMember(ctx, communityPublicID)
	if err != nil {
		s.logger.Error("Failed to get community member: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get community member",
		})
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
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	usersList, err := s.store.GetCommunityByCommunityName(ctx, req.CommunitiesName)
	if err != nil {
		s.logger.Error("Failed to get community by community name: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get community by community name",
		})
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
	var reqJSON updateCommunityName

	err := ctx.ShouldBindUri(&reqUri)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	err = ctx.ShouldBindJSON(&reqJSON)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	publicID, err := uuid.Parse(reqUri.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	community, err := s.store.GetCommunity(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get community: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get community",
		})
		return
	}

	if authPayload.UserID != community.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed"})
		return
	}

	response, err := s.store.UpdateCommunityName(ctx, publicID, reqJSON.CommunityName)
	if err != nil {
		s.logger.Error("Failed to update the community name: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to update community name",
		})
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
	var reqJSON updateCommunityDescription

	err := ctx.ShouldBindUri(&reqUri)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	err = ctx.ShouldBindJSON(&reqJSON)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}

	publicID, err := uuid.Parse(reqUri.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	community, err := s.store.GetCommunity(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get community: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get community",
		})
		return
	}

	if authPayload.UserID != community.UserID {
		s.logger.Error("User is not allowed to update the description")
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"code":    "FORBIDDEN_ERROR",
			"message": "You are not allowed",
		})
		return
	}

	response, err := s.store.UpdateCommunityDescription(ctx, publicID, reqJSON.CommunityDescription)
	if err != nil {
		s.logger.Error("Failed to update the  description: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to update community description",
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
