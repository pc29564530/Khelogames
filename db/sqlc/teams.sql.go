// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: teams.sql

package db

import (
	"context"
)

const addTeamPlayers = `-- name: AddTeamPlayers :one
INSERT INTO team_players (
    team_id,
    player_id
) VALUES ($1, $2)
RETURNING team_id, player_id
`

type AddTeamPlayersParams struct {
	TeamID   int64 `json:"team_id"`
	PlayerID int64 `json:"player_id"`
}

func (q *Queries) AddTeamPlayers(ctx context.Context, arg AddTeamPlayersParams) (TeamPlayer, error) {
	row := q.db.QueryRowContext(ctx, addTeamPlayers, arg.TeamID, arg.PlayerID)
	var i TeamPlayer
	err := row.Scan(&i.TeamID, &i.PlayerID)
	return i, err
}

const getMatchByTeam = `-- name: GetMatchByTeam :many
SELECT t.id AS tournament_id, t.tournament_name, tm.id AS match_id, tm.home_team_id, tm.away_team_id, c1.name AS home_team_name, c2.name AS away_team_name, tm.start_timestamp, t.sports
FROM matches tm
JOIN tournaments t ON tm.tournament_id = t.id
JOIN teams c1 ON tm.home_team_id = c1.id
JOIN teams c2 ON tm.away_team_id = c2.id
WHERE c1.id=$1 OR c2.id=$1
ORDER BY tm.id DESC, tm.start_timestamp DESC
`

type GetMatchByTeamRow struct {
	TournamentID   int64  `json:"tournament_id"`
	TournamentName string `json:"tournament_name"`
	MatchID        int64  `json:"match_id"`
	HomeTeamID     int64  `json:"home_team_id"`
	AwayTeamID     int64  `json:"away_team_id"`
	HomeTeamName   string `json:"home_team_name"`
	AwayTeamName   string `json:"away_team_name"`
	StartTimestamp int64  `json:"start_timestamp"`
	Sports         string `json:"sports"`
}

