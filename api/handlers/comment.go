package handlers

import (
	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createCommentRequest struct {
	CommentText string `json:"comment_text"`
}

type createCommentThreadIdRequest struct {
	ThreadPublicID uuid.UUID `uri:"thread_public_id"`
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

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	comment, err := s.store.CreateComment(ctx, uriReq.ThreadPublicID, authPayload.ID, bodyReq.CommentText)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create comment"})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

type getAllCommentRequest struct {
	PublicID uuid.UUID `uri:"public_id"`
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

	comments, err := s.store.GetAllComment(ctx, req.PublicID)
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
	PublicID uuid.UUID `json:"public_id"`
}

func (s *HandlersServer) DeleteCommentByUserFunc(ctx *gin.Context) {
	var req deleteCommentByUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	s.logger.Debug("bind the request: ", req)

	authPayload := ctx.MustGet(pkg.AuthorizationHeaderKey).(*token.Payload)

	comments, err := s.store.DeleteComment(ctx, req.PublicID, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get comment by user: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Info("successfully get comment by user")
	ctx.JSON(http.StatusOK, comments)
}
