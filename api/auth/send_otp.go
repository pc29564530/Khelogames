package auth

import (
	"database/sql"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sfreiberg/gotwilio"
)

type createSendOtpRequest struct {
	MobileNumber string `json:"mobile_number"`
}

func generateOtp() string {
	rand.Seed(time.Now().UnixNano())
	otp := strconv.Itoa(rand.Intn(899999))
	return otp
}

func (s *AuthServer) Otp(ctx *gin.Context) {

	var reqSendOTP createSendOtpRequest
	err := ctx.ShouldBindJSON(&reqSendOTP)

	if err != nil {
		s.logger.Error("unable to bind mobile number: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
	}

	s.logger.Debug("mobile number request: %v", err)

	otp := generateOtp()

	s.logger.Debug("Otp generate: %v", otp)

	err = s.sendOTP(reqSendOTP.MobileNumber, otp)
	if err != nil {
		s.logger.Error("unable to send otp: %v", err)
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	s.logger.Info("Otp has been send successfully")

	s.logger.Debug("signup arg: %v", err)

	signup, err := s.store.CreateSignup(ctx, reqSendOTP.MobileNumber, otp)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("no row in signup: %v", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("unable to bind signup: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Info("successfully singup")

	ctx.JSON(http.StatusOK, signup)
}

func (s *AuthServer) sendOTP(mobileNumber string, otp string) error {

	err := godotenv.Load("./app.env")
	if err != nil {
		s.logger.Error("Unable to read env file: ", err)
		return err
	}
	AccountSid := os.Getenv("ACCOUNT_SID")
	AuthToken := os.Getenv("AUTH_TOKEN")
	TwilioPhoneNumber := os.Getenv("YOUR_TWILIO_PHONE_NUMBER")

	twilioClient := gotwilio.NewTwilioClient(AccountSid, AuthToken)

	// Send SMS OTP to the user's phone number
	_, _, err = twilioClient.SendSMS(TwilioPhoneNumber, mobileNumber, "Your OTP is: "+otp, "", "")
	return err

}
