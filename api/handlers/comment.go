package handlers

import (
	"database/sql"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/pkg"
	"khelogames/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CommentServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewCommentServer(store *db.Store, logger *logger.Logger) *CommentServer {
	return &CommentServer{store: store, logger: logger}
}

type createCommentRequest struct {
	CommentText string `json:"comment_text"`
}

type createCommentThreadIdRequest struct {
	ThreadID int64 `uri:"threadId"`
}

func (s *CommentServer) CreateCommentFunc(ctx *gin.Context) {
	var req createCommentRequest
	var reqThreadId createCommentThreadIdRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row error: %v", err)
			return
		}
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("successfully bind: %v", req)
	err = ctx.ShouldBindUri(&reqThreadId)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("successfully bind: %v", reqThreadId)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateCommentParams{
		ThreadID:    reqThreadId.ThreadID,
		Owner:       authPayload.Username,
		CommentText: req.CommentText,
	}
	s.logger.Debug("params arg: %v", arg)

	comment, err := s.store.CreateComment(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create comment: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Info("successfully create the comment")

	ctx.JSON(http.StatusOK, comment)
	return
}

type getAllCommentRequest struct {
	ThreadID int64 `uri:"thread_id"`
}

func (s *CommentServer) GetAllCommentFunc(ctx *gin.Context) {
	var req getAllCommentRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Debug("bind the request: %v", req)

	comments, err := s.store.GetAllComment(ctx, req.ThreadID)
	if err != nil {
		s.logger.Error("Failed to get comment: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get all the comments :%v ", comments)
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
	return
}

type getCommentByUserRequest struct {
	Owner string `uri:"owner"`
}

func (s *CommentServer) GetCommentByUserFunc(ctx *gin.Context) {
	var req getCommentByUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)
	if req.Owner == "undefined" {
		s.logger.Error("Failed to get defined owner: %v", err)
		ctx.JSON(http.StatusBadRequest, (err))
		return
	}

	comments, err := s.store.GetCommentByUser(ctx, req.Owner)
	if err != nil {
		s.logger.Error("Failed to get comment by user: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Info("successfully get comment by user")
	ctx.JSON(http.StatusOK, comments)
	return
}
