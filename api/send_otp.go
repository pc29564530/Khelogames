package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
	db "khelogames/db/sqlc"
	"khelogames/util"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type createSendOtpRequest struct {
	MobileNumber string `json:"mobile_number"`
}

func generateOtp() string {
	rand.Seed(time.Now().UnixNano())
	otp := strconv.Itoa(rand.Intn(899999))
	return otp
}

var config util.Config

func (server *Server) Otp(ctx *gin.Context) {
	err := godotenv.Load("./app.env")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	AccountSid := os.Getenv("ACCOUNT_SID")
	AuthToken := os.Getenv("AUTH_TOKEN")
	VerifyServiceToken := os.Getenv("VERIFY_SERVICE_SID")

	var reqSendOTP createSendOtpRequest
	fmt.Println(config.AccountSid)
	client := twilio.NewRestClientWithParams(twilio.ClientParams{Username: AccountSid, Password: AuthToken})
	err = ctx.ShouldBindJSON(&reqSendOTP)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	otp := generateOtp()
	err = server.sendOTP(client, reqSendOTP.MobileNumber, otp, VerifyServiceToken)
	if err != nil {
		fmt.Println("unable to send otp")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	fmt.Println("Otp has been send successfully")

	arg := db.CreateSignupParams{
		MobileNumber: reqSendOTP.MobileNumber,
		Otp:          otp,
	}

	signup, err := server.store.CreateSignup(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, signup)

}

// Sending a otp to verify the user mobile number
func (server *Server) sendOTP(client *twilio.RestClient, to string, otp string, VerifyServiceToken string) error {
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")
	resp, err := client.VerifyV2.CreateVerification(VerifyServiceToken, params)
	if err != nil {
		fmt.Errorf("Unable to create a message: %s", err)
		return err
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Printf(string(response))
	}

	return nil
}
