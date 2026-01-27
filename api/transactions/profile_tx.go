package transactions

import (
	"context"
	"fmt"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/google/uuid"
)

func (store *SQLStore) UpdateProfileTx(ctx context.Context, publicID uuid.UUID, bio string, avatarUrl string, fullName, city, state, country string, latitude, longitude float64, h3Index string) (*models.UserProfiles, error) {
	var userProfiles *models.UserProfiles
	err := store.execTx(ctx, func(q *database.Queries) error {
		location, err := q.AddLocation(ctx, city, state, country, latitude, longitude, h3Index)
		if err != nil {
			store.logger.Error("Failed to add location: ", err)
			return err
		}

		profile, err := store.GetProfile(ctx, publicID)
		if err != nil {
			store.logger.Error("Failed to get profile: ", err)
			return err
		}

		arg := database.EditProfileParams{
			PublicID:   profile.PublicID,
			Bio:        bio,
			AvatarUrl:  avatarUrl,
			LocationID: int32(location.ID),
		}
		fmt.Println("Arg: ", arg)
		updatedProfile, err := q.EditProfile(ctx, arg)
		if err != nil {
			store.logger.Error("Failed to update profile: ", err)
			return err
		}
		fmt.Println("Update Profile: ", updatedProfile)
		_, err = q.UpdateUser(ctx, int32(updatedProfile.UserID), fullName)
		if err != nil {
			store.logger.Error("Failed to update the user full name: ", err)
			return err
		}

		userProfiles = updatedProfile
		return nil
	})
	return userProfiles, err
}
