package handlers

import (
	db "khelogames/db/sqlc"

	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createLikeRequest struct {
	ThreadID int64 `uri:"thread_id"`
}

func (s *HandlersServer) CreateLikeFunc(ctx *gin.Context) {
	var req createLikeRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind : ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateLikeParams{
		ThreadID: req.ThreadID,
		Username: authPayload.Username,
	}
	s.logger.Debug("params arg: ", arg)

	likeThread, err := s.store.CreateLike(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create like : ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("liked the thread: ", likeThread)

	ctx.JSON(http.StatusOK, likeThread)
	return
}

type countLikeRequest struct {
	ThreadID int64 `uri:"thread_id"`
}

func (s *HandlersServer) CountLikeFunc(ctx *gin.Context) {
	var req countLikeRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)

	countLike, err := s.store.CountLikeUser(ctx, req.ThreadID)
	if err != nil {
		s.logger.Error("Failed to count like user: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get like count: ", countLike)
	ctx.JSON(http.StatusOK, countLike)
	return
}

type checkUserRequest struct {
	ThreadID int64 `uri:"thread_id"`
}

func (s *HandlersServer) CheckLikeByUserFunc(ctx *gin.Context) {
	var req checkUserRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.CheckUserCountParams{
		ThreadID: req.ThreadID,
		Username: authPayload.Username,
	}
	s.logger.Debug("params arg: ", arg)

	userFound, err := s.store.CheckUserCount(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to check user: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("liked by user: ", userFound)
	ctx.JSON(http.StatusOK, userFound)
	return
}
