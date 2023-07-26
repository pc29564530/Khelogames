package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//const (
//	accountSid       = "AC678672a16c66b33b075c556dfd805ad1"
//	authToken        = "7c83405dfef243da0cd68c22792444af"
//	verifyServiceSid = "VAd8d998e67283e3bb85bcf5c4f21682ad"
//)

type createSignupResponse struct {
	MobileNumber string    `json:"mobile_number"`
	Otp          string    `json:"otp"`
	CreatedAt    time.Time `json:"created_at"`
}

type createSignupRequest struct {
	MobileNumber string `json:"mobileNumber"`
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

	verifyOTP, err := server.store.GetSignup(ctx, req.MobileNumber)
	if verifyOTP.MobileNumber != req.MobileNumber {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	fmt.Println(verifyOTP.Otp)
	fmt.Println(req.Otp)
	if verifyOTP.Otp != req.Otp {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, verifyOTP)
	fmt.Printf("Successfully created account %w", http.StatusOK)
	return
}
