package database

import (
	"context"
	"khelogames/database/models"
)

const getMatch = `
SELECT * FROM matches
WHERE id=$1 AND tournament_id=$2
`

type GetMatchParams struct {
	ID           int64 `json:"id"`
	TournamentID int64 `json:"tournament_id"`
}

func (q *Queries) GetMatch(ctx context.Context, arg GetMatchParams) (models.Match, error) {
	row := q.db.QueryRowContext(ctx, getMatch, arg.ID, arg.TournamentID)
	var i models.Match
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.AwayTeamID,
		&i.HomeTeamID,
		&i.StartTimestamp,
		&i.EndTimestamp,
		&i.Type,
		&i.StatusCode,
		&i.Result,
	)
	return i, err
}

const getMatchByID = `
WITH cricket_groups AS (
    SELECT team_id, group_id
    FROM cricket_standing
    WHERE tournament_id = $1
),
football_groups AS (
    SELECT team_id, group_id
    FROM football_standing
    WHERE tournament_id = $1
)
SELECT DISTINCT
    m.id,
    m.tournament_id,
    m.away_team_id,
    m.home_team_id,
    m.start_timestamp,
    m.end_timestamp,
    m.type,
    m.status_code,
    m.result,
    m.stage,
    CASE  WHEN m.stage = 'Knockout' THEN m.knockout_level_id
		  ELSE NULL
	END AS knockout_level_id,
    CASE
        WHEN m.stage = 'Group' AND g.name = 'cricket' THEN cg.group_id
        WHEN m.stage = 'Group' AND g.name = 'football' THEN fg.group_id
        ELSE NULL
    END AS group_id,
    t1.name AS home_team_name,
    t1.slug AS home_team_slug,
    t1.shortName AS home_team_shortName,
    t1.media_url AS home_team_media_url,
    t1.gender AS home_team_gender,
    t1.country AS home_team_country,
    t1.national AS home_team_national,
    t1.type AS home_team_type,
    t1.player_count AS home_team_player_count,
    t1.game_id AS home_game_id,
    t2.name AS away_team_name,
    t2.slug AS away_team_slug,
    t2.shortName AS away_team_shortName,
    t2.media_url AS away_team_media_url,
    t2.gender AS away_team_gender,
    t2.country AS away_team_country,
    t2.national AS away_team_national,
    t2.type AS away_team_type,
    t2.player_count AS away_team_player_count,
    t1.game_id AS away_game_id
FROM matches m
JOIN teams AS t1 ON m.home_team_id = t1.id
JOIN teams AS t2 ON m.away_team_id = t2.id
JOIN games AS g ON g.id = t1.game_id
LEFT JOIN cricket_groups cg ON cg.team_id = m.home_team_id OR cg.team_id = m.away_team_id
LEFT JOIN football_groups fg ON fg.team_id = m.home_team_id OR fg.team_id = m.away_team_id
LEFT JOIN groups AS gr ON gr.id = 
    CASE
        WHEN m.stage = 'Group' AND g.name = 'cricket' THEN cg.group_id
        WHEN m.stage = 'Group' AND g.name = 'football' THEN fg.group_id
        ELSE NULL
    END
WHERE m.tournament_id = $1
ORDER BY m.id ASC;
`

type GetMatchByIDRow struct {
	ID                  int64   `json:"id"`
	TournamentID        int64   `json:"tournament_id"`
	AwayTeamID          int64   `json:"away_team_id"`
	HomeTeamID          int64   `json:"home_team_id"`
	StartTimestamp      int64   `json:"start_timestamp"`
	EndTimestamp        int64   `json:"end_timestamp"`
	Type                string  `json:"type"`
	StatusCode          string  `json:"status_code"`
	Result              *int64  `json:"result"`
	Stage               *string `json:"stage"`
	KnockoutLevelID     *int32  `json:"knockout_level_id"`
	GroupID             *int64  `json:"group_id"`
	HomeTeamName        string  `json:"home_team_name"`
	HomeTeamSlug        string  `json:"home_team_slug"`
	HomeTeamShortname   string  `json:"home_team_shortname"`
	HomeTeamMediaUrl    string  `json:"home_team_media_url"`
	HomeTeamGender      string  `json:"home_team_gender"`
	HomeTeamCountry     string  `json:"home_team_country"`
	HomeTeamNational    bool    `json:"home_team_national"`
	HomeTeamType        string  `json:"home_team_type"`
	HomeTeamPlayerCount int32   `json:"home_team_player_count"`
	HomeGameID          int64   `json:"home_game_id"`
	AwayTeamName        string  `json:"away_team_name"`
	AwayTeamSlug        string  `json:"away_team_slug"`
	AwayTeamShortname   string  `json:"away_team_shortname"`
	AwayTeamMediaUrl    string  `json:"away_team_media_url"`
	AwayTeamGender      string  `json:"away_team_gender"`
	AwayTeamCountry     string  `json:"away_team_country"`
	AwayTeamNational    bool    `json:"away_team_national"`
	AwayTeamType        string  `json:"away_team_type"`
	AwayTeamPlayerCount int32   `json:"away_team_player_count"`
	AwayGameID          int64   `json:"away_game_id"`
}

