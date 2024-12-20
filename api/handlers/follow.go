package handlers

import (
	"database/sql"
	db "khelogames/database"

	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type createFollowingRequest struct {
	FollowingOwner string `uri:"following_owner"`
}

// this is function i have to call the get_following endpoint so that using that i can verify the following list
func (s *HandlersServer) CreateFollowingFunc(ctx *gin.Context) {
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

	arg := db.CreateFollowingParams{
		FollowerOwner:  authPayload.Username,
		FollowingOwner: req.FollowingOwner,
	}
	s.logger.Debug("params arg: ", arg)

	follower, err := s.store.CreateFollowing(ctx, arg)
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
	follower, err := s.store.GetAllFollower(ctx, authPayload.Username)
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
	follower, err := s.store.GetAllFollowing(ctx, authPayload.Username)
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
	FollowingOwner string `uri:"following_owner"`
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

	following, err := s.store.DeleteFollowing(ctx, authPayload.Username, req.FollowingOwner)
	if err != nil {
		s.logger.Error("Failed to unfollow user: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("successfully get: ", following)

	ctx.JSON(http.StatusOK, following)
	return
}

func (s *HandlersServer) CheckConnectionFunc(ctx *gin.Context) {

	followingOwner := ctx.Query("following_owner")
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	connectionEstablished, err := s.store.CheckConnection(ctx, authPayload.Username, followingOwner)
	if err != nil {
		s.logger.Error("Failed to check connection ", err)
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	s.logger.Info("Successfully checked connection ")
	ctx.JSON(http.StatusAccepted, connectionEstablished)
}

func (s *HandlersServer) IsFollowingFunc(ctx *gin.Context) {
	followingOnwer := ctx.Query("following_owner")
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	isFollowing, err := s.store.IsFollowingF(ctx, authPayload.Username, followingOnwer)
	if err != nil {
		s.logger.Error("Failed to check following ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check following status"})
		return
	}

	s.logger.Info("Successfully checked following status")
	ctx.JSON(http.StatusOK, gin.H{"is_following": isFollowing})
}
