package transactions

import (
	"fmt"
	coreToken "khelogames/core/token"
	"khelogames/database"
	"khelogames/database/models"
	util "khelogames/util"

	"github.com/gin-gonic/gin"
)

// CreateEmailSignUpTx performs email signup transaction
func (store *SQLStore) CreateEmailSignUpTx(ctx *gin.Context, config util.Config, st *database.Store, fullName, username, email, hashPassword string) (*models.UserProfiles, *models.Users, map[string]interface{}, error) {
	clientIP := ctx.ClientIP()
	accessDuration := config.AccessTokenDuration
	refreshDuration := config.RefreshTokenDuration
	var tokens map[string]interface{}
	var userSignUp *models.Users
	var profile *models.UserProfiles

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		userSignUp, err = q.CreateEmailSignUp(ctx, fullName, username, email, hashPassword)
		if err != nil {
			store.logger.Error("Failed to create email signup: ", err)
			return err
		}

		tokens, err = coreToken.CreateNewTokenTx(ctx, q, *store.tokenMaker, int32(userSignUp.ID), userSignUp.PublicID, accessDuration, refreshDuration, ctx.Request.UserAgent(), clientIP)
		if err != nil {
			store.logger.Error("Failed to create new token: ", err)
			return err
		}
		store.logger.Info("Successfully created email sign-up for: ", email)

		arg := database.CreateProfileParams{
			UserID:    int32(userSignUp.ID),
			Bio:       "",
			AvatarUrl: "",
		}

		profile, err = q.CreateProfile(ctx, arg)
		if err != nil {
			store.logger.Error("Failed to create profile: ", err)
			return err
		}

		return nil
	})

	return profile, userSignUp, tokens, err
}

func (store *SQLStore) CreateGoogleSignUpTx(ctx *gin.Context, config util.Config, name, username, email, googleID string, avatarUrl string) (*models.UserProfiles, *models.Users, map[string]interface{}, error) {
	accessDuration := config.AccessTokenDuration
	refreshDuration := config.RefreshTokenDuration
	clientIP := ctx.ClientIP()
	var tokens map[string]interface{}
	var userSignUp *models.Users
	var profile *models.UserProfiles

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		userSignUp, err = q.CreateGoogleSignUp(ctx, name, username, email, googleID)
		if err != nil {
			return fmt.Errorf("Failed to create google signup: ", err)
		}

		// Create tokens and session
		tokens, err = coreToken.CreateNewTokenTx(ctx, q, *store.tokenMaker, int32(userSignUp.ID), userSignUp.PublicID, accessDuration, refreshDuration, ctx.Request.UserAgent(), clientIP)
		if err != nil {
			return fmt.Errorf("Failed to create new token: ", err)
		}

		// Create user profile
		arg := database.CreateProfileParams{
			UserID:    int32(userSignUp.ID),
			Bio:       "",
			AvatarUrl: avatarUrl,
		}

		profile, err = q.CreateProfile(ctx, arg)
		if err != nil {
			store.logger.Error("Failed to create profile: ", err)
			return err
		}

		store.logger.Info("Successfully created user_profile")
		return nil
	})
	return profile, userSignUp, tokens, err
}
