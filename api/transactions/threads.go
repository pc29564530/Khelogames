package transactions

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/google/uuid"
)

// CreateLikeTx atomically inserts a like row and increments the thread like_count.
// Returns the updated Thread with the new like_count.
func (store *SQLStore) CreateLikeTx(ctx context.Context, userPublicID uuid.UUID, threadPublicID uuid.UUID) (*models.Thread, error) {
	var thread *models.Thread

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		_, err = q.CreateLike(ctx, userPublicID, threadPublicID)
		if err != nil {
			return err
		}

		thread, err = q.UpdateThreadLike(ctx, threadPublicID)
		if err != nil {
			return err
		}

		return nil
	})

	return thread, err
}

// DeleteLikeThreadTx atomically removes a like row and decrements the thread like_count.
// Returns the updated Thread with the new like_count.
func (store *SQLStore) DeleteLikeThreadTx(ctx context.Context, userPublicID uuid.UUID, threadPublicID uuid.UUID) (*models.Thread, error) {
	var thread *models.Thread

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		err = q.DeleteLike(ctx, userPublicID, threadPublicID)
		if err != nil {
			store.logger.Error("Failed to delete like: ", err)
			return err
		}

		thread, err = q.DecrementThreadLike(ctx, threadPublicID)
		if err != nil {
			store.logger.Error("Failed to decrement like count: ", err)
			return err
		}

		return nil
	})

	return thread, err
}
