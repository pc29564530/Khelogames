package transcation_setup

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/google/uuid"
)

func (store *SQLStore) CreateCommunityTx(ctx context.Context, arg database.CreateCommunityParams) (models.Communities, error) {
	var community models.Communities

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		community, err = q.CreateCommunity(ctx, arg)
		if err != nil {
			return err
		}
		return err
	})
	return community, err
}

func (store *SQLStore) AddJoinCommunityTx(ctx context.Context, communityPublicID, userPublicID uuid.UUID) (models.JoinCommunity, error) {
	var communityUser models.JoinCommunity

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		// Add user to join_community table
		communityUser, err := store.AddJoinCommunity(ctx, communityPublicID, userPublicID)
		if err != nil {
			store.logger.Error("Failed to join community: ", err)
			return err
		}
		store.logger.Debug("Successfully joined community: ", communityUser)

		// Increment member count
		err = store.IncrementCommunityMemberCount(ctx, communityPublicID)
		if err != nil {
			store.logger.Error("Failed to increment member count: ", err)
		}
		return err
	})
	return communityUser, err
}
