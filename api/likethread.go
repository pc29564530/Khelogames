package api

import (
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createLikeRequest struct {
	ThreadID int64 `uri:"thread_id"`
}

func (server *Server) createLike(ctx *gin.Context) {
	var req createLikeRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateLikeParams{
		ThreadID: req.ThreadID,
		Username: authPayload.Username,
	}

	likeThread, err := server.store.CreateLike(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, likeThread)
	return
}

type countLikeRequest struct {
	ThreadID int64 `uri:"thread_id"`
}

func (server *Server) countLike(ctx *gin.Context) {
	var req countLikeRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	countLike, err := server.store.CountLikeUser(ctx, req.ThreadID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, countLike)
	return
}

type checkUserRequest struct {
	ThreadID int64 `uri:"thread_id"`
}

func (server *Server) checkLikeByUser(ctx *gin.Context) {
	var req checkUserRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CheckUserCountParams{
		ThreadID: req.ThreadID,
		Username: authPayload.Username,
	}

	userFound, err := server.store.CheckUserCount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userFound)
	return
}
