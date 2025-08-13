package transcation_setup

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/google/uuid"
)

func (store *SQLStore) CreateMatchTx(ctx context.Context, arg database.NewMatchParams) (models.Match, error) {
	var match models.Match

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		// Create user
		match, err = q.NewMatch(ctx, arg)
		if err != nil {
			return err
		}
		return err
	})
	return match, err
}

func (store *SQLStore) UpdateMatchStatusTx(ctx context.Context, matchPublicID uuid.UUID, status string) (models.Match, error) {
	var match models.Match

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		// Create user
		match, err = q.UpdateMatchStatus(ctx, matchPublicID, status)
		if err != nil {
			return err
		}
		return err
	})
	return match, err
}

func (store *SQLStore) UpdateMatchResultTx(ctx context.Context, matchID int32, resultID int32) (models.Match, error) {
	var match models.Match

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		// Create user
		match, err = q.UpdateMatchResult(ctx, matchID, resultID)
		if err != nil {
			return err
		}
		return err
	})
	return match, err
}
