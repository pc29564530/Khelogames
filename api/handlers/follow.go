package handlers

import (
	"database/sql"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/pkg"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type FollowServer struct {
	store  *db.Store
	logger *logger.Logger
}

func NewFollowServer(store *db.Store, logger *logger.Logger) *FollowServer {
	return &FollowServer{store: store, logger: logger}
}

type createFollowingRequest struct {
	FollowingOwner string `uri:"following_owner"`
}

// this is function i have to call the get_following endpoint so that using that i can verify the following list
func (s *FollowServer) CreateFollowingFunc(ctx *gin.Context) {
	var req createFollowingRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)

	arg := db.CreateFollowingParams{
		FollowerOwner:  authPayload.Username,
		FollowingOwner: req.FollowingOwner,
	}

	follower, err := s.store.CreateFollowing(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to create following: %v", err)
		ctx.JSON(http.StatusBadRequest, (err))
		return
	}
	ctx.JSON(http.StatusOK, follower)
	return

}

func (s *FollowServer) GetAllFollowerFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	follower, err := s.store.GetAllFollower(ctx, authPayload.Username)
	if err != nil {
		fmt.Errorf("Failed to get follwer: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	ctx.JSON(http.StatusOK, follower)
	return
}

func (s *FollowServer) GetAllFollowingFunc(ctx *gin.Context) {
	authPayload := ctx.MustGet(pkg.AuthorizationPayloadKey).(*token.Payload)
	follower, err := s.store.GetAllFollowing(ctx, authPayload.Username)
	if err != nil {
		fmt.Errorf("Failed to get following: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, follower)
	return
}

type deleteFollowingRequest struct {
	FollowingOwner string `uri:"following_owner"`
}

func (s *FollowServer) DeleteFollowingFunc(ctx *gin.Context) {

	var req deleteFollowingRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("No row error: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	following, err := s.store.DeleteFollowing(ctx, req.FollowingOwner)
	if err != nil {
		fmt.Errorf("Failed to unfollow user: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, following)
	return
}
