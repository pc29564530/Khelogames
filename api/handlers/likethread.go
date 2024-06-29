package handlers

import (
	"fmt"
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
		fmt.Errorf("Failed to bind : %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	arg := db.CreateLikeParams{
		ThreadID: req.ThreadID,
		Username: authPayload.Username,
	}

	likeThread, err := s.store.CreateLike(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to create like : %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

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
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	countLike, err := s.store.CountLikeUser(ctx, req.ThreadID)
	if err != nil {
		fmt.Errorf("Failed to count like user: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

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
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.CheckUserCountParams{
		ThreadID: req.ThreadID,
		Username: authPayload.Username,
	}

	userFound, err := s.store.CheckUserCount(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to check user count: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, userFound)
	return
}