func (q *Queries) GetMatchByID(ctx context.Context, tournamentID int64) ([]GetMatchByIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getMatchByID, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMatchByIDRow
	for rows.Next() {
		var i GetMatchByIDRow
		if err := rows.Scan(
			&i.ID,
			&i.TournamentID,
			&i.AwayTeamID,
			&i.HomeTeamID,
			&i.StartTimestamp,
			&i.EndTimestamp,
			&i.Type,
			&i.StatusCode,
			&i.Result,
			&i.Stage,
			&i.KnockoutLevelID,
			&i.GroupID,
			&i.HomeTeamName,
			&i.HomeTeamSlug,
			&i.HomeTeamShortname,
			&i.HomeTeamMediaUrl,
			&i.HomeTeamGender,
			&i.HomeTeamCountry,
			&i.HomeTeamNational,
			&i.HomeTeamType,
			&i.HomeTeamPlayerCount,
			&i.HomeGameID,
			&i.AwayTeamName,
			&i.AwayTeamSlug,
			&i.AwayTeamShortname,
			&i.AwayTeamMediaUrl,
			&i.AwayTeamGender,
			&i.AwayTeamCountry,
			&i.AwayTeamNational,
			&i.AwayTeamType,
			&i.AwayTeamPlayerCount,
			&i.AwayGameID,
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

const getMatchByMatchID = `
SELECT
    m.id, m.tournament_id, m.away_team_id, m.home_team_id, m.start_timestamp, m.end_timestamp, m.type, m.status_code, m.result, m.stage, m.knockout_level_id,
    t1.name AS home_team_name, t1.slug AS home_team_slug, t1.shortName AS home_team_shortName, t1.media_url AS home_team_media_url, t1.gender AS home_team_gender, t1.country AS home_team_country, t1.national AS home_team_national, t1.type AS home_team_type, t1.player_count AS home_team_player_count, t1.game_id AS home_game_id,
    t2.name AS away_team_name, t2.slug AS away_team_slug, t2.shortName AS away_team_shortName, t2.media_url AS away_team_media_url, t2.gender AS away_team_gender, t2.country AS away_team_country, t2.national AS away_team_national, t2.type AS away_team_type, t2.player_count AS away_team_player_count, t1.game_id AS away_game_id
FROM matches m
JOIN teams AS t1 ON m.home_team_id=t1.id
JOIN teams AS t2 ON m.away_team_id=t2.id
WHERE m.id=$1
`

func (q *Queries) GetMatchByMatchID(ctx context.Context, id int64) (GetMatchByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getMatchByMatchID, id)
	var i GetMatchByIDRow
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.AwayTeamID,
		&i.HomeTeamID,
		&i.StartTimestamp,
		&i.EndTimestamp,
		&i.Type,
		&i.StatusCode,
		&i.Result,
		&i.Stage,
		&i.KnockoutLevelID,
		&i.HomeTeamName,
		&i.HomeTeamSlug,
		&i.HomeTeamShortname,
		&i.HomeTeamMediaUrl,
		&i.HomeTeamGender,
		&i.HomeTeamCountry,
		&i.HomeTeamNational,
		&i.HomeTeamType,
		&i.HomeTeamPlayerCount,
		&i.HomeGameID,
		&i.AwayTeamName,
		&i.AwayTeamSlug,
		&i.AwayTeamShortname,
		&i.AwayTeamMediaUrl,
		&i.AwayTeamGender,
		&i.AwayTeamCountry,
		&i.AwayTeamNational,
		&i.AwayTeamType,
		&i.AwayTeamPlayerCount,
		&i.AwayGameID,
	)
	return i, err
}

const getMatches = `-
SELECT id, tournament_id, away_team_id, home_team_id, start_timestamp, end_timestamp, type, status_code FROM matches
WHERE tournament_id=$1
ORDER BY id DESC
`

func (q *Queries) GetMatches(ctx context.Context, tournamentID int64) ([]models.Match, error) {
	rows, err := q.db.QueryContext(ctx, getMatches, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Match
	for rows.Next() {
		var i models.Match
		if err := rows.Scan(
			&i.ID,
			&i.TournamentID,
			&i.AwayTeamID,
			&i.HomeTeamID,
			&i.StartTimestamp,
			&i.EndTimestamp,
			&i.Type,
			&i.StatusCode,
			&i.Result,
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

const getMatchesByTournamentID = `
SELECT id, tournament_id, away_team_id, home_team_id, start_timestamp, end_timestamp, type, status_code, result, stage FROM matches
WHERE tournament_id=$1
ORDER BY id ASC
`

func (q *Queries) GetMatchesByTournamentID(ctx context.Context, tournamentID int64) ([]models.Match, error) {
	rows, err := q.db.QueryContext(ctx, getMatchesByTournamentID, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Match
	for rows.Next() {
		var i models.Match
		if err := rows.Scan(
			&i.ID,
			&i.TournamentID,
			&i.AwayTeamID,
			&i.HomeTeamID,
			&i.StartTimestamp,
			&i.EndTimestamp,
			&i.Type,
			&i.StatusCode,
			&i.Result,
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

const newMatch = `
INSERT INTO matches (
    tournament_id,
    away_team_id,
    home_team_id,
    start_timestamp,
    end_timestamp,
    type,
    status_code,
	result,
	stage,
	knockout_level_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, &10, &11
) RETURNING id, tournament_id, away_team_id, home_team_id, start_timestamp, end_timestamp, type, status_code, result, stage, knockout_level_id
`

type NewMatchParams struct {
	TournamentID    int64  `json:"tournament_id"`
	AwayTeamID      int64  `json:"away_team_id"`
	HomeTeamID      int64  `json:"home_team_id"`
	StartTimestamp  int64  `json:"start_timestamp"`
	EndTimestamp    int64  `json:"end_timestamp"`
	Type            string `json:"type"`
	StatusCode      string `json:"status_code"`
	Result          *int64 `json:"result"`
	Stage           string `json:"stage"`
	KnockoutLevelID *int32 `json:"knockout_level_id"`
}

func (q *Queries) NewMatch(ctx context.Context, arg NewMatchParams) (models.Match, error) {
	row := q.db.QueryRowContext(ctx, newMatch,
		arg.TournamentID,
		arg.AwayTeamID,
		arg.HomeTeamID,
		arg.StartTimestamp,
		arg.EndTimestamp,
		arg.Type,
		arg.StatusCode,
		arg.Result,
		arg.Stage,
		arg.KnockoutLevelID,
	)
	var i models.Match
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.AwayTeamID,
		&i.HomeTeamID,
		&i.StartTimestamp,
		&i.EndTimestamp,
		&i.Type,
		&i.StatusCode,
		&i.Result,
		&i.Stage,
		&i.KnockoutLevelID,
	)
	return i, err
}

const updateMatchSchedule = `
UPDATE matches
SET start_timestamp=$1
WHERE id=$2
RETURNING id, tournament_id, away_team_id, home_team_id, start_timestamp, end_timestamp, type, status_code, result, stage
`

type UpdateMatchScheduleParams struct {
	StartTimestamp int64 `json:"start_timestamp"`
	ID             int64 `json:"id"`
}

func (q *Queries) UpdateMatchSchedule(ctx context.Context, arg UpdateMatchScheduleParams) (models.Match, error) {
	row := q.db.QueryRowContext(ctx, updateMatchSchedule, arg.StartTimestamp, arg.ID)
	var i models.Match
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.AwayTeamID,
		&i.HomeTeamID,
		&i.StartTimestamp,
		&i.EndTimestamp,
		&i.Type,
		&i.StatusCode,
		&i.Result,
		&i.Stage,
	)
	return i, err
}

const updateMatchStatus = `
UPDATE matches
SET status_code=$1
WHERE id=$2
RETURNING id, tournament_id, away_team_id, home_team_id, start_timestamp, end_timestamp, type, status_code, result, stage
`

type UpdateMatchStatusParams struct {
	StatusCode string `json:"status_code"`
	ID         int64  `json:"id"`
}

func (q *Queries) UpdateMatchStatus(ctx context.Context, arg UpdateMatchStatusParams) (models.Match, error) {
	row := q.db.QueryRowContext(ctx, updateMatchStatus, arg.StatusCode, arg.ID)
	var i models.Match
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.AwayTeamID,
		&i.HomeTeamID,
		&i.StartTimestamp,
		&i.EndTimestamp,
		&i.Type,
		&i.StatusCode,
		&i.Result,
		&i.Stage,
	)
	return i, err
}

const updateMatchResult = `
UPDATE matches
SET result=$1
WHERE id=$2
RETURNING *
`

func (q *Queries) UpdateMatchResult(ctx context.Context, id, result int64) (models.Match, error) {
	row := q.db.QueryRowContext(ctx, updateMatchResult, result, id)
	var i models.Match
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.AwayTeamID,
		&i.HomeTeamID,
		&i.StartTimestamp,
		&i.EndTimestamp,
		&i.Type,
		&i.StatusCode,
		&i.Result,
		&i.Stage,
	)
	return i, err
}
