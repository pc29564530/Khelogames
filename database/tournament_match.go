package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const getMatchByTournamentPublicID = `
WITH cricket_groups AS (
    SELECT cs.team_id, cs.group_id
    FROM cricket_standing cs
    JOIN tournaments t ON t.id = cs.tournament_id
    WHERE t.public_id = $1
),
football_groups AS (
    SELECT fs.team_id, fs.group_id
    FROM football_standing fs
    JOIN tournaments t ON t.id = fs.tournament_id
    WHERE t.public_id = $1
)
SELECT DISTINCT
    m.id,
    m.public_id,
    m.tournament_id,
    m.away_team_id,
    m.home_team_id,
    m.start_timestamp,
    m.end_timestamp,
    m.type,
    m.status_code,
    m.result,
    m.stage,
    CASE  
        WHEN m.stage = 'Knockout' THEN m.knockout_level_id
        ELSE NULL
    END AS knockout_level_id,
	m.match_format,
	m.day_number,
	m.sub_status,
    CASE
        WHEN m.stage = 'Group' AND g.name = 'cricket' THEN cg.group_id
        WHEN m.stage = 'Group' AND g.name = 'football' THEN fg.group_id
        ELSE NULL
    END AS group_id,
    t1.id AS home_team_id,
    t1.public_id AS home_team_public_id,
    t1.user_id AS home_team_user_id,
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
    t2.id AS away_team_id,
    t2.public_id AS away_team_public_id,
    t2.user_id AS away_team_user_id,
    t2.name AS away_team_name,
    t2.slug AS away_team_slug,
    t2.shortName AS away_team_shortName,
    t2.media_url AS away_team_media_url,
    t2.gender AS away_team_gender,
    t2.country AS away_team_country,
    t2.national AS away_team_national,
    t2.type AS away_team_type,
    t2.player_count AS away_team_player_count,
    t2.game_id AS away_game_id
FROM matches m
LEFT JOIN teams t1 ON m.home_team_id = t1.id
LEFT JOIN teams t2 ON m.away_team_id = t2.id
JOIN games g ON g.id = t1.game_id
LEFT JOIN cricket_groups cg ON cg.team_id = m.home_team_id OR cg.team_id = m.away_team_id
LEFT JOIN football_groups fg ON fg.team_id = m.home_team_id OR fg.team_id = m.away_team_id
LEFT JOIN groups gr ON gr.id = 
    CASE
        WHEN m.stage = 'Group' AND g.name = 'cricket' THEN cg.group_id
        WHEN m.stage = 'Group' AND g.name = 'football' THEN fg.group_id
        ELSE NULL
    END
JOIN tournaments t ON t.id = m.tournament_id
WHERE t.public_id = $1
ORDER BY m.id ASC;
`

type GetMatchByIDRow struct {
	ID                  int64     `json:"id"`
	PublicID            uuid.UUID `json:"public_id"`
	TournamentID        int32     `json:"tournament_id"`
	AwayTeamID          int32     `json:"away_team_id"`
	HomeTeamID          int32     `json:"home_team_id"`
	StartTimestamp      int64     `json:"start_timestamp"`
	EndTimestamp        int64     `json:"end_timestamp"`
	Type                string    `json:"type"`
	StatusCode          string    `json:"status_code"`
	Result              *int64    `json:"result"`
	Stage               *string   `json:"stage"`
	KnockoutLevelID     *int32    `json:"knockout_level_id"`
	MatchFormat         string    `json:"match_format"`
	DayNumber           *int      `json:"day_number"`
	SubStatus           *string   `json:"sub_status"`
	LocationID          *int32    `json:"location_id"`
	LocationLocked      bool      `json:"location_locked"`
	GroupID             *int64    `json:"group_id"`
	HomeTeamPublicID    uuid.UUID `json:"home_team_public_id"`
	HomeTeamUserID      int32     `json:"home_team_user_id"`
	HomeTeamName        string    `json:"home_team_name"`
	HomeTeamSlug        string    `json:"home_team_slug"`
	HomeTeamShortname   string    `json:"home_team_shortname"`
	HomeTeamMediaUrl    string    `json:"home_team_media_url"`
	HomeTeamGender      string    `json:"home_team_gender"`
	HomeTeamCountry     string    `json:"home_team_country"`
	HomeTeamNational    bool      `json:"home_team_national"`
	HomeTeamType        string    `json:"home_team_type"`
	HomeTeamPlayerCount int32     `json:"home_team_player_count"`
	HomeGameID          int64     `json:"home_game_id"`
	HomeTeamLocationID  *int32    `json:"home_team_location_id"`
	AwayTeamPublicID    uuid.UUID `json:"away_team_public_id"`
	AwayTeamUserID      int32     `json:"away_team_user_id"`
	AwayTeamName        string    `json:"away_team_name"`
	AwayTeamSlug        string    `json:"away_team_slug"`
	AwayTeamShortname   string    `json:"away_team_shortname"`
	AwayTeamMediaUrl    string    `json:"away_team_media_url"`
	AwayTeamGender      string    `json:"away_team_gender"`
	AwayTeamCountry     string    `json:"away_team_country"`
	AwayTeamNational    bool      `json:"away_team_national"`
	AwayTeamType        string    `json:"away_team_type"`
	AwayTeamPlayerCount int32     `json:"away_team_player_count"`
	AwayGameID          int64     `json:"away_game_id"`
	AwayTeamLocationID  *int32    `json:"away_team_location_id"`
}

