package transactions

import (
	"context"
	"fmt"
	"khelogames/database"
	"khelogames/database/models"

	"github.com/google/uuid"
	"github.com/uber/h3-go/v4"
)

func (store *SQLStore) CreateTeamsTx(ctx context.Context, userPublicID uuid.UUID, name, slug, shortName, mediaUrl, gender string,
	national bool,
	types string,
	playerCount int32,
	gameID int32,
	city, state, country string,
	latitude, longitude float64,
) (models.Team, error) {
	var team models.Team

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		// Create user

		latLng := h3.NewLatLng(latitude, longitude)
		cell, err := h3.LatLngToCell(latLng, 9)
		if err != nil {
			store.logger.Error("Unable to get cell of h3: ", err)
			return err
		}

		h3Index := cell.String()

		location, err := q.AddLocation(ctx,
			city,
			state,
			country,
			latitude,
			longitude,
			h3Index,
		)
		if err != nil {
			return fmt.Errorf("Failed to add the location: ", err)
		}

		locationID := int32(location.ID)

		arg := database.NewTeamsParams{
			UserPublicID: userPublicID,
			Name:         name,
			Slug:         slug,
			Shortname:    shortName,
			MediaUrl:     mediaUrl,
			Gender:       gender,
			National:     false,
			Country:      country,
			Type:         types,
			PlayerCount:  playerCount,
			GameID:       gameID,
			LocationID:   &locationID,
		}

		team, err = q.NewTeams(ctx, arg)
		if err != nil {
			return err
		}
		return err
	})
	return team, err
}

func (store *SQLStore) UpdateTeamTx(ctx context.Context, teamPublicID uuid.UUID, city, state, country string, latitude, longitude float64) (models.Team, error) {
	var team models.Team

	err := store.execTx(ctx, func(q *database.Queries) error {
		var err error
		// Create user

		latLng := h3.NewLatLng(latitude, longitude)
		cell, err := h3.LatLngToCell(latLng, 9)
		if err != nil {
			store.logger.Error("Unable to get cell of h3: ", err)
			return err
		}

		h3Index := cell.String()

		location, err := q.AddLocation(ctx,
			city,
			state,
			country,
			latitude,
			longitude,
			h3Index,
		)
		if err != nil {
			return fmt.Errorf("Failed to add the location: ", err)
		}

		locationID := int32(location.ID)

		res, err := store.UpdateTeamLocation(ctx, teamPublicID, locationID)
		if err != nil {
			return fmt.Errorf("Failed to update team location: ", err)
		}
		team = *res
		return err
	})
	return team, err
}
