package transactions

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/google/uuid"
)

func (store *SQLStore) UpdateProfileTx(ctx context.Context, publicID uuid.UUID, bio string, avatarUrl string, fullName string) (models.UserProfiles, error) {
	var userProfiles models.UserProfiles
	err := store.execTx(ctx, func(q *database.Queries) error {
		arg := database.EditProfileParams{
			PublicID:  publicID,
			Bio:       bio,
			AvatarUrl: avatarUrl,
		}
		updatedProfile, err := store.EditProfile(ctx, arg)
		if err != nil {
			store.logger.Error("Failed to update profile: ", err)
			return err
		}

		_, err = store.UpdateUser(ctx, int32(updatedProfile.UserID), fullName)
		if err != nil {
			store.logger.Error("Failed to update the user full name: ", err)
			return err
		}
		return nil
	})
	return userProfiles, err
}
