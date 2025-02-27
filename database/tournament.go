package database

import (
	"context"
	"khelogames/database/models"
)

const getTournament = `
SELECT id, name, slug, sports, country, status_code, level, start_timestamp, game_id, group_count, max_group_team, stage, has_knockout FROM tournaments
WHERE id=$1
`

func (q *Queries) GetTournament(ctx context.Context, id int64) (models.Tournament, error) {
	row := q.db.QueryRowContext(ctx, getTournament, id)
	var i models.Tournament
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.Sports,
		&i.Country,
		&i.StatusCode,
		&i.Level,
		&i.StartTimestamp,
		&i.GameID,
		&i.GroupCount,
		&i.MaxGroupTeam,
		&i.Stage,
	)
	return i, err
}

const getTournaments = `
SELECT * FROM tournaments
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
			&i.Name,
			&i.Slug,
			&i.Sports,
			&i.Country,
			&i.StatusCode,
			&i.Level,
			&i.StartTimestamp,
			&i.GameID,
			&i.GroupCount,
			&i.MaxGroupTeam,
			&i.Stage,
			&i.HasKnockout,
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
SELECT * FROM tournaments
WHERE game_id=$1 AND level=$2
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
			&i.Name,
			&i.Slug,
			&i.Sports,
			&i.Country,
			&i.StatusCode,
			&i.Level,
			&i.StartTimestamp,
			&i.GameID,
			&i.GroupCount,
			&i.MaxGroupTeam,
			&i.Stage,
			&i.HasKnockout,
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
    g.id, g.name, g.min_players, JSON_BUILD_OBJECT('id', t.id, 'name', t.name, 'slug', t.slug, 'country', t.country, 'status_code', t.status_code, 'level', t.level, 'start_timestamp', t.start_timestamp, 'game_id', t.game_id, 'group_count', t.group_count, 'max_group_team', t.max_group_team, 'stage', t.stage, 'has_knockout', t.has_knockout) AS tournament_data
FROM tournaments t
JOIN games AS g ON g.id = t.game_id
WHERE t.game_id=$1
`

type GetTournamentsBySportRow struct {
	ID         int64       `json:"id"`
	Name       string      `json:"name"`
	MinPlayers int32       `json:"min_players"`
	Tournament interface{} `json:"tournament_data"`
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
			&i.Tournament,
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
    name,
    slug,
    sports,
    country,
    status_code,
    level,
    start_timestamp,
    game_id,
	group_count,
	max_group_team,
	stage,
	has_knockout
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id, name, slug, sports, country, status_code, level, start_timestamp, game_id, group_count, max_group_team, stage, has_knockout
`

type NewTournamentParams struct {
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Sports         string `json:"sports"`
	Country        string `json:"country"`
	StatusCode     string `json:"status_code"`
	Level          string `json:"level"`
	StartTimestamp int64  `json:"start_timestamp"`
	GameID         *int64 `json:"game_id"`
	GroupCount     *int32 `json:"group_count"`
	MaxGroupTeam   *int32 `json:"max_group_team"`
	Stage          string `json:"stage"`
	HasKnockout    bool   `json:"has_knockout"`
}

func (q *Queries) NewTournament(ctx context.Context, arg NewTournamentParams) (models.Tournament, error) {
	row := q.db.QueryRowContext(ctx, newTournament,
		arg.Name,
		arg.Slug,
		arg.Sports,
		arg.Country,
		arg.StatusCode,
		arg.Level,
		arg.StartTimestamp,
		arg.GameID,
		arg.GroupCount,
		arg.MaxGroupTeam,
		arg.Stage,
		arg.HasKnockout,
	)
	var i models.Tournament
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.Sports,
		&i.Country,
		&i.StatusCode,
		&i.Level,
		&i.StartTimestamp,
		&i.GameID,
		&i.GroupCount,
		&i.MaxGroupTeam,
		&i.Stage,
		&i.HasKnockout,
	)
	return i, err
}

const updateTournamentDate = `
UPDATE tournaments
SET start_timestamp=$1
WHERE id=$2
RETURNING *
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
		&i.Name,
		&i.Slug,
		&i.Sports,
		&i.Country,
		&i.StatusCode,
		&i.Level,
		&i.StartTimestamp,
		&i.GameID,
		&i.GroupCount,
		&i.MaxGroupTeam,
		&i.Stage,
		&i.HasKnockout,
	)
	return i, err
}

const updateTournamentStatus = `
UPDATE tournaments
SET status_code=$1
WHERE id=$2
RETURNING *
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
		&i.Name,
		&i.Slug,
		&i.Sports,
		&i.Country,
		&i.StatusCode,
		&i.Level,
		&i.StartTimestamp,
		&i.GameID,
		&i.GroupCount,
		&i.MaxGroupTeam,
		&i.Stage,
		&i.HasKnockout,
	)
	return i, err
}
