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
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = ctx.ShouldBindUri(&reqThreadId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	comments, err := server.store.GetAllComment(ctx, req.ThreadID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, comments)
	return
}
