package handlers

import (
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LikethreadServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewLikeThreadServer(store *db.Store, logger *logger.Logger) *LikethreadServer {
	return &LikethreadServer{store: store, logger: logger}
}

type createLikeRequest struct {
	ThreadID int64 `uri:"thread_id"`
}

func (s *LikethreadServer) CreateLikeFunc(ctx *gin.Context) {
	var req createLikeRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind : %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateLikeParams{
		ThreadID: req.ThreadID,
		Username: authPayload.Username,
	}
	s.logger.Debug("params arg: ", arg)

	likeThread, err := s.store.CreateLike(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to create like : %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("liked the thread: %v", likeThread)

	ctx.JSON(http.StatusOK, likeThread)
	return
}

type countLikeRequest struct {
	ThreadID int64 `uri:"thread_id"`
}

func (s *LikethreadServer) CountLikeFunc(ctx *gin.Context) {
	var req countLikeRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)

	countLike, err := s.store.CountLikeUser(ctx, req.ThreadID)
	if err != nil {
		s.logger.Error("Failed to count like user: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get like count: %v", countLike)
	ctx.JSON(http.StatusOK, countLike)
	return
}

type checkUserRequest struct {
	ThreadID int64 `uri:"thread_id"`
}

func (s *LikethreadServer) CheckLikeByUserFunc(ctx *gin.Context) {
	var req checkUserRequest

	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.CheckUserCountParams{
		ThreadID: req.ThreadID,
		Username: authPayload.Username,
	}
	s.logger.Debug("params arg: %v", arg)

	userFound, err := s.store.CheckUserCount(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to check user: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("liked by user: %v", userFound)
	ctx.JSON(http.StatusOK, userFound)
	return
}
