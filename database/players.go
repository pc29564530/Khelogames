package database

import (
	"context"
	"database/sql"
	"khelogames/database/models"
)

const getAllPlayer = `
SELECT id, username, slug, short_name, media_url, positions, sports, country, player_name, game_id FROM players
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
			&i.Username,
			&i.Slug,
			&i.ShortName,
			&i.MediaUrl,
			&i.Positions,
			&i.Sports,
			&i.Country,
			&i.PlayerName,
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

const getPlayer = `
SELECT id, username, slug, short_name, media_url, positions, sports, country, player_name, game_id FROM players
WHERE id=$1
`

func (q *Queries) GetPlayer(ctx context.Context, id int64) (*models.Player, error) {
	row := q.db.QueryRowContext(ctx, getPlayer, id)
	var i models.Player
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Slug,
		&i.ShortName,
		&i.MediaUrl,
		&i.Positions,
		&i.Sports,
		&i.Country,
		&i.PlayerName,
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
SELECT id, username, slug, short_name, media_url, positions, sports, country, player_name, game_id FROM players
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
			&i.Username,
			&i.Slug,
			&i.ShortName,
			&i.MediaUrl,
			&i.Positions,
			&i.Sports,
			&i.Country,
			&i.PlayerName,
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

const newPlayer = `
INSERT INTO players (
    username,
    slug,
    short_name,
    media_url,
    positions,
    sports,
    country,
    player_name,
    game_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id, username, slug, short_name, media_url, positions, sports, country, player_name, game_id
`

type NewPlayerParams struct {
	Username   string `json:"username"`
	Slug       string `json:"slug"`
	ShortName  string `json:"short_name"`
	MediaUrl   string `json:"media_url"`
	Positions  string `json:"positions"`
	Sports     string `json:"sports"`
	Country    string `json:"country"`
	PlayerName string `json:"player_name"`
	GameID     int64  `json:"game_id"`
}

func (q *Queries) NewPlayer(ctx context.Context, arg NewPlayerParams) (models.Player, error) {
	row := q.db.QueryRowContext(ctx, newPlayer,
		arg.Username,
		arg.Slug,
		arg.ShortName,
		arg.MediaUrl,
		arg.Positions,
		arg.Sports,
		arg.Country,
		arg.PlayerName,
		arg.GameID,
	)
	var i models.Player
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Slug,
		&i.ShortName,
		&i.MediaUrl,
		&i.Positions,
		&i.Sports,
		&i.Country,
		&i.PlayerName,
		&i.GameID,
	)
	return i, err
}

const searchPlayer = `
SELECT id, username, slug, short_name, media_url, positions, sports, country, player_name, game_id FROM players
WHERE player_name LIKE $1
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
			&i.Username,
			&i.Slug,
			&i.ShortName,
			&i.MediaUrl,
			&i.Positions,
			&i.Sports,
			&i.Country,
			&i.PlayerName,
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

const updatePlayerMedia = `
UPDATE players
SET media_url=$1
WHERE id=$2
RETURNING id, username, slug, short_name, media_url, positions, sports, country, player_name, game_id
`

type UpdatePlayerMediaParams struct {
	MediaUrl string `json:"media_url"`
	ID       int64  `json:"id"`
}

func (q *Queries) UpdatePlayerMedia(ctx context.Context, arg UpdatePlayerMediaParams) (models.Player, error) {
	row := q.db.QueryRowContext(ctx, updatePlayerMedia, arg.MediaUrl, arg.ID)
	var i models.Player
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Slug,
		&i.ShortName,
		&i.MediaUrl,
		&i.Positions,
		&i.Sports,
		&i.Country,
		&i.PlayerName,
		&i.GameID,
	)
	return i, err
}

const updatePlayerPosition = `
UPDATE players
SET positions=$1
WHERE id=$2
RETURNING id, username, slug, short_name, media_url, positions, sports, country, player_name, game_id
`

type UpdatePlayerPositionParams struct {
	Positions string `json:"positions"`
	ID        int64  `json:"id"`
}

func (q *Queries) UpdatePlayerPosition(ctx context.Context, arg UpdatePlayerPositionParams) (models.Player, error) {
	row := q.db.QueryRowContext(ctx, updatePlayerPosition, arg.Positions, arg.ID)
	var i models.Player
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Slug,
		&i.ShortName,
		&i.MediaUrl,
		&i.Positions,
		&i.Sports,
		&i.Country,
		&i.PlayerName,
		&i.GameID,
	)
	return i, err
}
