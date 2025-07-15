package handlers

import (
	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createLikeRequest struct {
	ThreadPublicID uuid.UUID `uri:"thread_public_id"`
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

	likeThread, err := s.store.CreateLike(ctx, authPayload.PublicID, req.ThreadPublicID)
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
	ThreadPublicID uuid.UUID `uri:"thread_public_id"`
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

	countLike, err := s.store.CountLikeUser(ctx, req.ThreadPublicID)
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
	ThreadPublicID uuid.UUID `uri:"thread_public_id"`
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

	userFound, err := s.store.CheckUserCount(ctx, authPayload.PublicID, req.ThreadPublicID)
	if err != nil {
		s.logger.Error("Failed to check user: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("liked by user: ", userFound)
	ctx.JSON(http.StatusOK, userFound)
	return
}
