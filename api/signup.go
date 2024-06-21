package api

import (
	"database/sql"
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
			server.logger.Error("No row found: %v", err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			server.logger.Error("Unique violation error: %v", err)
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		server.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.logger.Info("MobileNUmber: ", req.MobileNumber)
	server.logger.Info("Otp: ", req.Otp)

	verifyOTP, err := server.store.GetSignup(ctx, req.MobileNumber)
	if verifyOTP.MobileNumber != req.MobileNumber {
		server.logger.Error("Failed to verify mobile: %v", err)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if verifyOTP.Otp != req.Otp {
		server.logger.Error("Failed to verify otp: %v", err)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	server.logger.Info("Successfully created account %w", http.StatusOK)
	return
}
