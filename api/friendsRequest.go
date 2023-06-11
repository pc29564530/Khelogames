package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"
	"time"
)

type getUsername struct {
	RecieverUsername string `json:"reciever_username"`
}

func (server *Server) getRecieverUsername(ctx *gin.Context) {
	var req getUsername

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.RecieverUsername)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	sender_username := authPayload.Username

	argConnection := db.CreateFriendsRequestParams{
		SenderUsername:   sender_username,
		RecieverUsername: user.Username,
		Status:           "pending",
	}
	connection, err := server.store.CreateFriendsRequest(ctx, argConnection)
	if err != nil {
		fmt.Errorf("unable to create a connection ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, connection)
	fmt.Println("Successfully send the friend request")
	return

}

type ListConnections struct {
	Username string `json:"username"`
}

func (server *Server) ListConnections(ctx *gin.Context) {
	var req ListConnections
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	connections, err := server.store.ListFriends(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, connections)
}

type getSenderUsernameRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getSenderUsernameResponse struct {
	ID               int64     `uri:"id" binding:"required min=1"`
	SenderUsername   string    `json:"sender_username"`
	RecieverUsername string    `json:"reciever_username"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

type acceptFriendRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) acceptFriend(ctx *gin.Context) {
	var req acceptFriendRequest
	err := ctx.ShouldBindUri(&req)
	fmt.Println("check the err")
	if err != nil {
		fmt.Println("no row found")
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	currentUser := authPayload.Username

	err = server.store.UpdateFriendsRequest(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetFriendsRequest(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

	}

	rsp := getSenderUsernameResponse{
		ID:               user.ID,
		SenderUsername:   user.SenderUsername,
		RecieverUsername: currentUser,
		Status:           user.Status,
		CreatedAt:        user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}
