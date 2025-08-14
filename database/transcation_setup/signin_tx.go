package transcation_setup

import (
	"context"
	"khelogames/database"
)

func (store *SQLStore) CreateRenewAccessTx(ctx context.Context, fullName, username, email, hashPassword string) (CreateSignUpTxResult, error) {
	var result CreateSignUpTxResult

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		result.User, err = q.CreateEmailSignUp(ctx, fullName, username, email, hashPassword)
		if err != nil {
			return err
		}

		// Create profile
		profileArg := database.CreateProfileParams{
			UserID:    int32(result.User.ID),
			Bio:       "",
			AvatarUrl: "",
		}
		result.Profile, err = q.CreateProfile(ctx, profileArg)
		if err != nil {
			return err
		}
		return nil
	})

	return result, err
}
