package handlers

import (
	"khelogames/core/token"
	"khelogames/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createCommentRequest struct {
	CommentText string `json:"comment_text"`
}

type createCommentThreadIdRequest struct {
	ThreadPublicID string `uri:"thread_public_id"`
}

func (s *HandlersServer) CreateCommentFunc(ctx *gin.Context) {
	var uriReq createCommentThreadIdRequest
	var bodyReq createCommentRequest

	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid thread ID"})
		return
	}

	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment body"})
		return
	}

	threadPublicID, err := uuid.Parse(uriReq.ThreadPublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	comment, err := s.store.CreateComment(ctx, threadPublicID, authPayload.PublicID, bodyReq.CommentText)
	if err != nil {
		s.logger.Error("Failed to create comment: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create comment"})
		return
	}
	profile, err := s.store.GetProfileByUserID(ctx, comment.UserID)
	if err != nil {
		s.logger.Error("Failed to get profile by user ID: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get profile"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
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
	})
}

type getAllCommentRequest struct {
	PublicID string `uri:"public_id"`
}

func (s *HandlersServer) GetAllCommentFunc(ctx *gin.Context) {
	var req getAllCommentRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
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

	comments, err := s.store.GetAllComment(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get comment: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get all the comments : ", comments)
	s.logger.Debug("Received threads from database")
	ctx.JSON(http.StatusAccepted, comments)
}

type deleteCommentByUserRequest struct {
	PublicID string `json:"public_id"`
}

func (s *HandlersServer) DeleteCommentByUserFunc(ctx *gin.Context) {
	var req deleteCommentByUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	s.logger.Debug("bind the request: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	comments, err := s.store.DeleteComment(ctx, publicID, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get comment by user: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Info("successfully get comment by user")
	ctx.JSON(http.StatusOK, comments)
}
