package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
	db "khelogames/db/sqlc"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
)

type createUserRequest struct {
	Username     string `json:"username"`
	MobileNumber string `json:"mobile_number"`
}

type createUserResponse struct {
	Username     string    `json:"username"`
	MobileNumber string    `json:"mobile_number"`
	CreatedAt    time.Time `json:"created_at"`
}

const (
	accountSid       = "AC678672a16c66b33b075c556dfd805ad1"
	authToken        = "7c83405dfef243da0cd68c22792444af"
	verifyServiceSid = "VAd8d998e67283e3bb85bcf5c4f21682ad"
)

func generateOtp() string {
	rand.Seed(time.Now().UnixNano())
	otp := strconv.Itoa(rand.Intn(899999))
	return otp
}

func (server *Server) createUser(ctx *gin.Context) {

	var req createUserRequest

	client := twilio.NewRestClientWithParams(twilio.ClientParams{Username: accountSid, Password: authToken})

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusBadGateway, errorResponse(err))
		return
	}

	otp := generateOtp()

	arg := db.CreateUserParams{
		Username:     req.Username,
		MobileNumber: req.MobileNumber,
	}

	argSignup := db.CreateSignupParams{
		MobileNumber: req.MobileNumber,
		Otp:          otp,
	}

	signup, err := server.store.CreateSignup(ctx, argSignup)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	respSignup := createSignupResponse{
		MobileNumber: signup.MobileNumber,
		Otp:          signup.Otp,
		CreatedAt:    signup.CreatedAt,
	}

	ctx.JSON(http.StatusOK, respSignup)

	err = server.sendOTP(client, arg.MobileNumber, otp)
	if err != nil {
		fmt.Println("unable to send otp")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	resp := createUserResponse{
		Username:     user.Username,
		MobileNumber: user.MobileNumber,
		CreatedAt:    user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, resp)
	return
}

// Create Sign Up function
func (server *Server) sendOTP(client *twilio.RestClient, to string, otp string) error {

	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(verifyServiceSid, params)
	if err != nil {
		fmt.Errorf("Unable to create a message: %s", err)
		return err
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Printf(string(response))
	}

	return nil
}
