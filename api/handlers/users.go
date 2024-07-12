package handlers

import (
	"database/sql"
	db "khelogames/db/sqlc"

	"khelogames/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func authorizationCode(ctx *gin.Context, username string, mobileNumber string, role string, s *HandlersServer) {
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		s.logger.Error("Failed to create access token", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("created a accesstoken: %v", accessToken)

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		username,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		s.logger.Error("Failed to create refresh token", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Debug("created a refresh token: %v", refreshToken)

	session, err := s.store.CreateSessions(ctx, db.CreateSessionsParams{
		ID:           refreshPayload.ID,
		Username:     username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		s.logger.Error("Failed to create session", err)
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
	s.logger.Info("User logged in successfully")
	ctx.JSON(http.StatusAccepted, rsp)
	return
}

func (s *HandlersServer) CreateUserFunc(ctx *gin.Context) {

	var req createUserRequest

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		errCode := db.ErrorCode(err)
		if errCode == db.UniqueViolation {
			s.logger.Error("Unique violation error ", err)
			ctx.JSON(http.StatusForbidden, (err))
			return
		}
		if err == sql.ErrNoRows {
			s.logger.Error("No row data", err)
			ctx.JSON(http.StatusNotFound, (err))
			return
		}
		s.logger.Error("Error while binding JSON", err)
		ctx.JSON(http.StatusBadGateway, (err))
		return
	}

	s.logger.Debug("bind the request: %v", req)

	hashedPassword, err := util.HashPassword(req.HashedPassword)
	if err != nil {
		s.logger.Error("Failed to hash password", err)
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
		s.logger.Error("Unable to create user")
		ctx.JSON(http.StatusUnauthorized, (err))
		return
	}

	resp := createUserResponse{
		Username:     user.Username,
		MobileNumber: user.MobileNumber,
		Role:         user.Role,
	}

	s.logger.Debug("successfully created user: %v", resp)

	authorizationCode(ctx, resp.Username, resp.MobileNumber, resp.Role, s)

	_, err = s.store.DeleteSignup(ctx, req.MobileNumber)
	if err != nil {
		s.logger.Error("Unable to delete signup details: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	s.logger.Debug("delete the signup details")

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
	s.logger.Info("Profile created successfully")
	s.logger.Info("Successfully created the user")
	return
}

type getUserRequest struct {
	Username string `uri:"username"`
}

func (s *HandlersServer) GetUsersFunc(ctx *gin.Context) {
	var req getUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the reqeust: %v", req)

	users, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		s.logger.Error("Failed to get user: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get the user data: %v", users)
	ctx.JSON(http.StatusOK, users)
	return
}

type getListUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *HandlersServer) ListUsersFunc(ctx *gin.Context) {
	var req getListUsersRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		s.logger.Error("Failed to bind: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: %v", req)
	arg := db.ListUserParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	userList, err := s.store.ListUser(ctx, arg)
	if err != nil {
		s.logger.Error("Failed to get list: %v", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get the users list: %v", userList)
	ctx.JSON(http.StatusOK, userList)
}
