package handlers

import (
	"khelogames/core/token"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type createFollowingRequest struct {
	TargetPublicID string `uri:"target_public_id" binding:"required"`
}

// this is function i have to call the get_following endpoint so that using that i can verify the following list
func (s *HandlersServer) CreateUserConnectionsFunc(ctx *gin.Context) {
	var req createFollowingRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	publicID, err := uuid.Parse(req.TargetPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"target_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	follower, err := s.store.CreateUserConnections(ctx, authPayload.PublicID, publicID)
	if err != nil {
		s.logger.Error("Failed to create following: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to create following",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("Successfully created: ", follower)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    follower,
	})
}

func (s *HandlersServer) GetAllFollowerFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	follower, err := s.store.GetAllFollower(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get follower: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get follower",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("Successfully retrieved follower: ", follower)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    follower,
	})
}

func (s *HandlersServer) GetAllFollowingFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	following, err := s.store.GetAllFollowing(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get following: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get following",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("Successfully retrieved following: ", following)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    following,
	})
}

func (s *HandlersServer) GetFollowerCountFunc(ctx *gin.Context) {

	var req struct {
		PublicID string `uri:"public_id"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	follower, err := s.store.GetFollowerCount(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get follower: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get follower",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("Successfully retrieved follower: ", follower)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    follower,
	})
}

func (s *HandlersServer) GetFollowingCountFunc(ctx *gin.Context) {

	var req struct {
		PublicID string `uri:"public_id"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	following, err := s.store.GetFollowingCount(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get following: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get following",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("Successfully retrieved following: ", following)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    following,
	})
}

type deleteFollowingRequest struct {
	TargetPublicID string `uri:"target_public_id" binding:"required"`
}

func (s *HandlersServer) DeleteFollowingFunc(ctx *gin.Context) {
	var req deleteFollowingRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	targetPublicID, err := uuid.Parse(req.TargetPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format: ", err)
		fieldErrors := map[string]string{"target_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	err = s.store.DeleteUsersConnections(ctx, authPayload.PublicID, targetPublicID)
	if err != nil {
		s.logger.Error("Failed to unfollow user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to unfollow user",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message": "Successfully unfollowed user",
		},
	})
}

func (s *HandlersServer) IsFollowingFunc(ctx *gin.Context) {
	var req struct {
		TargetPublicID string `uri:"target_public_id" binding:"required"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	targetPublicID, err := uuid.Parse(req.TargetPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"target_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	isFollowing, err := s.store.IsFollowingF(ctx, authPayload.PublicID, targetPublicID)
	if err != nil {
		s.logger.Error("Failed to check following ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to check following status",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Info("Successfully checked following status: ", isFollowing)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"is_following": isFollowing,
		},
	})
}
