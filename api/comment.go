package api

import (
	"database/sql"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createCommentRequest struct {
	CommentText string `json:"comment_text"`
}

type createCommentThreadIdRequest struct {
	ThreadID int64 `uri:"threadId"`
}

func (server *Server) createComment(ctx *gin.Context) {
	var req createCommentRequest
	var reqThreadId createCommentThreadIdRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			server.logger.Error("No row error: %v", err)
			return
		}
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = ctx.ShouldBindUri(&reqThreadId)
	if err != nil {
		if err == sql.ErrNoRows {
			server.logger.Error("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateCommentParams{
		ThreadID:    reqThreadId.ThreadID,
		Owner:       authPayload.Username,
		CommentText: req.CommentText,
	}

	comment, err := server.store.CreateComment(ctx, arg)
	if err != nil {
		server.logger.Error("Failed to create comment: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comment)
	return
}

type getAllCommentRequest struct {
	ThreadID int64 `uri:"thread_id"`
}

func (server *Server) getAllComment(ctx *gin.Context) {
	var req getAllCommentRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	comments, err := server.store.GetAllComment(ctx, req.ThreadID)
	if err != nil {
		server.logger.Error("Failed to get comment: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comments)
	return
}

type getCommentByUserRequest struct {
	Owner string `uri:"owner"`
}

func (server *Server) getCommentByUser(ctx *gin.Context) {
	var req getCommentByUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if req.Owner == "undefined" {
		server.logger.Error("Failed to get defined owner: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	comments, err := server.store.GetCommentByUser(ctx, req.Owner)
	if err != nil {
		server.logger.Error("Failed to get comment by user: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comments)
	return
}
