package auth

// type createMobileRequest struct {
// 	MobileNumber string `json:"mobile_number"`
// 	Otp          string `json:"otp"`
// }

// func (s *AuthServer) CreateMobileSignUp(ctx *gin.Context) {

// 	tx, err := s.store.BeginTx(ctx)
// 	if err != nil {
// 		s.logger.Error("Failed to begin the transcation: ", err)
// 		return
// 	}
// 	defer tx.Rollback()

// 	var req createMobileRequest
// 	err = ctx.ShouldBindJSON(&req)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			s.logger.Error("No row found: ", err)
// 			ctx.JSON(http.StatusNotFound, (err))
// 			return
// 		}

// 		errCode := db.ErrorCode(err)
// 		if errCode == db.UniqueViolation {
// 			s.logger.Error("Unique violation error: ", err)
// 			ctx.JSON(http.StatusForbidden, (err))
// 			return
// 		}
// 		s.logger.Error("Failed to bind: ", err)
// 		ctx.JSON(http.StatusInternalServerError, (err))
// 		return
// 	}

// 	s.logger.Info("MobileNUmber: ", req.MobileNumber)
// 	s.logger.Info("Otp: ", req.Otp)

// 	verifyOTP, err := s.store.GetSignup(ctx, req.MobileNumber)
// 	if verifyOTP.MobileNumber != req.MobileNumber {
// 		s.logger.Error("Failed to verify mobile: ", err)
// 		ctx.JSON(http.StatusNotFound, (err))
// 		return
// 	}

// 	s.logger.Debug(fmt.Sprintf("successfully get the otp: %v ", verifyOTP))

// 	if verifyOTP.Otp != req.Otp {
// 		s.logger.Error("Failed to verify otp: ", err)
// 		ctx.JSON(http.StatusNotFound, (err))
// 		return
// 	}

// 	s.logger.Debug(fmt.Sprintf("Successfully verified OTP for mobile number: %v", req.MobileNumber))

// 	ctx.JSON(http.StatusAccepted, gin.H{"mobile_number": verifyOTP.MobileNumber})
// }

// func (s *AuthServer) CreateMobileSignIn(ctx *gin.Context) {

// 	tx, err := s.store.BeginTx(ctx)
// 	if err != nil {
// 		s.logger.Error("Failed to begin the transcation: ", err)
// 		return
// 	}

// 	defer tx.Rollback()

// 	var req createMobileRequest
// 	err = ctx.ShouldBindJSON(&req)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			s.logger.Error("No row found: ", err)
// 			ctx.JSON(http.StatusNotFound, (err))
// 			return
// 		}

// 		errCode := db.ErrorCode(err)
// 		if errCode == db.UniqueViolation {
// 			s.logger.Error("Unique violation error: ", err)
// 			ctx.JSON(http.StatusForbidden, (err))
// 			return
// 		}
// 		s.logger.Error("Failed to bind: ", err)
// 		ctx.JSON(http.StatusInternalServerError, (err))
// 		return
// 	}

// 	s.logger.Info("MobileNUmber: ", req.MobileNumber)
// 	s.logger.Info("Otp: ", req.Otp)

// 	verifyOTP, err := s.store.GetSignup(ctx, req.MobileNumber)
// 	if verifyOTP.MobileNumber != req.MobileNumber {
// 		s.logger.Error("Failed to verify mobile: ", err)
// 		ctx.JSON(http.StatusNotFound, (err))
// 		return
// 	}

// 	s.logger.Debug(fmt.Sprintf("successfully get the otp: %v ", verifyOTP))
// 	if verifyOTP.Otp != req.Otp {
// 		s.logger.Error("Failed to verify otp: ", err)
// 		ctx.JSON(http.StatusNotFound, (err))
// 		return
// 	}

// 	s.logger.Debug(fmt.Sprintf("Successfully verified OTP for mobile number: %v", req.MobileNumber))

// 	user, err := s.store.GetUserByMobileNumber(ctx, verifyOTP.MobileNumber)
// 	if err != nil && err == sql.ErrNoRows {
// 		s.logger.Error("Error checking if user does not exists: ", err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 		return
// 	}

// 	tokens := CreateNewToken(ctx, user.Username, s, tx)

// 	session := tokens["session"].(models.Session)
// 	accessToken := tokens["accessToken"].(string)
// 	accessPayload := tokens["accessPayload"].(*token.Payload)
// 	refreshToken := tokens["refreshToken"].(string)
// 	refreshPayload := tokens["refreshPayload"].(*token.Payload)

// 	// rsp := loginUserResponse{
// 	// 	SessionID:             session.ID,
// 	// 	AccessToken:           accessToken,
// 	// 	AccessTokenExpiresAt:  accessPayload.ExpiredAt,
// 	// 	RefreshToken:          refreshToken,
// 	// 	RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
// 	// 	User: userResponse{
// 	// 		Username:     user.Username,
// 	// 		MobileNumber: getStringValue(user.MobileNumber),
// 	// 		Role:         user.Role,
// 	// 		Gmail:        getStringValue(user.Gmail),
// 	// 	},
// 	// }

// 	// s.logger.Info("Successfully sign in account ", rsp)
// 	// ctx.JSON(http.StatusAccepted, rsp)

// }

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
