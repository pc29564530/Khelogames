package handlers

import (
	"khelogames/core/token"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *HandlersServer) AddJoinCommunityFunc(ctx *gin.Context) {
	var req struct {
		CommunityPublicID string `uri:"community_public_id" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("Bind the request: ", req)

	communityPublicID, err := uuid.Parse(req.CommunityPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"community_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	community, err := s.store.GetCommunity(ctx, communityPublicID)
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

	if community.UserID != authPayload.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "You are not allowed to join this community",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	communityUser, err := s.txStore.AddJoinCommunityTx(ctx.Request.Context(), communityPublicID, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to join community: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to join community",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Debug("Successfully joined community: ", communityUser)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message": "Successfully joined",
			"member":  communityUser,
		},
	})
}

func (s *HandlersServer) GetCommunityByUserFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	communityList, err := s.store.GetCommunityByUser(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get community by user: ", err)
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
	s.logger.Debug("community by user: ", communityList)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    communityList,
	})
}
