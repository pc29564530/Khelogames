package handlers

import (
	"fmt"
	"khelogames/core/token"
	db "khelogames/database"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"net/http"

	utils "khelogames/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createCommunitiesRequest struct {
	CommunityName string `json:"name" binding:"required,min=3,max=100"`
	Description   string `json:"description"`
	AvatarUrl     string `json:"avatar_url" binding:"omitempty,url"`
	CoverImageUrl string `json:"cover_image_url" binding:"omitempty,url"`
}

// Create communities function
func (s *HandlersServer) CreateCommunitesFunc(ctx *gin.Context) {
	var req createCommunitiesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		fmt.Println("Field Error: ", fieldErrors)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
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
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create community",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("created community : ", communities)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    communities,
	})
}

type getCommunityRequest struct {
	PublicID string `uri:"public_id" binding:"required,min=1"`
}

// get Community by id.
func (s *HandlersServer) GetCommunityFunc(ctx *gin.Context) {
	var req getCommunityRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("bind the request: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	community, err := s.store.GetCommunity(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get community: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get community",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    community,
	})
}

// Get all communities by owner.
func (s *HandlersServer) GetAllCommunitiesFunc(ctx *gin.Context) {

	communities, err := s.store.GetAllCommunities(ctx)
	if err != nil {
		s.logger.Error("Failed to  get communities: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get communities",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("get all community: ", communities)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    communities,
	})
}

// Get all users that have joined a particular communities
type getCommunitiesMemberRequest struct {
	CommunityPublicID string `uri:"community_public_id" binding:"required"`
}

func (s *HandlersServer) GetCommunitiesMemberFunc(ctx *gin.Context) {
	var req getCommunitiesMemberRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	communityPublicID, err := uuid.Parse(req.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"community_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	usersList, err := s.store.GetCommunitiesMember(ctx, communityPublicID)
	if err != nil {
		s.logger.Error("Failed to get community member: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get community member",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("get community member: ", usersList)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    usersList,
	})
}

type getCommunityByCommunityNameRequest struct {
	CommunitiesName string `uri:"communities_name" binding:"required"`
}

func (s *HandlersServer) GetCommunityByCommunityNameFunc(ctx *gin.Context) {
	var req getCommunityByCommunityNameRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	usersList, err := s.store.GetCommunityByCommunityName(ctx, req.CommunitiesName)
	if err != nil {
		s.logger.Error("Failed to get community by community name: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get community by community name",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("get community by commmunity name: ", usersList)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    usersList,
	})
}

type updateCommunityName struct {
	CommunityName string `json:"community_name" binding:"required,min=3,max=100"`
}

func (s *HandlersServer) UpdateCommunityByCommunityNameFunc(ctx *gin.Context) {
	var reqUri communityPublicIDReq
	var reqJSON updateCommunityName

	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	publicID, err := uuid.Parse(reqUri.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"community_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	community, err := s.store.GetCommunity(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get community: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get community",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	if authPayload.UserID != community.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You are not allowed",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	response, err := s.store.UpdateCommunityName(ctx, publicID, reqJSON.CommunityName)
	if err != nil {
		s.logger.Error("Failed to update the community name: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update community name",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

type communityPublicIDReq struct {
	CommunityPublicID string `uri:"community_public_id" binding:"required"`
}

type updateCommunityDescription struct {
	CommunityDescription string `json:"community_description" binding:"required,min=1,max=500"`
}

func (s *HandlersServer) UpdateCommunityByDescriptionFunc(ctx *gin.Context) {
	var reqUri communityPublicIDReq
	var reqJSON updateCommunityDescription

	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	if err := ctx.ShouldBindJSON(&reqJSON); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	publicID, err := uuid.Parse(reqUri.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"community_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	community, err := s.store.GetCommunity(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get community: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get community",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	if authPayload.UserID != community.UserID {
		s.logger.Error("User is not allowed to update the description")
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You are not allowed",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	response, err := s.store.UpdateCommunityDescription(ctx, publicID, reqJSON.CommunityDescription)
	if err != nil {
		s.logger.Error("Failed to update the  description: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to update community description",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
