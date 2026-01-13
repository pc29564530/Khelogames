package handlers

import (
	"khelogames/core/token"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createLikeRequest struct {
	ThreadPublicID string `uri:"thread_public_id" binding:"required"`
}

func (s *HandlersServer) CreateLikeFunc(ctx *gin.Context) {
	var req createLikeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("bind the request: ", req)

	threadPublicID, err := uuid.Parse(req.ThreadPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"thread_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	likeThread, err := s.store.CreateLike(ctx, authPayload.PublicID, threadPublicID)
	if err != nil {
		s.logger.Error("Failed to create like: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to like the thread",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("liked the thread: ", likeThread)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    likeThread,
	})
}

type countLikeRequest struct {
	ThreadPublicID string `uri:"thread_public_id" binding:"required"`
}

func (s *HandlersServer) CountLikeFunc(ctx *gin.Context) {
	var req countLikeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("bind the request: ", req)

	threadPublicID, err := uuid.Parse(req.ThreadPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"thread_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	countLike, err := s.store.CountLikeUser(ctx, threadPublicID)
	if err != nil {
		s.logger.Error("Failed to count like user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to count like users",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("get like count: ", countLike)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    countLike,
	})
}

type checkUserRequest struct {
	ThreadPublicID string `uri:"thread_public_id" binding:"required"`
}

func (s *HandlersServer) CheckLikeByUserFunc(ctx *gin.Context) {
	var req checkUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("bind the request: ", req)

	threadPublicID, err := uuid.Parse(req.ThreadPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"thread_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	userFound, err := s.store.CheckUserCount(ctx, authPayload.PublicID, threadPublicID)
	if err != nil {
		s.logger.Error("Failed to check user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to check like by user",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("liked by user: ", userFound)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userFound,
	})
}
