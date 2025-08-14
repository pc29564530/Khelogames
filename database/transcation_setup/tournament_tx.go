package transcation_setup

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"
)

func (store *SQLStore) CreateTournamentTx(ctx context.Context, arg database.NewTournamentParams) (models.Tournament, error) {
	var tournament models.Tournament

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		tournament, err = q.NewTournament(ctx, arg)
		if err != nil {
			return err
		}
		return err
	})
	return tournament, err
}

func (store *SQLStore) CreateTournamentStandingTx(ctx context.Context) (model.TournamentStanding, error) {
	var tournament models.Tournament

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		// Create user
		tournament, err = q.NewTournament(ctx, arg)
		if err != nil {
			return err
		}
		return err
	})
	return tournament, err
}
