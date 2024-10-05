package database

import (
	"context"
	"khelogames/database/models"
)

const createGroupTeams = `
INSERT INTO teams_group (
    group_id,
    team_id,
    tournament_id
) VALUES ( $1, $2, $3) RETURNING id, group_id, team_id, tournament_id
`

type CreateGroupTeamsParams struct {
	GroupID      int64 `json:"group_id"`
	TeamID       int64 `json:"team_id"`
	TournamentID int64 `json:"tournament_id"`
}

func (q *Queries) CreateGroupTeams(ctx context.Context, arg CreateGroupTeamsParams) (models.TeamsGroup, error) {
	row := q.db.QueryRowContext(ctx, createGroupTeams, arg.GroupID, arg.TeamID, arg.TournamentID)
	var i models.TeamsGroup
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.TeamID,
		&i.TournamentID,
	)
	return i, err
}

const createTournamentGroup = `
INSERT INTO groups (
    name,
    tournament_id,
    strength
) VALUES ( $1, $2, $3) RETURNING id, name, tournament_id, strength
`

type CreateTournamentGroupParams struct {
	Name         string `json:"name"`
	TournamentID int64  `json:"tournament_id"`
	Strength     int32  `json:"strength"`
}

func (q *Queries) CreateTournamentGroup(ctx context.Context, arg CreateTournamentGroupParams) (models.Group, error) {
	row := q.db.QueryRowContext(ctx, createTournamentGroup, arg.Name, arg.TournamentID, arg.Strength)
	var i models.Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.TournamentID,
		&i.Strength,
	)
	return i, err
}

const getGroupTeams = `
SELECT id, group_id, team_id, tournament_id FROM teams_group
WHERE tournament_id=$1 AND group_id=$2
`

type GetGroupTeamsParams struct {
	TournamentID int64 `json:"tournament_id"`
	GroupID      int64 `json:"group_id"`
}

func (q *Queries) GetGroupTeams(ctx context.Context, arg GetGroupTeamsParams) ([]models.TeamsGroup, error) {
	rows, err := q.db.QueryContext(ctx, getGroupTeams, arg.TournamentID, arg.GroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.TeamsGroup
	for rows.Next() {
		var i models.TeamsGroup
		if err := rows.Scan(
			&i.ID,
			&i.GroupID,
			&i.TeamID,
			&i.TournamentID,
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

const getTournamentGroup = `
SELECT id, name, tournament_id, strength FROM groups
WHERE tournament_id=$1 AND id=$2
`

type GetTournamentGroupParams struct {
	TournamentID int64 `json:"tournament_id"`
	ID           int64 `json:"id"`
}

func (q *Queries) GetTournamentGroup(ctx context.Context, arg GetTournamentGroupParams) (models.Group, error) {
	row := q.db.QueryRowContext(ctx, getTournamentGroup, arg.TournamentID, arg.ID)
	var i models.Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.TournamentID,
		&i.Strength,
	)
	return i, err
}

const getTournamentGroups = `
SELECT id, name, tournament_id, strength FROM groups
WHERE tournament_id=$1
`

func (q *Queries) GetTournamentGroups(ctx context.Context, tournamentID int64) ([]models.Group, error) {
	rows, err := q.db.QueryContext(ctx, getTournamentGroups, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Group
	for rows.Next() {
		var i models.Group
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.TournamentID,
			&i.Strength,
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
