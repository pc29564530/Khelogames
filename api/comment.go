package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"
)

type createCommentRequest struct {
	ThreadID    int64  `json: "thread_id"`
	CommentText string `json: "comment_text"`
}

func (server *Server) createComment(ctx *gin.Context) {
	var req createCommentRequest
	err := ctx.ShouldBindJSON(&req)
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
		ThreadID:    req.ThreadID,
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
