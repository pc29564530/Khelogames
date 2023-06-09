package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	db "khelogames/db/sqlc"
	"khelogames/token"
	"net/http"
)

type getUsername struct {
	RecieverUsername string `json:"reciever_username"`
}

type createConnectionsRequest struct {
	SenderUsername   string `json:"sender_username"`
	RecieverUsername string `json:"reciever_username"`
	Status           string `json:"status"`
}

type createConnectionResponse struct {
	SenderUsername   string `json:"sender_username"`
	RecieverUsername string `json:"reciever_username"`
	Status           string `json:"status"`
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
	//err = sendFriendRequest(ctx, sender_username, user.Username)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//}

	argConnection := db.CreateConnectionsParams{
		RecieverUsername: user.Username,
		SenderUsername:   sender_username,
		Status:           "pending",
	}
	fmt.Println("line67")
	fmt.Println(argConnection.SenderUsername)
	fmt.Println(argConnection.RecieverUsername)
	fmt.Println(argConnection.Status)
	fmt.Println("line69")
	connection, err := server.store.CreateConnections(ctx, argConnection)
	fmt.Println("ram ram ji")
	if err != nil {
		fmt.Errorf("unable to create a connection ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, connection)
	fmt.Println("Successfully send the friend request")
	ctx.JSON(http.StatusOK, user)
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	username := authPayload.Username
	connections, err := server.store.ListConnections(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, connections)
}
