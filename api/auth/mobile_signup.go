package auth

import (
	"database/sql"
	"fmt"
	db "khelogames/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createMobileRequest struct {
	MobileNumber string `json:"mobile_number"`
	Otp          string `json:"otp"`
}

func (s *AuthServer) CreateMobileSignUp(ctx *gin.Context) {

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin the transcation: ", err)
		return
	}
	defer tx.Rollback()

	var req createMobileRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row found: ", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}

		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			s.logger.Error("Unique violation error: ", err)
			ctx.JSON(http.StatusForbidden, (err))
			return
		}
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Info("MobileNUmber: ", req.MobileNumber)
	s.logger.Info("Otp: ", req.Otp)

	verifyOTP, err := s.store.GetSignup(ctx, req.MobileNumber)
	if verifyOTP.MobileNumber != req.MobileNumber {
		s.logger.Error("Failed to verify mobile: ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}

	s.logger.Debug(fmt.Sprintf("successfully get the otp: %v ", verifyOTP))

	if verifyOTP.Otp != req.Otp {
		s.logger.Error("Failed to verify otp: ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}

	s.logger.Debug(fmt.Sprintf("Successfully verified OTP for mobile number: %v", req.MobileNumber))

	ctx.JSON(http.StatusAccepted, gin.H{"mobile_number": verifyOTP.MobileNumber})
}

func (s *AuthServer) CreateMobileSignIn(ctx *gin.Context) {

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin the transcation: ", err)
		return
	}

	defer tx.Rollback()

	var req createMobileRequest
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("No row found: ", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}

		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			s.logger.Error("Unique violation error: ", err)
			ctx.JSON(http.StatusForbidden, (err))
			return
		}
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Info("MobileNUmber: ", req.MobileNumber)
	s.logger.Info("Otp: ", req.Otp)

	verifyOTP, err := s.store.GetSignup(ctx, req.MobileNumber)
	if verifyOTP.MobileNumber != req.MobileNumber {
		s.logger.Error("Failed to verify mobile: ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}

	s.logger.Debug(fmt.Sprintf("successfully get the otp: %v ", verifyOTP))

	if verifyOTP.Otp != req.Otp {
		s.logger.Error("Failed to verify otp: ", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}

	s.logger.Debug(fmt.Sprintf("Successfully verified OTP for mobile number: %v", req.MobileNumber))

	user, err := s.store.GetUserByMobileNumber(ctx, verifyOTP.MobileNumber)
	if err != nil && err == sql.ErrNoRows {
		s.logger.Error("Error checking if user does not exists: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	tokens := CreateNewToken(ctx, user.Username, s, tx)

	session := tokens["session"].(map[string]interface{})
	accessToken := tokens["accessToken"].(string)
	accessPayload := tokens["accessPayload"].(map[string]interface{})
	refreshToken := tokens["refreshToken"].(string)
	refreshPayload := tokens["refreshPayload"].(map[string]interface{})

	rsp := loginUserResponse{
		SessionID:             session["id"].(uuid.UUID),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload["expired_at"].(time.Time),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload["expired_at"].(time.Time),
		User: userResponse{
			Username:     user.Username,
			MobileNumber: *user.MobileNumber,
			Role:         user.Role,
			Gmail:        *user.Gmail,
		},
	}

	s.logger.Info("Successfully sign in account ", rsp)
	ctx.JSON(http.StatusAccepted, rsp)

}
