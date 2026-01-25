package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const getTournamentTeam = `
SELECT tt.tournament_id, JSON_BUILD_OBJECT('id', tm.id, 'public_id', tm.public_id, 'user_id', tm.user_id, 'name', tm.name, 'slug', tm.slug, 'short_name', tm.shortname, 'admin', tm.admin, 'media_url', tm.media_url, 'gender', tm.gender, 'national', tm.national, 'country', tm.country, 'type', tm.type, 'player_count', tm.player_count, 'game_id', tm.game_id) AS team_data
FROM tournament_team tt
LEFT JOIN teams AS tm ON tm.id = tt.team_id
WHERE tm.public_id=$1 AND tt.public_id=$2
`

type GetTournamentTeamInterface struct {
	TournamentID int64  `json:"tournament_id"`
	TeamData     string `json:"team_data"`
}

func (q *Queries) GetTournamentTeam(ctx context.Context, teamPublicID, tournamentPublicID uuid.UUID) (GetTournamentTeamInterface, error) {
	row := q.db.QueryRowContext(ctx, getTournamentTeam, teamPublicID, tournamentPublicID)
	var i GetTournamentTeamInterface
	err := row.Scan(&i.TournamentID, &i.TeamData)
	if err != nil {
		if err == sql.ErrNoRows {
			return GetTournamentTeamInterface{}, nil
		}
		return GetTournamentTeamInterface{}, err
	}
	return i, nil
}

const getTournamentTeams = `
SELECT
    tt.tournament_id, JSON_BUILD_OBJECT('id', tm.id, 'public_id', tm.public_id, 'user_id', tm.user_id, 'name', tm.name, 'slug', tm.slug, 'short_name', tm.shortname, 'admin', tm.admin, 'media_url', tm.media_url, 'gender', tm.gender, 'national', tm.national, 'country', tm.country, 'type', tm.type, 'player_count', tm.player_count, 'game_id', tm.game_id) AS team_data
FROM tournament_team tt
JOIN teams AS tm ON tm.id = tt.team_id
WHERE tt.public_id=$1
`

type GetTournamentTeamsRow struct {
	TournamentID int64           `json:"tournament_id"`
	TeamData     json.RawMessage `json:"team_data"`
}

func (q *Queries) GetTournamentTeams(ctx context.Context, publicID uuid.UUID) ([]GetTournamentTeamsRow, error) {
	rows, err := q.db.QueryContext(ctx, getTournamentTeams, publicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTournamentTeamsRow
	for rows.Next() {
		var i GetTournamentTeamsRow
		if err := rows.Scan(&i.TournamentID, &i.TeamData); err != nil {
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

const getTournamentTeamsCount = `
SELECT COUNT(*) FROM tournament_team tm
JOIN tournaments AS t ON t.id = tm.tournament_id
WHERE t.public_id=$1
`

func (q *Queries) GetTournamentTeamsCount(ctx context.Context, tournamentPublicID uuid.UUID) (int64, error) {
	row := q.db.QueryRowContext(ctx, getTournamentTeamsCount, tournamentPublicID)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

const newTournamentTeam = `
WITH tournamentID AS (
	SELECT * FROM tournaments WHERE public_id=$1
),
teamID AS (
	SELECT * FROM teams WHERE public_id=$2
)
INSERT INTO tournament_team (
    tournament_id,
    team_id
)
SELECT
	tournamentID.id,
	teamID.id
FROM tournamentID, teamID
RETURNING * 
`

func (q *Queries) NewTournamentTeam(ctx context.Context, tournamentPublicID, teamPublicID uuid.UUID) (models.TournamentTeam, error) {
	row := q.db.QueryRowContext(ctx, newTournamentTeam, tournamentPublicID, teamPublicID)
	var i models.TournamentTeam
	err := row.Scan(&i.TournamentID, &i.TeamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.TournamentTeam{}, nil
		}
		return models.TournamentTeam{}, err
	}
	return i, nil
}
