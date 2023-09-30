package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	db "khelogames/db/sqlc"
	"khelogames/util"
	"net/http"
	"time"
)

type createUserRequest struct {
	Username       string `json:"username"`
	MobileNumber   string `json:"mobileNumber"`
	HashedPassword string `json:"password"`
}

type createUserResponse struct {
	Username     string    `json:"username"`
	MobileNumber string    `json:"mobileNumber"`
	CreatedAt    time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {

	var req createUserRequest

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
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

	user, err := server.store.CreateUser(ctx, arg)
	fmt.Println("unable to create a new user: ", err)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	resp := createUserResponse{
		Username:     user.Username,
		MobileNumber: user.MobileNumber,
		CreatedAt:    user.CreatedAt,
	}
	fmt.Println(resp)
	ctx.JSON(http.StatusOK, resp)

	deleteSignUp, err := server.store.DeleteSignup(ctx, req.MobileNumber)
	if err != nil {
		fmt.Errorf("unable to delete the mobile number details: ", err)
		return
	}
	ctx.JSON(http.StatusOK, deleteSignUp)
	return
}

type getUserRequest struct {
	Username string `uri:"username"`
}

func (server *Server) getUsers(ctx *gin.Context) {
	var req getUserRequest
	err := ctx.ShouldBindUri(&req)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println(req.Username)

	users, err := server.store.GetUser(ctx, req.Username)
	fmt.Println(err)
	fmt.Println(req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Println(users)

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
