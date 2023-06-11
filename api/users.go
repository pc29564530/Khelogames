package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
	db "khelogames/db/sqlc"
	"khelogames/util"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type createUserRequest struct {
	Username       string `json:"username"`
	MobileNumber   string `json:"mobile_number"`
	HashedPassword string `json:"hashed_password"`
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

	hashedPassword, err := util.HashPassword(req.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		MobileNumber:   req.MobileNumber,
		HashedPassword: hashedPassword,
	}

	//signup parameter
	argSignup := db.CreateSignupParams{
		MobileNumber: req.MobileNumber,
		Otp:          otp,
	}
	// creating a signup to be store in datbase
	signup, err := server.store.CreateSignup(ctx, argSignup)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, signup)

	err = server.sendOTP(client, arg.MobileNumber, otp)
	if err != nil {
		fmt.Println("unable to send otp")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	fmt.Println("Otp has been send successfully")

	arg = db.CreateUserParams{
		Username:       req.Username,
		MobileNumber:   req.MobileNumber,
		HashedPassword: hashedPassword,
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

// Sending a otp to verify the user mobile number
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

type getUsersRequest struct {
	Username string `josn:"username"`
}

func (server *Server) getUsers(ctx *gin.Context) {
	var req getUsersRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	users, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
	return
}

type getListUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	var req getListUsersRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.ListUserParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	userList, err := server.store.ListUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, userList)
}
