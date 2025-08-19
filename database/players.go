package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const getPlayerByProfilePublicID = `
SELECT p.id, p.public_id, p.user_id, p.game_id, p.name, p.slug, p.short_name, p.media_url, p.positions, p.country, p.created_at, p.updated_at
FROM players p
JOIN user_profiles AS up ON up.user_id = p.user_id
JOIN users AS u ON u.id = p.user_id
WHERE up.public_id=$1
`

func (q *Queries) GetPlayerByProfile(ctx context.Context, profilePublicID uuid.UUID) (*models.Player, error) {
	row := q.db.QueryRowContext(ctx, getPlayerByProfilePublicID, profilePublicID)
	var i models.Player
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.GameID,
		&i.Name,
		&i.Slug,
		&i.ShortName,
		&i.MediaUrl,
		&i.Positions,
		&i.Country,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Failed to find the rows: ", err)
			return nil, nil
		}
		return nil, err
	}
	return &i, err
}

const getAllPlayer = `
SELECT * FROM players
`

func (q *Queries) GetAllPlayer(ctx context.Context) ([]models.Player, error) {
	rows, err := q.db.QueryContext(ctx, getAllPlayer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Player
	for rows.Next() {
		var i models.Player
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.UserID,
			&i.GameID,
			&i.Name,
			&i.Slug,
			&i.ShortName,
			&i.MediaUrl,
			&i.Positions,
			&i.Country,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// Existing method - gets player by public ID
const getPlayer = `
SELECT p.id, p.public_id, p.user_id, p.game_id, p.name, p.slug, p.short_name, p.media_url, p.positions, p.country, p.created_at, p.updated_at
FROM players p
WHERE p.public_id=$1
`

func (q *Queries) GetPlayer(ctx context.Context, userPublicID uuid.UUID) (*models.Player, error) {
	row := q.db.QueryRowContext(ctx, getPlayer, userPublicID)
	var i models.Player
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.GameID,
		&i.Name,
		&i.Slug,
		&i.ShortName,
		&i.MediaUrl,
		&i.Positions,
		&i.Country,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, err
}

// New method - gets player by player's own public ID
const getPlayerByPublicID = `
SELECT *
FROM players p
WHERE p.public_id=$1
`

func (q *Queries) GetPlayerByPublicID(ctx context.Context, playerPublicID uuid.UUID) (*models.Player, error) {
	row := q.db.QueryRowContext(ctx, getPlayerByPublicID, playerPublicID)
	var i models.Player
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.GameID,
		&i.Name,
		&i.Slug,
		&i.ShortName,
		&i.MediaUrl,
		&i.Positions,
		&i.Country,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, err
}

// Method to get player by internal ID (for cases where you have the int64 ID)
const getPlayerByID = `
SELECT p.id, p.public_id, p.user_id, p.name, p.slug, p.short_name, p.media_url, p.positions, p.country, p.game_id
FROM players p
WHERE p.id=$1
`

func (q *Queries) GetPlayerByID(ctx context.Context, playerID int64) (*models.Player, error) {
	row := q.db.QueryRowContext(ctx, getPlayerByID, playerID)
	var i models.Player
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.Name,
		&i.Slug,
		&i.ShortName,
		&i.MediaUrl,
		&i.Positions,
		&i.Country,
		&i.GameID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &i, err
}

const getPlayersCountry = `
SELECT * FROM players
WHERE country=$1
`

func (q *Queries) GetPlayersCountry(ctx context.Context, country string) ([]models.Player, error) {
	rows, err := q.db.QueryContext(ctx, getPlayersCountry, country)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Player
	for rows.Next() {
		var i models.Player
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.UserID,
			&i.Name,
			&i.Slug,
			&i.ShortName,
			&i.MediaUrl,
			&i.Positions,
			&i.Country,
			&i.GameID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPlayerBySport = `
	SELECT * FROM players
	WHERE game_id=$1;
`

func (q *Queries) GetPlayersBySport(ctx context.Context, gameID int32) ([]models.Player, error) {
	rows, err := q.db.QueryContext(ctx, getPlayerBySport, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Player
	for rows.Next() {
		var i models.Player
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.UserID,
			&i.GameID,
			&i.Name,
			&i.Slug,
			&i.ShortName,
			&i.MediaUrl,
			&i.Positions,
			&i.Country,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const newPlayer = `
WITH userID AS (
	SELECT * FROM users
	WHERE public_id = $1
)
INSERT INTO players (
	user_id,
    game_id,
    name,
    slug,
    short_name,
    media_url,
    positions,
    country
)
SELECT 
	userID.id,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8
FROM userID
RETURNING *
`

type NewPlayerParams struct {
	UserPublicID uuid.UUID `json:"userPublicID"`
	GameID       int64     `json:"game_id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	ShortName    string    `json:"short_name"`
	MediaUrl     string    `json:"media_url"`
	Positions    string    `json:"positions"`
	Country      string    `json:"country"`
}

func (q *Queries) NewPlayer(ctx context.Context, arg NewPlayerParams) (models.Player, error) {
	row := q.db.QueryRowContext(ctx, newPlayer,
		arg.UserPublicID,
		arg.GameID,
		arg.Name,
		arg.Slug,
		arg.ShortName,
		arg.MediaUrl,
		arg.Positions,
		arg.Country,
	)
	var i models.Player
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.GameID,
		&i.Name,
		&i.Slug,
		&i.ShortName,
		&i.MediaUrl,
		&i.Positions,
		&i.Country,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const searchPlayer = `
SELECT * FROM players
WHERE name LIKE $1
`

func (q *Queries) SearchPlayer(ctx context.Context, playerName string) ([]models.Player, error) {
	rows, err := q.db.QueryContext(ctx, searchPlayer, playerName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Player
	for rows.Next() {
		var i models.Player
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.UserID,
			&i.GameID,
			&i.Name,
			&i.Slug,
			&i.ShortName,
			&i.MediaUrl,
			&i.Positions,
			&i.Country,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePlayerMedia = `
UPDATE players
SET media_url=$2
WHERE public_id=$1
RETURNING *
`

type UpdatePlayerMediaParams struct {
	PublicID uuid.UUID `json:"public_id"`
	MediaUrl string    `json:"media_url"`
}

func (q *Queries) UpdatePlayerMedia(ctx context.Context, publicID uuid.UUID, mediaUrl string) (*models.Player, error) {
	row := q.db.QueryRowContext(ctx, updatePlayerMedia, publicID, mediaUrl)
	var i models.Player
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.GameID,
		&i.Name,
		&i.Slug,
		&i.ShortName,
		&i.MediaUrl,
		&i.Positions,
		&i.Country,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return &i, err
}

const updatePlayerPosition = `
UPDATE players
SET positions=$2
WHERE public_id=$1
RETURNING *
`

type UpdatePlayerPositionParams struct {
	PublicID  uuid.UUID `json:"public_id"`
	Positions string    `json:"positions"`
}

func (q *Queries) UpdatePlayerPosition(ctx context.Context, publicID uuid.UUID, positions string) (*models.Player, error) {
	row := q.db.QueryRowContext(ctx, updatePlayerPosition, publicID, positions)
	var i models.Player
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.GameID,
		&i.Name,
		&i.Slug,
		&i.ShortName,
		&i.MediaUrl,
		&i.Positions,
		&i.Country,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return &i, err
}
