package handlers

import (
	"fmt"
	errorhandler "khelogames/error_handler"
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
	Username     string `json:"username" binding:"required,min=3,max=50"`
	MobileNumber string `json:"mobile_number" binding:"required,min=10,max=15"`
	Role         string `json:"role" binding:"required,oneof=user admin moderator"`
	Gmail        string `json:"gmail" binding:"required,email"`
}

type userResponse struct {
	Username     string `json:"username"`
	MobileNumber string `json:"mobile_number"`
	Role         string `json:"role"`
	Gmail        string `json:"gmail"`
}

type getUserRequest struct {
	PublicID string `uri:"public_id" binding:"required"`
}

func (s *HandlersServer) GetUsersFunc(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("bind the reqeust: ", req)

	publicID, err := uuid.Parse(req.PublicID)
	if err != nil {
		s.logger.Error("Invalid UUID format", err)
		fieldErrors := map[string]string{"public_id": "Invalid UUID format"}
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}

	users, err := s.store.GetUser(ctx, publicID)
	if err != nil {
		s.logger.Error("Failed to get user: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("get the user data: ", users)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

type getListUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *HandlersServer) ListUsersFunc(ctx *gin.Context) {
	var req getListUsersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	s.logger.Debug("bind the request: ", req)
	Offset := (req.PageID - 1) * req.PageSize

	userList, err := s.store.ListUser(ctx, req.PageSize, Offset)
	if err != nil {
		s.logger.Error("Failed to get list: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get list of users",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}
	s.logger.Debug("get the users list: ", userList)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    userList,
	})
}

type searchUserRequest struct {
	Name string `json:"name" binding:"required,min=1"`
}

func (s *HandlersServer) SearchUserFunc(ctx *gin.Context) {
	var req searchUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fieldErrors := errorhandler.ExtractValidationErrors(err)
		errorhandler.ValidationErrorResponse(ctx, fieldErrors)
		return
	}
	searchQuery := "%" + req.Name + "%"
	fmt.Println("Search Query: ", searchQuery)

	response, err := s.store.SearchUser(ctx, searchQuery)
	if err != nil {
		s.logger.Error("Failed to search team : ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to search user",
			},
			"request_id": ctx.GetString("request_id"),
		})
		return
	}

	s.logger.Debug("User search: ", response)

	ctx.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data":    response,
	})
}
