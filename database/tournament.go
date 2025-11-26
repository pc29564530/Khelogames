package database

import (
	"context"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const getTournament = `
SELECT * FROM tournaments
WHERE public_id=$1
`

func (q *Queries) GetTournament(ctx context.Context, publicID uuid.UUID) (models.Tournament, error) {
	row := q.db.QueryRowContext(ctx, getTournament, publicID)
	var i models.Tournament
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.GameID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.Country,
		&i.Status,
		&i.Season,
		&i.Level,
		&i.StartTimestamp,
		&i.GroupCount,
		&i.MaxGroupTeam,
		&i.Stage,
		&i.HasKnockout,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTournamentByID = `
SELECT * FROM tournaments
WHERE id=$1
`

func (q *Queries) GetTournamentByID(ctx context.Context, id int64) (models.Tournament, error) {
	row := q.db.QueryRowContext(ctx, getTournament, id)
	var i models.Tournament
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.GameID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.Country,
		&i.Status,
		&i.Season,
		&i.Level,
		&i.StartTimestamp,
		&i.GroupCount,
		&i.MaxGroupTeam,
		&i.Stage,
		&i.HasKnockout,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
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
			&i.PublicID,
			&i.UserID,
			&i.GameID,
			&i.Name,
			&i.Slug,
			&i.Description,
			&i.Country,
			&i.Status,
			&i.Season,
			&i.Level,
			&i.StartTimestamp,
			&i.GroupCount,
			&i.MaxGroupTeam,
			&i.Stage,
			&i.HasKnockout,
			&i.IsPublic,
			&i.CreatedAt,
			&i.UpdatedAt,
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

func (q *Queries) GetTournamentsByLevel(ctx context.Context, gameID int64, level string) ([]models.Tournament, error) {
	rows, err := q.db.QueryContext(ctx, getTournamentsByLevel, gameID, level)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Tournament
	for rows.Next() {
		var i models.Tournament
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.UserID,
			&i.GameID,
			&i.Name,
			&i.Slug,
			&i.Description,
			&i.Country,
			&i.Status,
			&i.Season,
			&i.Level,
			&i.StartTimestamp,
			&i.GroupCount,
			&i.MaxGroupTeam,
			&i.Stage,
			&i.HasKnockout,
			&i.IsPublic,
			&i.CreatedAt,
			&i.UpdatedAt,
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
    g.id, g.name, g.min_players, JSON_BUILD_OBJECT('id', t.id, 'public_id', t.public_id, 'user_id', t.user_id, 'name', t.name, 'slug', t.slug, 'country', t.country, 'status', t.status, 'level', t.level, 'start_timestamp', t.start_timestamp, 'game_id', t.game_id, 'group_count', t.group_count, 'max_group_team', t.max_group_teams, 'stage', t.stage, 'has_knockout', t.has_knockout) AS tournament_data
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
WITH userID AS (
	SELECT * FROM users WHERE public_id=$1
)
INSERT INTO tournaments (
	user_id,
    name,
    slug,
	description,
    country,
    status,
	season,
    level,
    start_timestamp,
    game_id,
	group_count,
	max_group_teams,
	stage,
	has_knockout,
	is_public
) 
SELECT 
	userID.id,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9,
	$10,
	$11,
	$12,
	$13,
	$14,
	$15
FROM userID	
RETURNING *
`

type NewTournamentParams struct {
	UserPublicID   uuid.UUID `json:"user_public_id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	Description    string    `json:"description"`
	Country        string    `json:"country"`
	Status         string    `json:"status"`
	Season         int       `json:"season"`
	Level          string    `json:"level"`
	StartTimestamp int64     `json:"start_timestamp"`
	GameID         *int64    `json:"game_id"`
	GroupCount     *int32    `json:"group_count"`
	MaxGroupTeams  *int32    `json:"max_group_teams"`
	Stage          string    `json:"stage"`
	HasKnockout    bool      `json:"has_knockout"`
	IsPublic       bool      `json:"is_public"`
}

func (q *Queries) NewTournament(ctx context.Context, arg NewTournamentParams) (models.Tournament, error) {
	row := q.db.QueryRowContext(ctx, newTournament,
		arg.UserPublicID,
		arg.Name,
		arg.Slug,
		arg.Description,
		arg.Country,
		arg.Status,
		arg.Season,
		arg.Level,
		arg.StartTimestamp,
		arg.GameID,
		arg.GroupCount,
		arg.MaxGroupTeams,
		arg.Stage,
		arg.HasKnockout,
		arg.IsPublic,
	)
	var i models.Tournament
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.GameID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.Country,
		&i.Status,
		&i.Season,
		&i.Level,
		&i.StartTimestamp,
		&i.GroupCount,
		&i.MaxGroupTeam,
		&i.Stage,
		&i.HasKnockout,
		&i.IsPublic,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateTournamentDate = `
UPDATE tournaments
SET start_timestamp=$2
WHERE public_id=$1
RETURNING *
`

type UpdateTournamentDateParams struct {
	TournamentPublicID uuid.UUID `json:"tournament_public_id"`
	StartTimestamp     int64     `json:"start_timestamp"`
}

func (q *Queries) UpdateTournamentDate(ctx context.Context, arg UpdateTournamentDateParams) (models.Tournament, error) {
	row := q.db.QueryRowContext(ctx, updateTournamentDate, arg.TournamentPublicID, arg.StartTimestamp)
	var i models.Tournament
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.GameID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.Country,
		&i.Status,
		&i.Season,
		&i.Level,
		&i.StartTimestamp,
		&i.GroupCount,
		&i.MaxGroupTeam,
		&i.Stage,
		&i.HasKnockout,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateTournamentStatus = `
UPDATE tournaments
SET status_code=$2
WHERE public_id=$1
RETURNING *
`

type UpdateTournamentStatusParams struct {
	TournamentPublicID uuid.UUID `json:"tournament_public_id"`
	Status             string    `json:"status_code"`
}

func (q *Queries) UpdateTournamentStatus(ctx context.Context, arg UpdateTournamentStatusParams) (models.Tournament, error) {
	row := q.db.QueryRowContext(ctx, updateTournamentStatus, arg.TournamentPublicID, arg.Status)
	var i models.Tournament
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.GameID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.Country,
		&i.Status,
		&i.Season,
		&i.Level,
		&i.StartTimestamp,
		&i.GroupCount,
		&i.MaxGroupTeam,
		&i.Stage,
		&i.HasKnockout,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const addTournamentUserRoles = `
INSERT INTO tournament_user_roles (
    tournament_id,
    user_id,
    role
) VALUES ($1, $2, $3)
RETURNING *;
`

func (q *Queries) AddTournamentUserRoles(ctx context.Context, tournamentID, userID int32, role string) (*models.TournamentUserRoles, error) {
	var i models.TournamentUserRoles
	rows := q.db.QueryRowContext(ctx, addTournamentUserRoles, tournamentID, userID, role)
	err := rows.Scan(
		&i.ID,
		&i.TournamentID,
		&i.UserID,
		&i.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to scan tournament user role query: ", err)
	}
	return &i, nil
}

const getTournamentUserRoles = `
    SELECT EXISTS(
        SELECT 1 
        FROM tournament_user_roles
        WHERE tournament_id = $1 AND user_id = $2
    );
`

func (q *Queries) GetTournamentUserRole(ctx context.Context, tournamentID, userID int32) (bool, error) {
	var exists bool

	err := q.db.QueryRowContext(ctx, getTournamentUserRoles, tournamentID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check tournament user role: %w", err)
	}

	return exists, nil
}
