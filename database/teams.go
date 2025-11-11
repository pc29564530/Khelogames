package database

import (
	"context"
	"encoding/json"
	"fmt"
	"khelogames/database/models"
	"log"

	"github.com/google/uuid"
)

const addTeamPlayers = `
WITH resolve_ids AS (
	SELECT
		t.id AS team_id,
		p.id AS player_id
	FROM teams t, players p
	WHERE t.public_id = $1 AND p.public_id = $2
)
INSERT INTO team_players (
    team_id,
    player_id,
	join_date,
	leave_date
)
SELECT 
	ri.team_id,
	ri.player_id,
	$3,
	$4
FROM resolve_ids ri
RETURNING *
`

type AddTeamPlayersParams struct {
	TeamPublicID   uuid.UUID `json:"team_public_id"`
	PlayerPublicID uuid.UUID `json:"player_public_id"`
	JoinDate       int32     `json:"join_date"`
	LeaveDate      *int32    `json:"leave_date"`
}

func (q *Queries) AddTeamPlayers(ctx context.Context, arg AddTeamPlayersParams) (models.TeamPlayer, error) {
	row := q.db.QueryRowContext(ctx, addTeamPlayers, arg.TeamPublicID, arg.PlayerPublicID, arg.JoinDate, arg.LeaveDate)
	var i models.TeamPlayer
	err := row.Scan(&i.TeamID, &i.PlayerID, &i.JoinDate, &i.LeaveDate)
	return i, err
}

