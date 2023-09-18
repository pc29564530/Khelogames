package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"
)

type createFollowingRequest struct {
	FollowingOwner string `uri:"following_owner"`
}

func (server *Server) createFollowing(ctx *gin.Context) {
	var req createFollowingRequest
	fmt.Println("line no 18")
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println(req.FollowingOwner)
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateFollowingParams{
		FollowerOwner:  authPayload.Username,
		FollowingOwner: req.FollowingOwner,
	}
	fmt.Println(arg.FollowerOwner)
	fmt.Println(arg.FollowingOwner)
	follower, err := server.store.CreateFollowing(ctx, arg)
	fmt.Println(err)
	fmt.Println("line 39")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fmt.Println(follower)
	ctx.JSON(http.StatusOK, follower)
	return

}

func (server *Server) getAllFollower(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	fmt.Println(authPayload.Username)
	follower, err := server.store.GetAllFollower(ctx, authPayload.Username)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, follower)
	return
}

func (server *Server) getAllFollowing(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	fmt.Println(authPayload.Username)
	follower, err := server.store.GetAllFollowing(ctx, authPayload.Username)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println(follower)
	ctx.JSON(http.StatusOK, follower)
	return
}

func (server *Server) deleteFollowing(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	following, err := server.store.DeleteFollowing(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, following)
	return
}
