package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}
type createUserRequest struct {
	Username     string `json:"username"`
	MobileNumber string `json:"mobile_number"`
	Role         string `json:"role"`
	Gmail        string `json:"gmail"`
}

type userResponse struct {
	Username     string `json:"username"`
	MobileNumber string `json:"mobile_number"`
	Role         string `json:"role"`
	Gmail        string `json:"gmail"`
}

type getUserRequest struct {
	PublicID string `uri:"public_id"`
}

func (s *HandlersServer) GetUsersFunc(ctx *gin.Context) {
	var req getUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind user: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid request format",
		})
		return
	}
	s.logger.Debug("bind the reqeust: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"code":    "VALIDATION_ERROR",
			"message": "Invalid UUID format",
		})
		return
	}

	users, err := s.store.GetUser(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get user",
		})
		return
	}
	s.logger.Debug("get the user data: ", users)
	ctx.JSON(http.StatusOK, users)
	return
}

type getListUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *HandlersServer) ListUsersFunc(ctx *gin.Context) {
	var req getListUsersRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Failed to bind the request",
		})
		return
	}
	s.logger.Debug("bind the request: ", req)
	Offset := (req.PageID - 1) * req.PageSize

	userList, err := s.store.ListUser(ctx, req.PageSize, Offset)
	if err != nil {
		s.logger.Error("Failed to get list: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to get list of users",
		})
		return
	}
	s.logger.Debug("get the users list: ", userList)
	ctx.JSON(http.StatusOK, userList)
}

type searchUserRequest struct {
	Name string `json:"name"`
}

func (s *HandlersServer) SearchUserFunc(ctx *gin.Context) {
	var req searchUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"code":    "VALIDATION_ERROR",
			"message": "Failed to bind the request",
		})
		return
	}
	searchQuery := "%" + req.Name + "%"
	fmt.Println("Search Query: ", searchQuery)

	response, err := s.store.SearchUser(ctx, searchQuery)
	if err != nil {
		s.logger.Error("Failed to search team : ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"code":    "INTERNAL_ERROR",
			"message": "Failed to search user",
		})
		return
	}

	s.logger.Debug("User search: ", response)

	ctx.JSON(http.StatusAccepted, response)
	return
}