func (q *Queries) GetMatchByTeam(ctx context.Context, id int64) ([]GetMatchByTeamRow, error) {
	rows, err := q.db.QueryContext(ctx, getMatchByTeam, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMatchByTeamRow
	for rows.Next() {
		var i GetMatchByTeamRow
		if err := rows.Scan(
			&i.TournamentID,
			&i.TournamentName,
			&i.MatchID,
			&i.HomeTeamID,
			&i.AwayTeamID,
			&i.HomeTeamName,
			&i.AwayTeamName,
			&i.StartTimestamp,
			&i.Sports,
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

const getTeam = `-- name: GetTeam :one
SELECT id, name, slug, shortname, admin, media_url, gender, national, country, type, sports FROM teams
WHERE id=$1
`

func (q *Queries) GetTeam(ctx context.Context, id int64) (Team, error) {
	row := q.db.QueryRowContext(ctx, getTeam, id)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.Shortname,
		&i.Admin,
		&i.MediaUrl,
		&i.Gender,
		&i.National,
		&i.Country,
		&i.Type,
		&i.Sports,
	)
	return i, err
}

const getTeamPlayers = `-- name: GetTeamPlayers :many
SELECT team_id, player_id FROM team_players
WHERE team_id=$1
`

func (q *Queries) GetTeamPlayers(ctx context.Context, teamID int64) ([]TeamPlayer, error) {
	rows, err := q.db.QueryContext(ctx, getTeamPlayers, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TeamPlayer
	for rows.Next() {
		var i TeamPlayer
		if err := rows.Scan(&i.TeamID, &i.PlayerID); err != nil {
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

const getTeams = `-- name: GetTeams :many
SELECT id, name, slug, shortname, admin, media_url, gender, national, country, type, sports FROM teams
`

func (q *Queries) GetTeams(ctx context.Context) ([]Team, error) {
	rows, err := q.db.QueryContext(ctx, getTeams)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Team
	for rows.Next() {
		var i Team
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Slug,
			&i.Shortname,
			&i.Admin,
			&i.MediaUrl,
			&i.Gender,
			&i.National,
			&i.Country,
			&i.Type,
			&i.Sports,
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

const getTeamsBySport = `-- name: GetTeamsBySport :many
SELECT id, name, slug, shortname, admin, media_url, gender, national, country, type, sports FROM teams
WHERE sports=$1
`

func (q *Queries) GetTeamsBySport(ctx context.Context, sports string) ([]Team, error) {
	rows, err := q.db.QueryContext(ctx, getTeamsBySport, sports)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Team
	for rows.Next() {
		var i Team
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Slug,
			&i.Shortname,
			&i.Admin,
			&i.MediaUrl,
			&i.Gender,
			&i.National,
			&i.Country,
			&i.Type,
			&i.Sports,
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

const getTournamentsByTeam = `-- name: GetTournamentsByTeam :many
SELECT t.id, t.tournament_name, t.sports FROM tournaments t
JOIN tournament_team tt ON t.id=tt.tournament_id
JOIN teams c ON tt.team_id=c.id
WHERE c.id=$1
`

type GetTournamentsByTeamRow struct {
	ID             int64  `json:"id"`
	TournamentName string `json:"tournament_name"`
	Sports         string `json:"sports"`
}

func (q *Queries) GetTournamentsByTeam(ctx context.Context, id int64) ([]GetTournamentsByTeamRow, error) {
	rows, err := q.db.QueryContext(ctx, getTournamentsByTeam, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTournamentsByTeamRow
	for rows.Next() {
		var i GetTournamentsByTeamRow
		if err := rows.Scan(&i.ID, &i.TournamentName, &i.Sports); err != nil {
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

const newTeams = `-- name: NewTeams :one
INSERT INTO teams (
    name,
    slug,
    shortName,
    admin,
    media_url,
    gender,
    national,
    country,
    type,
    sports
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id, name, slug, shortname, admin, media_url, gender, national, country, type, sports
`

type NewTeamsParams struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Shortname string `json:"shortname"`
	Admin     string `json:"admin"`
	MediaUrl  string `json:"media_url"`
	Gender    string `json:"gender"`
	National  bool   `json:"national"`
	Country   string `json:"country"`
	Type      string `json:"type"`
	Sports    string `json:"sports"`
}

func (q *Queries) NewTeams(ctx context.Context, arg NewTeamsParams) (Team, error) {
	row := q.db.QueryRowContext(ctx, newTeams,
		arg.Name,
		arg.Slug,
		arg.Shortname,
		arg.Admin,
		arg.MediaUrl,
		arg.Gender,
		arg.National,
		arg.Country,
		arg.Type,
		arg.Sports,
	)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.Shortname,
		&i.Admin,
		&i.MediaUrl,
		&i.Gender,
		&i.National,
		&i.Country,
		&i.Type,
		&i.Sports,
	)
	return i, err
}

const searchTeam = `-- name: SearchTeam :many
SELECT id, name from teams
WHERE name LIKE $1
`

type SearchTeamRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) SearchTeam(ctx context.Context, name string) ([]SearchTeamRow, error) {
	rows, err := q.db.QueryContext(ctx, searchTeam, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchTeamRow
	for rows.Next() {
		var i SearchTeamRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

const updateMediaUrl = `-- name: UpdateMediaUrl :one
UPDATE teams
SET media_url=$1
WHERE id=$2
RETURNING id, name, slug, shortname, admin, media_url, gender, national, country, type, sports
`

type UpdateMediaUrlParams struct {
	MediaUrl string `json:"media_url"`
	ID       int64  `json:"id"`
}

func (q *Queries) UpdateMediaUrl(ctx context.Context, arg UpdateMediaUrlParams) (Team, error) {
	row := q.db.QueryRowContext(ctx, updateMediaUrl, arg.MediaUrl, arg.ID)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.Shortname,
		&i.Admin,
		&i.MediaUrl,
		&i.Gender,
		&i.National,
		&i.Country,
		&i.Type,
		&i.Sports,
	)
	return i, err
}

const updateTeamName = `-- name: UpdateTeamName :one

UPDATE teams
SET name=$1
WHERE id=$2
RETURNING id, name, slug, shortname, admin, media_url, gender, national, country, type, sports
`

type UpdateTeamNameParams struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

// -- name: UpdateTeamsSport :one
// UPDATE club
// SET sport=$1
// WHERE club_name=$2
// RETURNING *;
func (q *Queries) UpdateTeamName(ctx context.Context, arg UpdateTeamNameParams) (Team, error) {
	row := q.db.QueryRowContext(ctx, updateTeamName, arg.Name, arg.ID)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.Shortname,
		&i.Admin,
		&i.MediaUrl,
		&i.Gender,
		&i.National,
		&i.Country,
		&i.Type,
		&i.Sports,
	)
	return i, err
}