package auth

import (
	"database/sql"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignupServer struct {
	store  *db.Store
	logger *logger.Logger
}

type createSignupRequest struct {
	MobileNumber string `json:"mobile_number"`
	Otp          string `json:"otp"`
}

func NewSignupServer(store *db.Store, logger *logger.Logger) *SignupServer {
	return &SignupServer{store: store, logger: logger}
}

func (s *SignupServer) CreateSignupFunc(ctx *gin.Context) {

	var req createSignupRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("No row found: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}

		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			fmt.Errorf("Unique violation error: %v", err)
			ctx.JSON(http.StatusForbidden, (err))
			return
		}
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	fmt.Println("MobileNUmber: ", req.MobileNumber)
	fmt.Println("Otp: ", req.Otp)

	verifyOTP, err := s.store.GetSignup(ctx, req.MobileNumber)
	if verifyOTP.MobileNumber != req.MobileNumber {
		fmt.Errorf("Failed to verify mobile: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}

	if verifyOTP.Otp != req.Otp {
		fmt.Errorf("Failed to verify otp: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	fmt.Println("Successfully created account %w", http.StatusOK)
	return
}
