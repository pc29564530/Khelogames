package handlers

import (
	"khelogames/core/token"
	errorhandler "khelogames/error_handler"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createCommentRequest struct {
	CommentText string `json:"comment_text" binding:"required,min=1,max=5000"`
}

type createCommentThreadIdRequest struct {
	ThreadPublicID string `uri:"thread_public_id" binding:"required"`
}

func (s *HandlersServer) CreateCommentFunc(ctx *gin.Context) {
	var uriReq createCommentThreadIdRequest
	var bodyReq createCommentRequest

	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		s.logger.Error("Failed to bind thread public ID: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		s.logger.Error("Failed to bind comment body: ", err)
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	threadPublicID, err := uuid.Parse(uriReq.ThreadPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"thread_public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	comment, err := s.store.CreateComment(ctx, threadPublicID, authPayload.PublicID, bodyReq.CommentText)
	if err != nil {
		s.logger.Error("Failed to create comment: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Could not create comment",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	profile, err := s.store.GetProfileByUserID(ctx, comment.UserID)
	if err != nil {
		s.logger.Error("Failed to get profile by user ID: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Could not get user profile",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":                comment.ID,
			"public_id":         comment.PublicID,
			"thread_id":         comment.ThreadID,
			"user_id":           comment.UserID,
			"parent_comment_id": comment.ParentCommentID,
			"comment_text":      comment.CommentText,
			"like_count":        comment.LikeCount,
			"reply_count":       comment.ReplyCount,
			"is_deleted":        comment.IsDeleted,
			"is_edited":         comment.IsEdited,
			"created_at":        comment.CreatedAt,
			"updated_at":        comment.UpdatedAt,
			"profile":           profile,
		},
	})
}

type getAllCommentRequest struct {
	PublicID string `uri:"public_id" binding:"required"`
}

func (s *HandlersServer) GetAllCommentFunc(ctx *gin.Context) {
	var req getAllCommentRequest
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

	comments, err := s.store.GetAllComment(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get comment: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get comments",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("Successfully get all the comments : ", comments)
	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    comments,
	})
}

type deleteCommentByUserRequest struct {
	PublicID string `uri:"public_id" binding:"required"`
}

func (s *HandlersServer) DeleteCommentByUserFunc(ctx *gin.Context) {
	var req deleteCommentByUserRequest
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

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	comments, err := s.store.DeleteComment(ctx, publicID, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to delete comment: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to delete comment",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Info("successfully deleted comment")
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    comments,
	})
}
