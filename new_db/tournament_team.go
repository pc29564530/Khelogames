package new_db

import (
	"context"
	"encoding/json"
	"khelogames/new_db/models"
)

const getTournamentTeam = `
SELECT tournament_id, team_id FROM tournament_team
WHERE team_id=$1
`

func (q *Queries) GetTournamentTeam(ctx context.Context, teamID int64) (models.TournamentTeam, error) {
	row := q.db.QueryRowContext(ctx, getTournamentTeam, teamID)
	var i models.TournamentTeam
	err := row.Scan(&i.TournamentID, &i.TeamID)
	return i, err
}

const getTournamentTeams = `
SELECT
    tt.tournament_id, JSON_BUILD_OBJECT('id', tm.id, 'name', tm.name, 'slug', tm.slug, 'short_name', tm.shortname, 'admin', tm.admin, 'media_url', tm.media_url, 'gender', tm.gender, 'national', tm.national, 'country', tm.country, 'type', tm.type, 'player_count', tm.player_count, 'game_id', tm.game_id) AS team_data
FROM tournament_team tt
JOIN teams AS tm ON tm.id = tt.team_id
WHERE tt.tournament_id=$1
`

type GetTournamentTeamsRow struct {
	TournamentID int64           `json:"tournament_id"`
	TeamData     json.RawMessage `json:"team_data"`
}

func (q *Queries) GetTournamentTeams(ctx context.Context, tournamentID int64) ([]GetTournamentTeamsRow, error) {
	rows, err := q.db.QueryContext(ctx, getTournamentTeams, tournamentID)
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
SELECT COUNT(*) FROM tournament_team
WHERE tournament_id=$1
`

func (q *Queries) GetTournamentTeamsCount(ctx context.Context, tournamentID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getTournamentTeamsCount, tournamentID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const newTournamentTeam = `
INSERT INTO tournament_team (
    tournament_id,
    team_id
) VALUES ( $1, $2 )
RETURNING tournament_id, team_id
`

type NewTournamentTeamParams struct {
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (q *Queries) NewTournamentTeam(ctx context.Context, arg NewTournamentTeamParams) (models.TournamentTeam, error) {
	row := q.db.QueryRowContext(ctx, newTournamentTeam, arg.TournamentID, arg.TeamID)
	var i models.TournamentTeam
	err := row.Scan(&i.TournamentID, &i.TeamID)
	return i, err
}
