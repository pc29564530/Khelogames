package transcation_setup

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/google/uuid"
)

func (store *SQLStore) CreateUserConnectionTx(ctx context.Context, userPublicID, targetPublicID uuid.UUID) (models.UsersConnections, error) {
	var userConnection models.UsersConnections
	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		userConnection, err = store.CreateUserConnections(ctx, userPublicID, targetPublicID)
		if err != nil {
			store.logger.Error("Failed to create following: ", err)
			return err
		}
		return err
	})
	return userConnection, err
}

func (store *SQLStore) DeleteUserConnectionsTx(ctx context.Context, userPublicID, targetPublicID uuid.UUID) error {
	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		err = store.DeleteUserConnectionsTx(ctx, userPublicID, targetPublicID)
		if err != nil {
			store.logger.Error("Failed to create following: ", err)
			return err
		}
		return err
	})
	return err
}
