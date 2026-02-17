package handlers

import (
	"khelogames/core/token"
	"khelogames/database/models"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type likeThreadRequest struct {
	ThreadPublicID string `uri:"thread_public_id" binding:"required"`
}

// LikeThreadFunc toggles a like on a thread in a single API call.
// Check → like/unlike in one transaction → returns the updated thread.
func (s *HandlersServer) LikeThreadFunc(ctx *gin.Context) {
	var req likeThreadRequest
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

	// Check if user already liked this thread
	userLikeCount, err := s.store.CheckUserCount(ctx, authPayload.PublicID, threadPublicID)
	if err != nil {
		s.logger.Error("Failed to check user like: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to check like status",
			},
		})
		return
	}

	var thread *models.Thread

	if userLikeCount != nil && *userLikeCount > 0 {
		// Already liked → unlike
		thread, err = s.txStore.DeleteLikeThreadTx(ctx, authPayload.PublicID, threadPublicID)
		if err != nil {
			s.logger.Error("Failed to unlike thread: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to unlike thread",
				},
			})
			return
		}
	} else {
		// Not liked → like
		thread, err = s.txStore.CreateLikeTx(ctx, authPayload.PublicID, threadPublicID)
		if err != nil {
			s.logger.Error("Failed to like thread: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to like thread",
				},
			})
			return
		}
	}

	s.logger.Debugf("Thread %s like toggled, new count=%d", threadPublicID, thread.LikeCount)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    thread,
	})
}
