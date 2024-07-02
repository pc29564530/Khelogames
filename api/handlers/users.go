package handlers

import (
	"database/sql"
	"fmt"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/token"
	"khelogames/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserServer struct {
	store      *db.Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
}

type userResponse struct {
	Username     string    `json:"username"`
	MobileNumber string    `json:"mobile_number"`
	CreatedAt    time.Time `json:"created_at"`
	Role         string    `json:"role"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}
type createUserRequest struct {
	Username       string `json:"username"`
	MobileNumber   string `json:"mobile_number"`
	HashedPassword string `json:"password"`
	Role           string `json:"role"`
}

type createUserResponse struct {
	Username     string    `json:"username"`
	MobileNumber string    `json:"mobile_number"`
	CreatedAt    time.Time `json:"created_at"`
	Role         string    `json:"role"`
}

func NewUserServer(store *db.Store, logger *logger.Logger, tokenMaker token.Maker, config util.Config) *UserServer {
	return &UserServer{store: store, logger: logger, tokenMaker: tokenMaker, config: config}
}

func authorizationCode(ctx *gin.Context, username string, mobileNumber string, role string, s *UserServer) {
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		username,
		s.config.AccessTokenDuration,
	)
	fmt.Println("AccessToken: ", accessToken)
	if err != nil {
		fmt.Errorf("Failed to create access token", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		username,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		fmt.Errorf("Failed to create refresh token", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	session, err := s.store.CreateSessions(ctx, db.CreateSessionsParams{
		ID:           refreshPayload.ID,
		Username:     username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		fmt.Errorf("Failed to create session", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User: userResponse{
			Username:     username,
			MobileNumber: mobileNumber,
			Role:         role,
		},
	}
	fmt.Println("User logged in successfully")
	ctx.JSON(http.StatusAccepted, rsp)
	return
}

func (s *UserServer) CreateUserFunc(ctx *gin.Context) {

	var req createUserRequest

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			fmt.Errorf("Unique violation error ", err)
			ctx.JSON(http.StatusForbidden, (err))
			return
		}
		if err == sql.ErrNoRows {
			fmt.Errorf("No row data", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		fmt.Errorf("Error while binding JSON", err)
		ctx.JSON(http.StatusBadGateway, (err))
		return
	}

	hashedPassword, err := util.HashPassword(req.HashedPassword)
	if err != nil {
		fmt.Errorf("Failed to hash password", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		MobileNumber:   req.MobileNumber,
		HashedPassword: hashedPassword,
		Role:           req.Role,
	}

	user, err := s.store.CreateUser(ctx, arg)

	if err != nil {
		fmt.Errorf("Unable to create user")
		ctx.JSON(http.StatusUnauthorized, (err))

		return
	}

	resp := createUserResponse{
		Username:     user.Username,
		MobileNumber: user.MobileNumber,
		Role:         user.Role,
	}

	authorizationCode(ctx, resp.Username, resp.MobileNumber, resp.Role, s)

	_, err = s.store.DeleteSignup(ctx, req.MobileNumber)
	if err != nil {
		fmt.Errorf("Unable to delete signup details: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	//createProfile
	argProfile := db.CreateProfileParams{
		Owner:     resp.Username,
		FullName:  "",
		Bio:       "",
		AvatarUrl: "",
		CoverUrl:  "",
	}

	_, err = s.store.CreateProfile(ctx, argProfile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	fmt.Println("Profile created successfully")
	fmt.Println("Successfully created the user")
	return
}

type getUserRequest struct {
	Username string `uri:"username"`
}

func (s *UserServer) GetUsersFunc(ctx *gin.Context) {
	var req getUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	fmt.Println("Server: ", s)
	fmt.Println("Store: ", s.store)

	users, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		fmt.Errorf("Failed to get user: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	ctx.JSON(http.StatusOK, users)
	return
}

type getListUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *UserServer) ListUsersFunc(ctx *gin.Context) {
	var req getListUsersRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		fmt.Errorf("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	arg := db.ListUserParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	userList, err := s.store.ListUser(ctx, arg)
	if err != nil {
		fmt.Errorf("Failed to get list: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	ctx.JSON(http.StatusOK, userList)
}
