package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "khelogames/db/sqlc"
	"net/http"
)

type createThreadRequest struct {
	Username        string         `json:"username"`
	CommunitiesName sql.NullString `json:"communities_name"`
	Title           sql.NullString `json:"title"`
	Content         sql.NullString `json:"content"`
}

func (server *Server) createThread(ctx *gin.Context) {
	var req createThreadRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateThreadParams{
		Username: req.Username,
		Content:  req.Content,
	}

	thread, err := server.store.CreateThread(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, thread)
	return
}

func (server *Server) getAllThreads(ctx *gin.Context) {
	threads, err := server.store.GetAllThreads(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, threads)
	return
}

type getThreadsByCommunitiesRequest struct {
	CommunitiesName sql.NullString `json:"communities_name"`
}

func (server *Server) getAllThreadsByCommunities(ctx *gin.Context) {
	var req getThreadsByCommunitiesRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	//authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	threads, err := server.store.GetAllThreadsByCommunities(ctx, req.CommunitiesName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, threads)
	return
}
