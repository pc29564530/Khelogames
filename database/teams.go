package database

import (
	"context"
	"encoding/json"
	"khelogames/database/models"
	"log"
)

const addTeamPlayers = `
INSERT INTO team_players (
    team_id,
    player_id,
	join_date,
	leave_date
) VALUES ($1, $2, $3, $4)
RETURNING *
`

type AddTeamPlayersParams struct {
	TeamID    int64  `json:"team_id"`
	PlayerID  int64  `json:"player_id"`
	JoinDate  int32  `json:"join_date"`
	LeaveDate *int32 `json:"leave_date"`
}

func (q *Queries) AddTeamPlayers(ctx context.Context, arg AddTeamPlayersParams) (models.TeamPlayer, error) {
	row := q.db.QueryRowContext(ctx, addTeamPlayers, arg.TeamID, arg.PlayerID, arg.JoinDate, arg.LeaveDate)
	var i models.TeamPlayer
	err := row.Scan(&i.TeamID, &i.PlayerID, &i.JoinDate, &i.LeaveDate)
	return i, err
}

type GetMatchByTeamRow struct {
	TournamentID   int64  `json:"touranment_id"`
	TournamentName string `json:"tournament_name"`
	MatchID        int64  `json:"match_id"`
	HomeTeamID     int64  `json:"home_team_id"`
	AwayTeamID     int64  `json:"away_team_id"`
	HomeTeamName   string `json:"home_team_name"`
	AwayTeamName   string `json:"away_team_name"`
	StartTimestamp int64  `json:"start_timestamp"`
	StatusCode     string `json:"status_code"`
	Type           string `json:"type"`
}

const getMatchesByTeam = `
	SELECT 
		json_build_object(
			'id', m.id,
			'tournament_id', m.tournament_id,
			'away_team_id', m.away_team_id,
			'home_team_id', m.home_team_id,
			'start_timestamp', m.start_timestamp,
			'end_timestamp', m.end_timestamp,
			'type', m.type,
			'status_code', m.status_code,
			'result', m.result,
			'stage', m.stage,
			'knockout_level_id', m.knockout_level_id,
			'match_format', m.match_format,

			'homeTeam', json_build_object(
				'id', ht.id,
				'name', ht.name,
				'slug', ht.slug,
				'short_name', ht.shortname,
				'admin', ht.admin,
				'media_url', ht.media_url,
				'gender', ht.gender,
				'national', ht.national,
				'country', ht.country,
				'type', ht.type,
				'player_count', ht.player_count,
				'game_id', ht.game_id
			),

			'homeScore', CASE 
				WHEN g.name = 'football' THEN 
					json_build_object(
						'id', fs_home.id,
						'match_id', fs_home.match_id,
						'team_id', fs_home.team_id,
						'first_half', fs_home.first_half,
						'second_half', fs_home.second_half,
						'goals', fs_home.goals
					)
				WHEN g.name = 'cricket' THEN cricket_home_scores.scores
				ELSE NULL
			END,

			'awayTeam', json_build_object(
				'id', at.id,
				'name', at.name,
				'slug', at.slug,
				'short_name', at.shortname,
				'admin', at.admin,
				'media_url', at.media_url,
				'gender', at.gender,
				'national', at.national,
				'country', at.country,
				'type', at.type,
				'player_count', at.player_count,
				'game_id', at.game_id
			),

			'awayScore', CASE 
				WHEN g.name = 'football' THEN 
					json_build_object(
						'id', fs_away.id,
						'match_id', fs_away.match_id,
						'team_id', fs_away.team_id,
						'first_half', fs_away.first_half,
						'second_half', fs_away.second_half,
						'goals', fs_away.goals
					)
				WHEN g.name = 'cricket' THEN cricket_away_scores.scores
				ELSE NULL
			END,

			'tournament', json_build_object(
				'id', t.id,
				'name', t.name,
				'slug', t.slug,
				'country', t.country,
				'status_code', t.status_code,
				'level', t.level,
				'start_timestamp', t.start_timestamp,
				'game_id', t.game_id,
				'group_count', t.group_count,
				'max_group_team', t.max_group_team,
				'stage', t.stage,
				'has_knockout', t.has_knockout
			)
		) AS response

	FROM matches m

	JOIN teams ht ON m.home_team_id = ht.id
	JOIN teams at ON m.away_team_id = at.id
	LEFT JOIN tournaments t ON m.tournament_id = t.id
	JOIN games g ON t.game_id = g.id

	-- Football scores
	LEFT JOIN football_score fs_home ON fs_home.match_id = m.id AND fs_home.team_id = ht.id AND g.name = 'football'
	LEFT JOIN football_score fs_away ON fs_away.match_id = m.id AND fs_away.team_id = at.id AND g.name = 'football'

	-- Cricket scores for home team
	LEFT JOIN LATERAL (
		SELECT json_agg(
			json_build_object(
				'id', cs.id,
				'match_id', cs.match_id,
				'team_id', cs.team_id,
				'inning_number', cs.inning_number,
				'score', cs.score,
				'wickets', cs.wickets,
				'overs', cs.overs,
				'run_rate', cs.run_rate,
				'target_run_rate', cs.target_run_rate,
				'follow_on', cs.follow_on,
				'is_inning_completed', cs.is_inning_completed,
				'declared', cs.declared
			) ORDER BY cs.inning_number
		) AS scores
		FROM cricket_score cs
		WHERE cs.match_id = m.id AND cs.team_id = ht.id
	) AS cricket_home_scores ON true

	-- Cricket scores for away team
	LEFT JOIN LATERAL (
		SELECT json_agg(
			json_build_object(
				'id', cs.id,
				'match_id', cs.match_id,
				'team_id', cs.team_id,
				'inning_number', cs.inning_number,
				'score', cs.score,
				'wickets', cs.wickets,
				'overs', cs.overs,
				'run_rate', cs.run_rate,
				'target_run_rate', cs.target_run_rate,
				'follow_on', cs.follow_on,
				'is_inning_completed', cs.is_inning_completed,
				'declared', cs.declared
			) ORDER BY cs.inning_number
		) AS scores
		FROM cricket_score cs
		WHERE cs.match_id = m.id AND cs.team_id = at.id
	) AS cricket_away_scores ON true

	WHERE (ht.id = $1 OR at.id = $1) AND t.game_id = $2;
`

