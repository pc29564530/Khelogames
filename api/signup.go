package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type createSignupResponse struct {
	MobileNumber string    `json:"mobile_number"`
	Otp          string    `json:"otp"`
	CreatedAt    time.Time `json:"created_at"`
}

type createSignupRequest struct {
	MobileNumber string `json:"mobile_number"`
	Otp          string `json:"otp"`
}

func (server *Server) createSignup(ctx *gin.Context) {
	var req createSignupRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := server.store.GetSignup(ctx, req.MobileNumber)
	if user.MobileNumber != req.MobileNumber {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if user.Otp != req.Otp {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	fmt.Printf("Successfully signup %s", http.StatusOK)
	return
}
