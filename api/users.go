package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	db "khelogames/db/sqlc"
	"khelogames/util"
	"net/http"
	"time"
)

type createUserRequest struct {
	Username       string `json:"username"`
	MobileNumber   string `json:"mobile_number"`
	HashedPassword string `json:"hashed_password"`
}

type createUserResponse struct {
	Username     string    `json:"username"`
	MobileNumber string    `json:"mobile_number"`
	CreatedAt    time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {

	var req createUserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadGateway, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		MobileNumber:   req.MobileNumber,
		HashedPassword: hashedPassword,
	}

	arg = db.CreateUserParams{
		Username:       req.Username,
		MobileNumber:   req.MobileNumber,
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	resp := createUserResponse{
		Username:     user.Username,
		MobileNumber: user.MobileNumber,
		CreatedAt:    user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, resp)
	return
}

type getUsersRequest struct {
	Username string `josn:"username"`
}

func (server *Server) getUsers(ctx *gin.Context) {
	var req getUsersRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	users, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
	return
}

type getListUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	var req getListUsersRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.ListUserParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	userList, err := server.store.ListUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, userList)
}
