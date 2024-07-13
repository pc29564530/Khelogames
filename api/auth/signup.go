package auth

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

func (s *AuthServer) CreateSignupFunc(ctx *gin.Context) {

	var req createSignupRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row found: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}

		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			s.logger.Error("Unique violation error: %v", err)
			ctx.JSON(http.StatusForbidden, (err))
			return
		}
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Info("MobileNUmber: ", req.MobileNumber)
	s.logger.Info("Otp: ", req.Otp)

	verifyOTP, err := s.store.GetSignup(ctx, req.MobileNumber)
	if verifyOTP.MobileNumber != req.MobileNumber {
		s.logger.Error("Failed to verify mobile: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}

	s.logger.Debug(fmt.Sprintf("successfully get the otp: %v ", verifyOTP))

	if verifyOTP.Otp != req.Otp {
		s.logger.Error("Failed to verify otp: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	s.logger.Info("Successfully created account %w", http.StatusOK)
	return
}
