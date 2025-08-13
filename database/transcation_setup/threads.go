package transcation_setup

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/google/uuid"
)

func (store *SQLStore) CreateThreadTx(ctx context.Context, arg database.CreateThreadParams) (models.Thread, error) {
	var thread models.Thread

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create thread
		thread, err = q.CreateThread(ctx, arg)
		if err != nil {
			return err
		}
		return err
	})
	return thread, err
}

func (store *SQLStore) CreateCommentTx(ctx context.Context, threadPublicID uuid.UUID, userPublicID uuid.UUID, commentText string) (*models.Comment, error) {
	var comment *models.Comment

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create thread
		comment, err = q.CreateComment(ctx, threadPublicID, userPublicID, commentText)
		if err != nil {
			return err
		}
		return err
	})
	return comment, err
}

func (store *SQLStore) CreateLikeTx(ctx context.Context, userPublicID uuid.UUID, threadPublicID uuid.UUID) (models.UserLikeThread, error) {
	var likeThread models.UserLikeThread

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create thread
		likeThread, err = q.CreateLike(ctx, userPublicID, threadPublicID)
		if err != nil {
			return err
		}
		return err
	})
	return likeThread, err
}
