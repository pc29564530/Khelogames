// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: tournament_organization.sql

package db

import (
	"context"
	"time"
)

const createTournamentOrganization = `-- name: CreateTournamentOrganization :one
INSERT INTO tournament_organization (
    tournament_id,
    tournament_start,
    player_count,
    team_count,
    group_count,
    advanced_team
) VALUES 
($1, $2, $3, $4, $5, $6 )
RETURNING id, tournament_id, tournament_start, player_count, team_count, group_count, advanced_team
`

type CreateTournamentOrganizationParams struct {
	TournamentID    int64     `json:"tournament_id"`
	TournamentStart time.Time `json:"tournament_start"`
	PlayerCount     int64     `json:"player_count"`
	TeamCount       int64     `json:"team_count"`
	GroupCount      int64     `json:"group_count"`
	AdvancedTeam    int64     `json:"advanced_team"`
}

func (q *Queries) CreateTournamentOrganization(ctx context.Context, arg CreateTournamentOrganizationParams) (TournamentOrganization, error) {
	row := q.db.QueryRowContext(ctx, createTournamentOrganization,
		arg.TournamentID,
		arg.TournamentStart,
		arg.PlayerCount,
		arg.TeamCount,
		arg.GroupCount,
		arg.AdvancedTeam,
	)
	var i TournamentOrganization
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.TournamentStart,
		&i.PlayerCount,
		&i.TeamCount,
		&i.GroupCount,
		&i.AdvancedTeam,
	)
	return i, err
}

const getTournamentOrganization = `-- name: GetTournamentOrganization :one
SELECT id, tournament_id, tournament_start, player_count, team_count, group_count, advanced_team FROM tournament_organization
WHERE tournament_id=$1
`

func (q *Queries) GetTournamentOrganization(ctx context.Context, tournamentID int64) (TournamentOrganization, error) {
	row := q.db.QueryRowContext(ctx, getTournamentOrganization, tournamentID)
	var i TournamentOrganization
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.TournamentStart,
		&i.PlayerCount,
		&i.TeamCount,
		&i.GroupCount,
		&i.AdvancedTeam,
	)
	return i, err
}