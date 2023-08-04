package api

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sfreiberg/gotwilio"
	db "khelogames/db/sqlc"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type createSendOtpRequest struct {
	MobileNumber string `json:"mobileNumber"`
}

func generateOtp() string {
	rand.Seed(time.Now().UnixNano())
	otp := strconv.Itoa(rand.Intn(899999))
	return otp
}

func (server *Server) Otp(ctx *gin.Context) {

	var reqSendOTP createSendOtpRequest
	err := ctx.ShouldBindJSON(&reqSendOTP)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	otp := generateOtp()

	err = server.sendOTP(reqSendOTP.MobileNumber, otp)
	if err != nil {
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

func (server *Server) sendOTP(mobileNumber string, otp string) error {

	err := godotenv.Load("./app.env")
	if err != nil {
		fmt.Errorf("Unable to read env file: ", err)
	}
	AccountSid := os.Getenv("ACCOUNT_SID")
	AuthToken := os.Getenv("AUTH_TOKEN")
	TwilioPhoneNumber := os.Getenv("YOUR_TWILIO_PHONE_NUMBER")

	twilioClient := gotwilio.NewTwilioClient(AccountSid, AuthToken)

	// Send SMS OTP to the user's phone number
	_, _, err = twilioClient.SendSMS(TwilioPhoneNumber, mobileNumber, "Your OTP is: "+otp, "", "")
	return err

}
