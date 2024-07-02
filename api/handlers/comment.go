package handlers

import (
	"database/sql"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/pkg"
	"khelogames/token"
	"net/http"

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
			fmt.Errorf("No row error: %v", err)
			return
		}
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	err = ctx.ShouldBindUri(&reqThreadId)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateCommentParams{
		ThreadID:    reqThreadId.ThreadID,
		Owner:       authPayload.Username,
		CommentText: req.CommentText,
	}

	comment, err := s.store.CreateComment(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to create comment: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

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
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	comments, err := s.store.GetAllComment(ctx, req.ThreadID)
	if err != nil {
		fmt.Errorf("Failed to get comment: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, comments)
	return
}

type getCommentByUserRequest struct {
	Owner string `uri:"owner"`
}

func (s *CommentServer) GetCommentByUserFunc(ctx *gin.Context) {
	var req getCommentByUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	if req.Owner == "undefined" {
		fmt.Errorf("Failed to get defined owner: %v", err)
		ctx.JSON(http.StatusBadRequest, (err))
		return
	}

	comments, err := s.store.GetCommentByUser(ctx, req.Owner)
	if err != nil {
		fmt.Errorf("Failed to get comment by user: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, comments)
	return
}
