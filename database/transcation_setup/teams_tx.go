package transcation_setup

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"
)

func (store *SQLStore) CreateTeamsTx(ctx context.Context, arg database.NewTeamsParams) (models.Team, error) {
	var team models.Team

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		// Create user
		team, err = q.NewTeams(ctx, arg)
		if err != nil {
			return err
		}
		return err
	})
	return team, err
}
