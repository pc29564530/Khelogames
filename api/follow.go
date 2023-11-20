package api

import (
	"database/sql"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type createFollowingRequest struct {
	FollowingOwner string `uri:"following_owner"`
}

// this is function i have to call the get_following endpoint so that using that i can verify the following list
func (server *Server) createFollowing(ctx *gin.Context) {
	var req createFollowingRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateFollowingParams{
		FollowerOwner:  authPayload.Username,
		FollowingOwner: req.FollowingOwner,
	}

	follower, err := server.store.CreateFollowing(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, follower)
	return

}

func (server *Server) getAllFollower(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	follower, err := server.store.GetAllFollower(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, follower)
	return
}

func (server *Server) getAllFollowing(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	follower, err := server.store.GetAllFollowing(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, follower)
	return
}

type deleteFollowingRequest struct {
	FollowingOwner string `uri:"following_owner"`
}

func (server *Server) deleteFollowing(ctx *gin.Context) {

	var req deleteFollowingRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	following, err := server.store.DeleteFollowing(ctx, req.FollowingOwner)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, following)
	return
}
