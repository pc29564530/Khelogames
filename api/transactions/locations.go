package transactions

import (
	"khelogames/database"
	"khelogames/database/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (store *SQLStore) AddLocationTx(ctx *gin.Context, locationOF string, eventPublicID uuid.UUID, city, state, country string, latitude, longitude float64) (models.Locations, error) {
	var location *models.Locations

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error

		location, err = q.AddLocation(ctx, city, state, country, latitude, longitude)
		if err != nil {
			store.logger.Error("Unable to update match status: ", err)
			return err
		}

		if locationOF == "tournament" {
			_, err := q.UpdateTournamentLocation(ctx, eventPublicID, location.ID)
			if err != nil {
				store.logger.Error("Unable to update tournament location: ", err)
				return err
			}
		} else if locationOF == "team" {
			_, err := q.UpdateTeamLocation(ctx, eventPublicID, int32(location.ID))
			if err != nil {
				store.logger.Error("Unable to update tournament location: ", err)
				return err
			}
		} else if locationOF == "match" {
			_, err := q.UpdateMatchLocation(ctx, eventPublicID, int32(location.ID))
			if err != nil {
				store.logger.Error("Unable to update tournament location: ", err)
				return err
			}
		} else if locationOF == "user" {
			_, err := q.UpdateProfilesLocation(ctx, eventPublicID, int32(location.ID))
			if err != nil {
				store.logger.Error("Unable to update tournament location: ", err)
				return err
			}
		}
		return err
	})
	return *location, err
}
