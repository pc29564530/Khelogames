package transactions

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/google/uuid"
)

func (store *SQLStore) AddJoinCommunityTx(ctx context.Context, communityPublicID, userPublicID uuid.UUID) (*models.JoinCommunity, error) {
	var communityUser *models.JoinCommunity

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		// Add user to join_community table
		communityUser, err = q.AddJoinCommunity(ctx, communityPublicID, userPublicID)
		if err != nil {
			store.logger.Error("Failed to join community: ", err)
			return err
		}
		store.logger.Info("Successfully joined community: ", communityUser)

		// Increment member count
		err = q.IncrementCommunityMemberCount(ctx, communityPublicID)
		if err != nil {
			store.logger.Error("Failed to increment member count: ", err)
		}
		return err
	})
	return communityUser, err
}
