package new_db

import (
	"context"
	"encoding/json"
	"khelogames/new_db/models"
)

const getTournament = `-- name: GetTournament :one
SELECT id, tournament_name, slug, sports, country, status_code, level, start_timestamp, game_id FROM tournaments
WHERE id=$1
`

func (q *Queries) GetTournament(ctx context.Context, id int64) (models.Tournament, error) {
	row := q.db.QueryRowContext(ctx, getTournament, id)
	var i models.Tournament
	err := row.Scan(
		&i.ID,
		&i.TournamentName,
		&i.Slug,
		&i.Sports,
		&i.Country,
		&i.StatusCode,
		&i.Level,
		&i.StartTimestamp,
		&i.GameID,
	)
	return i, err
}

const getTournaments = `
SELECT id, tournament_name, slug, sports, country, status_code, level, start_timestamp, game_id FROM tournaments
`

func (q *Queries) GetTournaments(ctx context.Context) ([]models.Tournament, error) {
	rows, err := q.db.QueryContext(ctx, getTournaments)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Tournament
	for rows.Next() {
		var i models.Tournament
		if err := rows.Scan(
			&i.ID,
			&i.TournamentName,
			&i.Slug,
			&i.Sports,
			&i.Country,
			&i.StatusCode,
			&i.Level,
			&i.StartTimestamp,
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

const getTournamentsByLevel = `
SELECT id, tournament_name, slug, sports, country, status_code, level, start_timestamp, game_id FROM tournaments
WHERE sports=$1 AND level=$2
`

type GetTournamentsByLevelParams struct {
	Sports string `json:"sports"`
	Level  string `json:"level"`
}

func (q *Queries) GetTournamentsByLevel(ctx context.Context, arg GetTournamentsByLevelParams) ([]models.Tournament, error) {
	rows, err := q.db.QueryContext(ctx, getTournamentsByLevel, arg.Sports, arg.Level)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Tournament
	for rows.Next() {
		var i models.Tournament
		if err := rows.Scan(
			&i.ID,
			&i.TournamentName,
			&i.Slug,
			&i.Sports,
			&i.Country,
			&i.StatusCode,
			&i.Level,
			&i.StartTimestamp,
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

const getTournamentsBySport = `
SELECT 
    g.id, g.name, g.min_players, JSON_BUILD_OBJECT('id', t.id, 'name', t.name, 'slug', t.slug, 'country', t.country, 'status_code', t.status_code, 'level', t.level, 'start_timestamp', t.start_timestamp, 'game_id', t.game_id) AS tournament_data
FROM tournaments t
JOIN games AS g ON g.id = t.game_id
WHERE t.game_id=$1
`

type GetTournamentsBySportRow struct {
	ID             int64           `json:"id"`
	Name           string          `json:"name"`
	MinPlayers     int32           `json:"min_players"`
	TournamentData json.RawMessage `json:"tournament_data"`
}

func (q *Queries) GetTournamentsBySport(ctx context.Context, gameID int64) ([]GetTournamentsBySportRow, error) {
	rows, err := q.db.QueryContext(ctx, getTournamentsBySport, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTournamentsBySportRow
	for rows.Next() {
		var i GetTournamentsBySportRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.MinPlayers,
			&i.TournamentData,
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

const newTournament = `
INSERT INTO tournaments (
    tournament_name,
    slug,
    sports,
    country,
    status_code,
    level,
    start_timestamp,
    game_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, tournament_name, slug, sports, country, status_code, level, start_timestamp, game_id
`

type NewTournamentParams struct {
	TournamentName string `json:"tournament_name"`
	Slug           string `json:"slug"`
	Sports         string `json:"sports"`
	Country        string `json:"country"`
	StatusCode     string `json:"status_code"`
	Level          string `json:"level"`
	StartTimestamp int64  `json:"start_timestamp"`
	GameID         int64  `json:"game_id"`
}

func (q *Queries) NewTournament(ctx context.Context, arg NewTournamentParams) (models.Tournament, error) {
	row := q.db.QueryRowContext(ctx, newTournament,
		arg.TournamentName,
		arg.Slug,
		arg.Sports,
		arg.Country,
		arg.StatusCode,
		arg.Level,
		arg.StartTimestamp,
		arg.GameID,
	)
	var i models.Tournament
	err := row.Scan(
		&i.ID,
		&i.TournamentName,
		&i.Slug,
		&i.Sports,
		&i.Country,
		&i.StatusCode,
		&i.Level,
		&i.StartTimestamp,
		&i.GameID,
	)
	return i, err
}

const updateTournamentDate = `
UPDATE tournaments
SET start_timestamp=$1
WHERE id=$2
RETURNING id, tournament_name, slug, sports, country, status_code, level, start_timestamp, game_id
`

type UpdateTournamentDateParams struct {
	StartTimestamp int64 `json:"start_timestamp"`
	ID             int64 `json:"id"`
}

func (q *Queries) UpdateTournamentDate(ctx context.Context, arg UpdateTournamentDateParams) (models.Tournament, error) {
	row := q.db.QueryRowContext(ctx, updateTournamentDate, arg.StartTimestamp, arg.ID)
	var i models.Tournament
	err := row.Scan(
		&i.ID,
		&i.TournamentName,
		&i.Slug,
		&i.Sports,
		&i.Country,
		&i.StatusCode,
		&i.Level,
		&i.StartTimestamp,
		&i.GameID,
	)
	return i, err
}

const updateTournamentStatus = `
UPDATE tournaments
SET status_code=$1
WHERE id=$2
RETURNING id, tournament_name, slug, sports, country, status_code, level, start_timestamp, game_id
`

type UpdateTournamentStatusParams struct {
	StatusCode string `json:"status_code"`
	ID         int64  `json:"id"`
}

func (q *Queries) UpdateTournamentStatus(ctx context.Context, arg UpdateTournamentStatusParams) (models.Tournament, error) {
	row := q.db.QueryRowContext(ctx, updateTournamentStatus, arg.StatusCode, arg.ID)
	var i models.Tournament
	err := row.Scan(
		&i.ID,
		&i.TournamentName,
		&i.Slug,
		&i.Sports,
		&i.Country,
		&i.StatusCode,
		&i.Level,
		&i.StartTimestamp,
		&i.GameID,
	)
	return i, err
}