func (q *Queries) GetMatchByTournamentPublicID(ctx context.Context, tournamentPublicID uuid.UUID) ([]GetMatchByIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getMatchByTournamentPublicID, tournamentPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMatchByIDRow
	for rows.Next() {
		var i GetMatchByIDRow
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
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
			&i.MatchFormat,
			&i.DayNumber,
			&i.SubStatus,
			&i.GroupID,
			&i.HomeTeamID,
			&i.HomeTeamPublicID,
			&i.HomeTeamUserID,
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
			&i.AwayTeamID,
			&i.AwayTeamPublicID,
			&i.AwayTeamUserID,
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
    m.id, m.public_id, m.tournament_id, m.away_team_id, m.home_team_id, m.start_timestamp, m.end_timestamp, m.type, m.status_code, m.result, m.stage, m.knockout_level_id, COALESCE(m.match_format, '') AS match_format, m.day_number, m.sub_status, m.location_id, m.location_locked,
    t1.id AS id, t1.public_id, t1.user_id, t1.name AS home_team_name, t1.slug AS home_team_slug, t1.shortName AS home_team_shortName, t1.media_url AS home_team_media_url, t1.gender AS home_team_gender, t1.country AS home_team_country, t1.national AS home_team_national, t1.type AS home_team_type, t1.player_count AS home_team_player_count, t1.game_id AS home_game_id, t1.location_id AS home_team_location_id,
    t2.id AS id,t2.public_id, t2.user_id, t2.name AS away_team_name, t2.slug AS away_team_slug, t2.shortName AS away_team_shortName, t2.media_url AS away_team_media_url, t2.gender AS away_team_gender, t2.country AS away_team_country, t2.national AS away_team_national, t2.type AS away_team_type, t2.player_count AS away_team_player_count, t1.game_id AS away_game_id, t2.location_id AS away_team_location_id
FROM matches m
JOIN teams AS t1 ON m.home_team_id=t1.id
JOIN teams AS t2 ON m.away_team_id=t2.id
WHERE m.public_id=$1
`

func (q *Queries) GetTournamentMatchByMatchID(ctx context.Context, matchPublicID uuid.UUID) (GetMatchByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getMatchByMatchID, matchPublicID)
	var i GetMatchByIDRow
	err := row.Scan(
		&i.ID,
		&i.PublicID,
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
		&i.MatchFormat,
		&i.DayNumber,
		&i.SubStatus,
		&i.LocationID,
		&i.LocationLocked,

		&i.HomeTeamID,
		&i.HomeTeamPublicID,
		&i.HomeTeamUserID,
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
		&i.HomeTeamLocationID,

		&i.AwayTeamID,
		&i.AwayTeamPublicID,
		&i.AwayTeamUserID,
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
		&i.AwayTeamLocationID,
	)
	return i, err
}

const newMatch = `
WITH tournamentID AS (
	SELECT * FROM tournaments WHERE public_id = $1
),
awayTeamID AS (
	SELECT * FROM teams WHERE public_id = $2
),
homeTeamID AS (
	SELECT * FROM teams WHERE public_id = $3
)
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
	knockout_level_id,
	match_format,
	day_number,
	sub_status,
	location_id,
	location_locked,
	game_id
)
SELECT 
	tournamentID.id,
	awayTeamID.id,
	homeTeamID.id,
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
	$15,
	$16
FROM tournamentID, awayTeamID, homeTeamID
RETURNING *;`

type NewMatchParams struct {
	TournamentPublicID uuid.UUID `json:"tournament_id"`
	AwayTeamPublicID   uuid.UUID `json:"away_team_id"`
	HomeTeamPublicID   uuid.UUID `json:"home_team_id"`
	StartTimestamp     int64     `json:"start_timestamp"`
	EndTimestamp       int64     `json:"end_timestamp"`
	Type               string    `json:"type"`
	StatusCode         string    `json:"status_code"`
	Result             *int64    `json:"result"`
	Stage              string    `json:"stage"`
	KnockoutLevelID    *int32    `json:"knockout_level_id"`
	MatchFormat        *string   `json:"match_format"`
	DayNumber          *int      `json:"day_number"`
	SubStatus          *string   `json:"sub_status"`
	LocationID         int32     `json:"location_id"`
	LocationLocked     bool      `json:"location_locked"`
	GameID             int32     `json:"game_id"`
}

func (q *Queries) NewMatch(ctx context.Context, arg NewMatchParams) (models.Match, error) {
	row := q.db.QueryRowContext(ctx, newMatch,
		arg.TournamentPublicID,
		arg.AwayTeamPublicID,
		arg.HomeTeamPublicID,
		arg.StartTimestamp,
		arg.EndTimestamp,
		arg.Type,
		arg.StatusCode,
		arg.Result,
		arg.Stage,
		arg.KnockoutLevelID,
		arg.MatchFormat,
		arg.DayNumber,
		arg.SubStatus,
		arg.LocationID,
		arg.LocationLocked,
		arg.GameID,
	)
	var i models.Match
	err := row.Scan(
		&i.ID,
		&i.PublicID,
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
		&i.MatchFormat,
		&i.DayNumber,
		&i.SubStatus,
		&i.LocationID,
		&i.LocationLocked,
		&i.GameID,
	)
	return i, err
}

const updateMatchSchedule = `
UPDATE matches
SET start_timestamp=$2
WHERE public_id=$1
RETURNING *
`

type UpdateMatchScheduleParams struct {
	PublicID       uuid.UUID `json:"public_id"`
	StartTimestamp int64     `json:"start_timestamp"`
}

func (q *Queries) UpdateMatchSchedule(ctx context.Context, matchPublicID uuid.UUID, startTimestamp int64) (models.Match, error) {
	row := q.db.QueryRowContext(ctx, updateMatchSchedule, matchPublicID, startTimestamp)
	var i models.Match
	err := row.Scan(
		&i.ID,
		&i.PublicID,
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
		&i.MatchFormat,
		&i.DayNumber,
		&i.SubStatus,
		&i.LocationID,
		&i.LocationLocked,
	)
	return i, err
}

const updateMatchStatus = `
UPDATE matches
SET status_code=$2
WHERE public_id=$1
RETURNING *
`

func (q *Queries) UpdateMatchStatus(ctx context.Context, matchPublicID uuid.UUID, status string) (models.Match, error) {
	row := q.db.QueryRowContext(ctx, updateMatchStatus, matchPublicID, status)
	var i models.Match
	err := row.Scan(
		&i.ID,
		&i.PublicID,
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
		&i.MatchFormat,
		&i.DayNumber,
		&i.SubStatus,
		&i.LocationID,
		&i.LocationLocked,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Match{}, fmt.Errorf("Failed to get data", err)
		}
		return models.Match{}, fmt.Errorf("Failed to scan: ", err)
	}
	return i, err
}

const updateMatchResult = `
UPDATE matches
SET 
    result = $2,
    status_code = 'finished'
WHERE id = $1
RETURNING *
`

func (q *Queries) UpdateMatchResult(ctx context.Context, matchID, resultID int32) (models.Match, error) {
	row := q.db.QueryRowContext(ctx, updateMatchResult, matchID, resultID)
	var i models.Match
	err := row.Scan(
		&i.ID,
		&i.PublicID,
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
		&i.MatchFormat,
		&i.DayNumber,
		&i.SubStatus,
		&i.LocationID,
		&i.LocationLocked,
	)
	return i, err
}

const updateMatchSubStatus = `
	UPDATE matches m
	SET
		sub_status = CASE WHEN m.status_code = 'in_progress' THEN $2 ELSE NULL
		END 
	WHERE m.public_id = $1
	RETURNING *
`

func (q *Queries) UpdateMatchSubStatus(ctx context.Context, matchPublicID uuid.UUID, subStatus string) (models.Match, error) {
	row := q.db.QueryRowContext(ctx, updateMatchSubStatus, matchPublicID, subStatus)
	var i models.Match
	err := row.Scan(
		&i.ID,
		&i.PublicID,
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
		&i.MatchFormat,
		&i.DayNumber,
		&i.SubStatus,
		&i.LocationID,
		&i.LocationLocked,
	)
	return i, err
}
