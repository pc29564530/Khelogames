package handlers

import (
	"database/sql"
	db "khelogames/database"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}
type createUserRequest struct {
	Username     string `json:"username"`
	MobileNumber string `json:"mobile_number"`
	Role         string `json:"role"`
	Gmail        string `json:"gmail"`
	SignUpType   string `json:"signup_type"`
}

type userResponse struct {
	Username     string `json:"username"`
	MobileNumber string `json:"mobile_number"`
	Role         string `json:"role"`
	Gmail        string `json:"gmail"`
}

func CreateNewToken(ctx *gin.Context, username string, s *HandlersServer, tx *sql.Tx) map[string]interface{} {
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		username,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create access token", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return nil
	}
	s.logger.Debug("created a accesstoken: ", accessToken)

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		username,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create refresh token", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return nil
	}

	s.logger.Debug("created a refresh token: ", refreshToken)

	session, err := s.store.CreateSessions(ctx, db.CreateSessionsParams{
		ID:           refreshPayload.ID,
		Username:     username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create session", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return nil
	}

	return map[string]interface{}{
		"accessToken":    accessToken,
		"accessPayload":  accessPayload,
		"refreshToken":   refreshToken,
		"refreshPayload": refreshPayload,
		"session":        session,
	}
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

	s.logger.Debug("bind the request: ", req)

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin the transcation: ", err)
		return
	}
	defer tx.Rollback()

	user, err := s.store.CreateUser(ctx, req.Username, req.MobileNumber, req.Role, req.Gmail)
	if err != nil {
		tx.Rollback()
		s.logger.Error("Unable to create user: ", err)
		ctx.JSON(http.StatusUnauthorized, (err))
		return
	}

	s.logger.Debug("successfully created user: ", user)

	tokens := CreateNewToken(ctx, user.Username, s, tx)

	session := tokens["session"].(map[string]interface{})
	accessToken := tokens["accessToken"].(string)
	accessPayload := tokens["accessPayload"].(map[string]interface{})
	refreshToken := tokens["refreshToken"].(string)
	refreshPayload := tokens["refreshPayload"].(map[string]interface{})

	rsp := loginUserResponse{
		SessionID:             session["id"].(uuid.UUID),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload["expired_at"].(time.Time),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload["expired_at"].(time.Time),
		User: userResponse{
			Username:     user.Username,
			MobileNumber: *user.MobileNumber,
			Role:         user.Role,
			Gmail:        *user.Gmail,
		},
	}

	ctx.JSON(http.StatusAccepted, rsp)

	if req.SignUpType == "mobile" {
		_, err = s.store.DeleteSignup(ctx, req.MobileNumber)
		if err != nil {
			tx.Rollback()
			s.logger.Error("Unable to delete signup details: ", err)
			ctx.JSON(http.StatusInternalServerError, (err))
			return
		}
		s.logger.Debug("delete the signup details")
	}

	//createProfile
	argProfile := db.CreateProfileParams{
		Owner:     user.Username,
		FullName:  "",
		Bio:       "",
		AvatarUrl: "",
	}

	_, err = s.store.CreateProfile(ctx, argProfile)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transcation: ", err)
		return
	}

	s.logger.Info("Profile created successfully")
	s.logger.Info("Successfully created the user")
}

type getUserRequest struct {
	Username string `uri:"username"`
}

func (s *HandlersServer) GetUsersFunc(ctx *gin.Context) {
	var req getUserRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the reqeust: ", req)

	users, err := s.store.GetUser(ctx, req.Username)
	if err != nil {
		s.logger.Error("Failed to get user: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("get the user data: ", users)
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
		s.logger.Error("Failed to bind: ", err)
		ctx.JSON(http.StatusInternalServerError, (err))
		return
	}
	s.logger.Debug("bind the request: ", req)
	// arg := db.ListUserParams{
	// 	Limit:  req.PageSize,
	// 	Offset: (req.PageID - 1) * req.PageSize,
	// }

	// userList, err := s.store.ListUser(ctx, username)
	// if err != nil {
	// 	s.logger.Error("Failed to get list: ", err)
	// 	ctx.JSON(http.StatusInternalServerError, (err))
	// 	return
	// }
	// s.logger.Debug("get the users list: ", userList)
	// ctx.JSON(http.StatusOK, userList)
}
