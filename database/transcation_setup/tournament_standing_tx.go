package transcation_setup

import (
	"context"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/google/uuid"
)

func (store *SQLStore) FootballStandingTx(ctx context.Context, tournamentPublicID, teamPublicID uuid.UUID, groupID int32) (models.FootballStanding, error) {
	var footballStanding models.FootballStanding

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		// Create user
		footballStanding, err = q.CreateFootballStanding(ctx, tournamentPublicID, int32(groupID), teamPublicID)
		if err != nil {
			return err
		}
		return err
	})
	return footballStanding, err
}

func (store *SQLStore) CricketStandingTx(ctx context.Context, tournamentPublicID uuid.UUID, teamPublicID uuid.UUID, groupID int32) (*models.CricketStanding, error) {
	var cricketStanding *models.CricketStanding

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		// Create user
		cricketStanding, err = q.CreateCricketStanding(ctx, tournamentPublicID, int32(groupID), teamPublicID)
		if err != nil {
			return err
		}
		return err
	})
	return cricketStanding, err
}
