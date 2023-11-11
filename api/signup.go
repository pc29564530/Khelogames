package api

import (
	"database/sql"
	"fmt"
	db "khelogames/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println("MobileNUmber: ", req.MobileNumber)
	fmt.Println("Otp: ", req.Otp)

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