type GetMatchByTeamRow struct {
	TournamentID   int32  `json:"touranment_id"`
	TournamentName string `json:"tournament_name"`
	MatchID        int64  `json:"match_id"`
	HomeTeamID     int32  `json:"home_team_id"`
	AwayTeamID     int32  `json:"away_team_id"`
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
			'public_id', m.public_id,
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
				'public_id', ht.public_id,
				'name', ht.name,
				'slug', ht.slug,
				'short_name', ht.shortname,
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
						'public_id', fs_home.public_id,
						'match_id', fs_home.match_id,
						'team_id', fs_home.team_id,
						'first_half', fs_home.first_half,
						'second_half', fs_home.second_half,
						'goals', fs_home.goals,
						'penalty_shootout', fs_home.penalty_shootout
					)
				WHEN g.name = 'cricket' THEN cricket_home_scores.scores
				ELSE NULL
			END,

			'awayTeam', json_build_object(
				'id', at.id,
				'public_id', at.public_id,
				'name', at.name,
				'slug', at.slug,
				'short_name', at.shortname,
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
						'public_id', fs_away.public_id,
						'match_id', fs_away.match_id,
						'team_id', fs_away.team_id,
						'first_half', fs_away.first_half,
						'second_half', fs_away.second_half,
						'goals', fs_away.goals,
						'penalty_shootout', fs_away.penalty_shootout
					)
				WHEN g.name = 'cricket' THEN cricket_away_scores.scores
				ELSE NULL
			END,

			'tournament', json_build_object(
				'id', t.id,
				'public_id', t.public_id,
				'name', t.name,
				'slug', t.slug,
				'country', t.country,
				'status', t.status,
				'level', t.level,
				'start_timestamp', t.start_timestamp,
				'game_id', t.game_id,
				'group_count', t.group_count,
				'max_group_team', t.max_group_teams,
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
				'public_id', cs.public_id,
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
				'public_id', cs.public_id,
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

	WHERE (ht.public_id = $1 OR at.public_id = $1) AND t.game_id = $2;
`

func (q *Queries) GetMatchesByTeam(ctx context.Context, teamPublicID uuid.UUID, gameID int64) ([]map[string]interface{}, error) {

	rows, err := q.db.QueryContext(ctx, getMatchesByTeam, teamPublicID, gameID)
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
WHERE c1.public_id=$1 OR c2.public_id=$1
ORDER BY tm.id DESC, tm.start_timestamp DESC
`

func (q *Queries) GetMatchByTeam(ctx context.Context, teamPublicID uuid.UUID) ([]GetMatchByTeamRow, error) {
	rows, err := q.db.QueryContext(ctx, getMatchByTeam, teamPublicID)
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
SELECT JSON_BUILD_OBJECT(
		'public_id', p.public_id,
		'user_id', p.user_id,
		'name', p.name,
		'slug', p.slug,
		'short_name', p.short_name,
		'media_url', p.media_url,
		'positions', p.positions,
		'country', p.country,
		'game', p.game_id,
		'join_date', tp.join_date
	)
FROM team_players tp
JOIN players p ON tp.player_id = p.id
JOIN teams t ON t.id = tp.team_id
WHERE t.public_id = $1 AND tp.leave_date IS NULL;
`

func (q *Queries) GetPlayerByTeam(ctx context.Context, teamPublicID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getPlayerByTeam, teamPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []map[string]interface{}
	for rows.Next() {
		var i map[string]interface{}
		var jsonByte []byte
		if err := rows.Scan(&jsonByte); err != nil {
			return nil, err
		}
		err := json.Unmarshal(jsonByte, &i)
		if err != nil {
			return nil, err
		}
		fmt.Println("Player: ", i)
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

// team by public_id
const getTeamByPublicID = `
	SELECT * FROM teams
	WHERE public_id=$1
`

func (q *Queries) GetTeamByPublicID(ctx context.Context, publicID uuid.UUID) (models.Team, error) {
	row := q.db.QueryRowContext(ctx, getTeamByPublicID, publicID)
	var i models.Team
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.Name,
		&i.Slug,
		&i.Shortname,
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

// team by public_id
const getTeamByID = `
	SELECT * FROM teams
	WHERE id=$1
`

func (q *Queries) GetTeamByID(ctx context.Context, id int64) (models.Team, error) {
	row := q.db.QueryRowContext(ctx, getTeamByID, id)
	var i models.Team
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.Name,
		&i.Slug,
		&i.Shortname,
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
SELECT JSON_BUILD_OBJECT(
	'id', t.id, 'public_id', t.public_id, 'name', t.name, 'slug', t.slug, 'shortname',t.shortname,
	'media_url', t.media_url, 'national', t.national, 'country', t.country, 'join_date', tp.join_date, 'leave_date', tp.leave_date 
) FROM team_players tp
JOIN teams AS t ON tp.team_id=t.id
JOIN players AS p ON tp.player_id = p.id
WHERE p.public_id=$1;
`

func (q *Queries) GetTeamByPlayer(ctx context.Context, playerPublicID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getTeamByPlayer, playerPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []map[string]interface{}
	for rows.Next() {
		var i map[string]interface{}
		var jsonByte []byte
		if err := rows.Scan(&jsonByte); err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonByte, &i)
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
SELECT tp.* FROM team_players tp
JOIN teams AS t ON t.id = tp.team_id
JOIN players AS p ON p.id = tp.user_id
WHERE t.public_id=$1
`

func (q *Queries) GetTeamPlayers(ctx context.Context, teamPublicID uuid.UUID) ([]models.TeamPlayer, error) {
	rows, err := q.db.QueryContext(ctx, getTeamPlayers, teamPublicID)
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
SELECT * FROM teams
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
			&i.PublicID,
			&i.UserID,
			&i.Name,
			&i.Slug,
			&i.Shortname,
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
	SELECT  JSON_BUILD_OBJECT('id', t.id, 'public_id', t.public_id, 'user_id', t.user_id, 'name', t.name, 'slug', t.slug, 'short_name', t.shortname, 'media_url', t.media_url, 'gender', t.gender, 'national', t.national, 'country', t.country, 'type', t.type, 'player_count', t.player_count, 'game_id', t.game_id)
	FROM teams t
	JOIN games AS g ON g.id = t.game_id
	WHERE t.game_id=$1
`

func (q *Queries) GetTeamsBySport(ctx context.Context, gameID int64) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getTeamsBySport, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []map[string]interface{}
	for rows.Next() {
		var i map[string]interface{}
		var jsonByte []byte
		if err := rows.Scan(&jsonByte); err != nil {
			return nil, err
		}
		err := json.Unmarshal(jsonByte, &i)
		if err != nil {
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
WHERE c.public_id=$1
`

type GetTournamentsByTeamRow struct {
	ID             int64     `json:"id"`
	PublicID       uuid.UUID `json:"public_id"`
	UserID         int32     `json:"user_id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	Country        string    `json:"country"`
	StatusCode     string    `json:"status_code"`
	Level          string    `json:"level"`
	StartTimestamp int64     `json:"start_timestamp"`
	GameID         int32     `json:"game_id"`
	GroupCount     int       `json:"group_count"`
	MaxGroupTeam   int       `json:"max_group_team"`
	Stage          string    `json:"stage"`
}

func (q *Queries) GetTournamentsByTeam(ctx context.Context, publicID uuid.UUID) ([]GetTournamentsByTeamRow, error) {
	rows, err := q.db.QueryContext(ctx, getTournamentsByTeam, publicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTournamentsByTeamRow
	for rows.Next() {
		var i GetTournamentsByTeamRow
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.UserID,
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
WITH userID AS (
	SELECT * FROM users
	WHERE public_id=$1
)
INSERT INTO teams (
	user_id,
    name,
    slug,
    shortName,
    media_url,
    gender,
    national,
    country,
    type,
    player_count,
    game_id 
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
	$11
FROM userID
RETURNING *;
`

type NewTeamsParams struct {
	UserPublicID uuid.UUID `json:"user_public_id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Shortname    string    `json:"shortname"`
	MediaUrl     string    `json:"media_url"`
	Gender       string    `json:"gender"`
	National     bool      `json:"national"`
	Country      string    `json:"country"`
	Type         string    `json:"type"`
	PlayerCount  int32     `json:"player_count"`
	GameID       int32     `json:"game_id"`
}

func (q *Queries) NewTeams(ctx context.Context, arg NewTeamsParams) (models.Team, error) {
	row := q.db.QueryRowContext(ctx, newTeams,
		arg.UserPublicID,
		arg.Name,
		arg.Slug,
		arg.Shortname,
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
		&i.PublicID,
		&i.UserID,
		&i.Name,
		&i.Slug,
		&i.Shortname,
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
SET media_url=$2
WHERE public_id=$1
RETURNING *
`

type UpdateMediaUrlParams struct {
	PublicID uuid.UUID `json:"public_id"`
	MediaUrl string    `json:"media_url"`
}

func (q *Queries) UpdateMediaUrl(ctx context.Context, publicID uuid.UUID, mediaUrl string) (models.Team, error) {
	row := q.db.QueryRowContext(ctx, updateMediaUrl, publicID, mediaUrl)
	var i models.Team
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.Name,
		&i.Slug,
		&i.Shortname,
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
SET name=$2
WHERE public_id=$1
RETURNING *
`

type UpdateTeamNameParams struct {
	PublicID uuid.UUID `json:"public_id"`
	Name     string    `json:"name"`
}

func (q *Queries) UpdateTeamName(ctx context.Context, publicID uuid.UUID, name string) (models.Team, error) {
	row := q.db.QueryRowContext(ctx, updateTeamName, publicID, name)
	var i models.Team
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.UserID,
		&i.Name,
		&i.Slug,
		&i.Shortname,
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
UPDATE team_players AS tp
SET leave_date = CASE WHEN $3 <> 0 THEN $3 ELSE NULL END
FROM teams t, players p
WHERE t.public_id = $1
  AND p.public_id = $2
  AND t.id = tp.team_id
  AND p.id = tp.player_id
RETURNING tp.*;
`

func (q *Queries) RemovePlayerFromTeam(ctx context.Context, teamPublicID, playerPublicID uuid.UUID, leaveDate int32) (models.TeamPlayer, error) {
	row := q.db.QueryRowContext(ctx, removePlayerFromTeam, teamPublicID, playerPublicID, leaveDate)
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
JOIN teams t ON t.id = team_players.team_id
JOIN players p ON p.id = team_players.player_id
WHERE t.public_id = $1 AND p.public_id = $2 AND leave_date IS NOT NULL
`

func (q *Queries) GetTeamPlayer(ctx context.Context, teamPublicID, playerPublicID uuid.UUID) bool {
	var exists bool
	err := q.db.QueryRowContext(ctx, getTeamPlayer, teamPublicID, playerPublicID).Scan(&exists)
	if err != nil {
		ctx.Err()
		return false
	}
	return exists
}
