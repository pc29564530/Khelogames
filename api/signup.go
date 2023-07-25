package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//const (
//	accountSid       = "AC678672a16c66b33b075c556dfd805ad1"
//	authToken        = "7c83405dfef243da0cd68c22792444af"
//	verifyServiceSid = "VAd8d998e67283e3bb85bcf5c4f21682ad"
//)

type createSignupResponse struct {
	MobileNumber string    `json:"mobile_number"`
	Otp          string    `json:"otp"`
	CreatedAt    time.Time `json:"created_at"`
}

type createSignupRequest struct {
	MobileNumber string `json:"mobile_number"`
	Otp          string `json:"otp"`
}

//func generateOtp() string {
//	rand.Seed(time.Now().UnixNano())
//	otp := strconv.Itoa(rand.Intn(899999))
//	return otp
//}

func (server *Server) createSignup(ctx *gin.Context) {

	//client := twilio.NewRestClientWithParams(twilio.ClientParams{Username: accountSid, Password: authToken})

	var req createSignupRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//otp := generateOtp()

	//arg := db.CreateSignupParams{
	//	MobileNumber: req.MobileNumber,
	//	Otp:          otp,
	//}

	//signup, err := server.store.CreateSignup(ctx, arg)
	//if err != nil {
	//	if err == sql.ErrNoRows {
	//		ctx.JSON(http.StatusNotFound, errorResponse(err))
	//		return
	//	}
	//	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	//	return
	//}
	//
	//ctx.JSON(http.StatusOK, signup)

	//err = server.sendOTP(client, arg.MobileNumber, otp)
	//if err != nil {
	//	fmt.Println("unable to send otp")
	//	ctx.JSON(http.StatusNotFound, errorResponse(err))
	//	return
	//}
	//fmt.Println("Otp has been send successfully")

	verifyOTP, err := server.store.GetSignup(ctx, req.MobileNumber)
	if verifyOTP.MobileNumber != req.MobileNumber {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	if verifyOTP.Otp != req.Otp {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, verifyOTP)
	//fmt.Printf("Successfully created account %w", http.StatusOK)
	return
}