func (q *Queries) GetMatchesByTeam(ctx context.Context, teamID, gameID int64) ([]map[string]interface{}, error) {

	rows, err := q.db.QueryContext(ctx, getMatchesByTeam, teamID, gameID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var matches []map[string]interface{}

	for rows.Next() {
		var jsonBytes []byte
		if err := rows.Scan(&jsonBytes); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		var match map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &match); err != nil {
			log.Printf("Error unmarshaling row JSON: %v", err)
			return nil, err
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

const getMatchByTeam = `
SELECT t.id AS tournament_id, t.name, tm.id AS match_id, tm.home_team_id, tm.away_team_id, c1.name AS home_team_name, c2.name AS away_team_name, tm.start_timestamp, tm.status_code, tm.type
FROM matches tm
JOIN tournaments t ON tm.id = t.id
JOIN teams c1 ON tm.home_team_id = c1.id
JOIN teams c2 ON tm.away_team_id = c2.id
WHERE c1.id=$1 OR c2.id=$1
ORDER BY tm.id DESC, tm.start_timestamp DESC
`

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
			&i.StatusCode,
			&i.Type,
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

const getPlayerByTeam = `
SELECT 
  tp.team_id, 
  tp.player_id, 
  tp.join_date, 
  tp.leave_date, 
  p.id AS player_id,
  p.username, 
  p.slug, 
  p.short_name, 
  p.media_url, 
  p.positions, 
  p.country, 
  p.player_name, 
  p.game_id,
  p.profile_id
FROM team_players tp
JOIN players p ON tp.player_id = p.id
WHERE tp.team_id = $1 AND tp.leave_date IS NULL;
`

func (q *Queries) GetPlayerByTeam(ctx context.Context, teamID int64) ([]models.GetPlayerByTeam, error) {
	rows, err := q.db.QueryContext(ctx, getPlayerByTeam, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.GetPlayerByTeam
	for rows.Next() {
		var i models.GetPlayerByTeam
		if err := rows.Scan(
			&i.TeamID,
			&i.PlayerID,
			&i.JoinDate,
			&i.LeaveDate,
			&i.ID,
			&i.Username,
			&i.Slug,
			&i.ShortName,
			&i.MediaUrl,
			&i.Positions,
			&i.Country,
			&i.PlayerName,
			&i.GameID,
			&i.ProfileID,
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

const getTeam = `
SELECT id, name, slug, shortname, admin, media_url, gender, national, country, type, player_count, game_id FROM teams
WHERE id=$1
`

func (q *Queries) GetTeam(ctx context.Context, id int64) (models.Team, error) {
	row := q.db.QueryRowContext(ctx, getTeam, id)
	var i models.Team
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
		&i.PlayerCount,
		&i.GameID,
	)
	return i, err
}

const getTeamByPlayer = `
SELECT tm.* FROM team_players
JOIN teams AS tm ON team_players.team_id=tm.id
WHERE player_id=$1
`

func (q *Queries) GetTeamByPlayer(ctx context.Context, playerID int64) ([]models.GetTeamByPlayer, error) {
	rows, err := q.db.QueryContext(ctx, getTeamByPlayer, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.GetTeamByPlayer
	for rows.Next() {
		var i models.GetTeamByPlayer
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
			&i.PlayerCount,
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

const getTeamPlayers = `
SELECT * FROM team_players
WHERE team_id=$1
`

func (q *Queries) GetTeamPlayers(ctx context.Context, teamID int64) ([]models.TeamPlayer, error) {
	rows, err := q.db.QueryContext(ctx, getTeamPlayers, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.TeamPlayer
	for rows.Next() {
		var i models.TeamPlayer
		if err := rows.Scan(&i.TeamID, &i.PlayerID, &i.JoinDate, &i.LeaveDate); err != nil {
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

const getTeams = `
SELECT id, name, slug, shortname, admin, media_url, gender, national, country, type, player_count, game_id FROM teams
`

func (q *Queries) GetTeams(ctx context.Context) ([]models.Team, error) {
	rows, err := q.db.QueryContext(ctx, getTeams)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Team
	for rows.Next() {
		var i models.Team
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
			&i.PlayerCount,
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

const getTeamsBySport = `
SELECT 
    g.id, g.name, g.min_players,JSON_BUILD_OBJECT('id', t.id, 'name', t.name, 'slug', t.slug, 'short_name', t.shortname, 'admin', t.admin, 'media_url', t.media_url, 'gender', t.gender, 'national', t.national, 'country', t.country, 'type', t.type, 'player_count', t.player_count, 'game_id', t.game_id) AS team_data
FROM teams t
JOIN games AS g ON g.id = t.game_id
WHERE t.game_id=$1
`

type GetTeamsBySportRow struct {
	ID         int64           `json:"id"`
	Name       string          `json:"name"`
	MinPlayers int32           `json:"min_players"`
	TeamData   json.RawMessage `json:"team_data"`
}

func (q *Queries) GetTeamsBySport(ctx context.Context, gameID int64) ([]GetTeamsBySportRow, error) {
	rows, err := q.db.QueryContext(ctx, getTeamsBySport, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTeamsBySportRow
	for rows.Next() {
		var i GetTeamsBySportRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.MinPlayers,
			&i.TeamData,
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

const getTournamentsByTeam = `

SELECT * FROM tournaments t
JOIN tournament_team tt ON t.id=tt.id
JOIN teams c ON tt.team_id=c.id
WHERE c.id=$1
`

type GetTournamentsByTeamRow struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Country        string `json:"country"`
	StatusCode     string `json:"status_code"`
	Level          string `json:"level"`
	StartTimestamp int64  `json:"start_timestamp"`
	GameID         int32  `json:"game_id"`
	GroupCount     int    `json:"group_count"`
	MaxGroupTeam   int    `json:"max_group_team"`
	Stage          string `json:"stage"`
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
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Slug,
			&i.Country,
			&i.StatusCode,
			&i.Level,
			&i.StartTimestamp,
			&i.GameID,
			&i.GroupCount,
			&i.MaxGroupTeam,
			&i.Stage,
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

const newTeams = `
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
    player_count,
    game_id 
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING id, name, slug, shortname, admin, media_url, gender, national, country, type, player_count, game_id
`

type NewTeamsParams struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Shortname   string `json:"shortname"`
	Admin       string `json:"admin"`
	MediaUrl    string `json:"media_url"`
	Gender      string `json:"gender"`
	National    bool   `json:"national"`
	Country     string `json:"country"`
	Type        string `json:"type"`
	PlayerCount int32  `json:"player_count"`
	GameID      int32  `json:"game_id"`
}

func (q *Queries) NewTeams(ctx context.Context, arg NewTeamsParams) (models.Team, error) {
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
		arg.PlayerCount,
		arg.GameID,
	)
	var i models.Team
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
		&i.PlayerCount,
		&i.GameID,
	)
	return i, err
}

const searchTeam = `
SELECT id, name from teams
WHERE name LIKE $1
`

type SearchTeamRow struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) SearchTeam(ctx context.Context, name string) ([]models.SearchTeam, error) {
	rows, err := q.db.QueryContext(ctx, searchTeam, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.SearchTeam
	for rows.Next() {
		var i models.SearchTeam
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

const updateMediaUrl = `
UPDATE teams
SET media_url=$1
WHERE id=$2
RETURNING id, name, slug, shortname, admin, media_url, gender, national, country, type, player_count, game_id
`

type UpdateMediaUrlParams struct {
	MediaUrl string `json:"media_url"`
	ID       int64  `json:"id"`
}

func (q *Queries) UpdateMediaUrl(ctx context.Context, arg UpdateMediaUrlParams) (models.Team, error) {
	row := q.db.QueryRowContext(ctx, updateMediaUrl, arg.MediaUrl, arg.ID)
	var i models.Team
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
		&i.PlayerCount,
		&i.GameID,
	)
	return i, err
}

const updateTeamName = `
UPDATE teams
SET name=$1
WHERE id=$2
RETURNING id, name, slug, shortname, admin, media_url, gender, national, country, type, player_count, game_id
`

type UpdateTeamNameParams struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

func (q *Queries) UpdateTeamName(ctx context.Context, arg UpdateTeamNameParams) (models.Team, error) {
	row := q.db.QueryRowContext(ctx, updateTeamName, arg.Name, arg.ID)
	var i models.Team
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
		&i.PlayerCount,
		&i.GameID,
	)
	return i, err
}

const removePlayerFromTeam = `
UPDATE team_players
SET leave_date=$3
WHERE team_id=$1 AND player_id=$2
RETURNING *;
`

type UpdateLeaveDateParams struct {
	TeamID    int64  `json:"team_id"`
	PlayerID  int64  `json:"player_id"`
	LeaveDate *int32 `json:"leave_date"`
}

func (q *Queries) RemovePlayerFromTeam(ctx context.Context, arg UpdateLeaveDateParams) (models.TeamPlayer, error) {
	row := q.db.QueryRowContext(ctx, removePlayerFromTeam, arg.TeamID, arg.PlayerID, arg.LeaveDate)
	var i models.TeamPlayer
	err := row.Scan(
		&i.TeamID,
		&i.PlayerID,
		&i.JoinDate,
		&i.LeaveDate,
	)
	return i, err
}

const getTeamPlayer = `
SELECT COUNT(*) > 0
FROM team_players
WHERE team_id = $1 AND player_id = $2 AND leave_date IS NOT NULL
`

func (q *Queries) GetTeamPlayer(ctx context.Context, teamID int64, playerID int64) bool {
	var exists bool
	err := q.db.QueryRowContext(ctx, getTeamPlayer, teamID, playerID).Scan(&exists)
	if err != nil {
		ctx.Err()
		return false
	}
	return exists
}
