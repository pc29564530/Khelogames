package transcation_setup

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"
)

func (store *SQLStore) CreatePlayerTx(ctx context.Context, arg database.NewPlayerParams) (models.Player, error) {
	var player models.Player

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create thread
		player, err = q.NewPlayer(ctx, arg)
		if err != nil {
			return err
		}
		return err
	})
	return player, err
}
