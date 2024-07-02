package auth

import (
	"database/sql"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sfreiberg/gotwilio"
)

type OtpServer struct {
	store  *db.Store
	logger *logger.Logger
}

type createSendOtpRequest struct {
	MobileNumber string `json:"mobile_number"`
}

func NewOtpServer(store *db.Store, logger *logger.Logger) *OtpServer {
	return &OtpServer{
		store:  store,
		logger: logger,
	}
}

func generateOtp() string {
	rand.Seed(time.Now().UnixNano())
	otp := strconv.Itoa(rand.Intn(899999))
	return otp
}

func (s *OtpServer) Otp(ctx *gin.Context) {

	var reqSendOTP createSendOtpRequest
	err := ctx.ShouldBindJSON(&reqSendOTP)

	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	otp := generateOtp()

	err = s.sendOTP(reqSendOTP.MobileNumber, otp)
	fmt.Println("line no 40: ", err)
	if err != nil {
		ctx.JSON(http.StatusNotFound, (err))
		return
	}
	fmt.Println("Otp has been send successfully")

	arg := db.CreateSignupParams{
		MobileNumber: reqSendOTP.MobileNumber,
		Otp:          otp,
	}

	signup, err := s.store.CreateSignup(ctx, arg)
	if err != nil {

		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, signup)

}

func (s *OtpServer) sendOTP(mobileNumber string, otp string) error {

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
