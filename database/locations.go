package database

import (
	"context"
	"fmt"
	"khelogames/database/models"
)

const addLocation = `
	INSERT INTO locations (
		city,
		state,
		country,
		latitude,
		longitude,
		created_at,
		updated_at,
		h3_index
	) VALUES (
	 	$1, $2, $3, $4, $5, NOW(), NOW(), $6
	) RETURNING *;
`

func (q *Queries) AddLocation(ctx context.Context, city, state, country string, latitude, longitude float64, h3Index string) (*models.Locations, error) {
	row := q.db.QueryRowContext(ctx, addLocation, city, state, country, latitude, longitude, h3Index)
	i := &models.Locations{}
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.City,
		&i.State,
		&i.Country,
		&i.Latitude,
		&i.Longitude,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.H3Index,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to get query row context: ", err)
	}
	return i, err
}

const updateUserLocation = `
	UPDATE locations
	SET
		latitude = $2,
		longitude = $3,
		h3_index = $4
	WHERE id = $1
	RETURNING *
`

func (q *Queries) UpdateUserLocation(ctx context.Context, locationID int32, latitude, longitude float64, h3Index string) (*models.Locations, error) {
	rows := q.db.QueryRowContext(ctx, updateUserLocation, locationID, latitude, longitude, h3Index)
	var i models.Locations
	err := rows.Scan(
		&i.ID,
		&i.PublicID,
		&i.City,
		&i.State,
		&i.Country,
		&i.Latitude,
		&i.Longitude,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.H3Index,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to scan location by user: ", err)
	}

	return &i, err
}

const updateLocationWithAddress = `
	UPDATE locations
	SET
		city = $2,
		state = $3,
		country = $4,
		latitude = $5,
		longitude = $6,
		h3_index = $7,
		updated_at = NOW()
	WHERE id = $1
	RETURNING *
`

func (q *Queries) UpdateLocationWithAddress(ctx context.Context, locationID int32, city, state, country string, latitude, longitude float64, h3Index string) (*models.Locations, error) {
	rows := q.db.QueryRowContext(ctx, updateLocationWithAddress, locationID, city, state, country, latitude, longitude, h3Index)
	var i models.Locations
	err := rows.Scan(
		&i.ID,
		&i.PublicID,
		&i.City,
		&i.State,
		&i.Country,
		&i.Latitude,
		&i.Longitude,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.H3Index,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to scan location: %w", err)
	}

	return &i, err
}
