package handlers

import (
	"database/sql"

	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type createFollowingRequest struct {
	TargetPublicID uuid.UUID `uri:"target_public_id"`
}

// this is function i have to call the get_following endpoint so that using that i can verify the following list
func (s *HandlersServer) CreateUserConnectionsFunc(ctx *gin.Context) {
	var req createFollowingRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row error: ", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	follower, err := s.store.CreateUserConnections(ctx, authPayload.PublicID, req.TargetPublicID)
	if err != nil {
		s.logger.Error("Failed to create following: ", err)
		ctx.JSON(http.StatusBadRequest, (err))
		return
	}
	s.logger.Debug("successfully created: ", follower)
	ctx.JSON(http.StatusOK, follower)
	return

}

func (s *HandlersServer) GetAllFollowerFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	follower, err := s.store.GetAllFollower(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get follwer: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("successfully get ", follower)
	ctx.JSON(http.StatusOK, follower)
	return
}

func (s *HandlersServer) GetAllFollowingFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	follower, err := s.store.GetAllFollowing(ctx, authPayload.PublicID)
	if err != nil {
		s.logger.Error("Failed to get following: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("successfully get: ", follower)

	ctx.JSON(http.StatusOK, follower)
	return
}

type deleteFollowingRequest struct {
	TargetPublicID uuid.UUID `uri:"target_public_id"`
}

func (s *HandlersServer) DeleteFollowingFunc(ctx *gin.Context) {

	var req deleteFollowingRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row error: ", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)

	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	err = s.store.DeleteUsersConnections(ctx, authPayload.PublicID, req.TargetPublicID)
	if err != nil {
		s.logger.Error("Failed to unfollow user: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"deleted": "Unfollow the user"})
	return
}

func (s *HandlersServer) IsFollowingFunc(ctx *gin.Context) {
	var req struct {
		TargetPublicID uuid.UUID `uri:"target_public_id"`
	}
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	isFollowing, err := s.store.IsFollowingF(ctx, authPayload.PublicID, req.TargetPublicID)
	if err != nil {
		s.logger.Error("Failed to check following ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check following status"})
		return
	}

	s.logger.Info("Successfully checked following status")
	ctx.JSON(http.StatusOK, gin.H{"is_following": isFollowing})
}
