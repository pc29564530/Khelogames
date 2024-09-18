package handlers

import (
	"database/sql"
	db "khelogames/db/sqlc"

	"khelogames/pkg"
	"khelogames/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type createCommentRequest struct {
	CommentText string `json:"comment_text"`
}

type createCommentThreadIdRequest struct {
	ThreadID int64 `uri:"threadId"`
}

func (s *HandlersServer) CreateCommentFunc(ctx *gin.Context) {
	var req createCommentRequest
	var reqThreadId createCommentThreadIdRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row error: ", err)
			return
		}
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("successfully bind: ", req)
	err = ctx.ShouldBindUri(&reqThreadId)
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
	s.logger.Debug("successfully bind: ", reqThreadId)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateCommentParams{
		ThreadID:    reqThreadId.ThreadID,
		Owner:       authPayload.Username,
		CommentText: req.CommentText,
	}
	s.logger.Debug("params arg: ", arg)

	comment, err := s.store.CreateComment(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create comment: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Info("successfully create the comment")

	ctx.JSON(http.StatusOK, comment)
}

type getAllCommentRequest struct {
	ThreadID int64 `uri:"thread_id"`
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

	comments, err := s.store.GetAllComment(ctx, req.ThreadID)
	if err != nil {
		s.logger.Error("Failed to get comment: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get all the comments : ", comments)
	s.logger.Debug("Received threads from database")
	var commentsDetails []map[string]interface{}

	for _, comment := range comments {
		profile, err := s.store.GetProfile(ctx, comment.Owner)
		if err != nil {
			s.logger.Error("Failed to find the profile ", err)
			return
		}
		var displayText string
		if profile.AvatarUrl == "" {
			displayText = strings.ToUpper(string(profile.FullName[0]))
		}

		commentDetail := map[string]interface{}{
			"id":           comment.ID,
			"username":     comment.Owner,
			"comment":      comment.CommentText,
			"display_text": displayText,
			"full_name":    profile.FullName,
			"avatar_url":   profile.AvatarUrl,
			"created_at":   comment.CreatedAt,
		}
		commentsDetails = append(commentsDetails, commentDetail)
	}
	s.logger.Info("successfully get all comment details")
	ctx.JSON(http.StatusOK, commentsDetails)
}

type getCommentByUserRequest struct {
	Owner string `uri:"owner"`
}

func (s *HandlersServer) GetCommentByUserFunc(ctx *gin.Context) {
	var req getCommentByUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)
	if req.Owner == "undefined" {
		s.logger.Error("Failed to get defined owner: ", err)
		ctx.JSON(http.StatusBadRequest, (err))
		return
	}

	comments, err := s.store.GetCommentByUser(ctx, req.Owner)
	if err != nil {
		s.logger.Error("Failed to get comment by user: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Info("successfully get comment by user")
	ctx.JSON(http.StatusOK, comments)
}

type deleteCommentByUserRequest struct {
	ID    int64  `json:"id"`
	Owner string `uri:"owner"`
}

func (s *HandlersServer) DeleteCommentByUserFunc(ctx *gin.Context) {
	var req deleteCommentByUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		return
	}
	s.logger.Debug("bind the request: ", req)
	if req.Owner == "undefined" {
		s.logger.Error("Failed to get defined owner: ", err)
		ctx.JSON(http.StatusBadRequest, (err))
		return
	}

	arg := db.DeleteCommentParams{
		ID:    req.ID,
		Owner: req.Owner,
	}

	comments, err := s.store.DeleteComment(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get comment by user: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Info("successfully get comment by user")
	ctx.JSON(http.StatusOK, comments)
}
