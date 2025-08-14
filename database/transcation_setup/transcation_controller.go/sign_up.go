package transcationcontroller

import (
	"context"
)

func (store *SQLStore) CreateEmailSignUpTx(ctx context.Context, arg CreateEmailSignUpTxParams) (CreateEmailSignUpTxResult, error) {
	var result CreateEmailSignUpTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create user
		result.User, err = q.CreateEmailSignUp(ctx, arg.FullName, arg.Username, arg.Email, arg.HashPassword)
		if err != nil {
			return err
		}

		// Create profile
		profileArg := CreateProfileParams{
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
