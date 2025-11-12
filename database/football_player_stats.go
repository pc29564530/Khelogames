package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"khelogames/database/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const addFootballPlayerStats = `
	WITH resolved_ids AS (
		SELECT 
			p.id AS player_id
		FROM players p
		WHERE p.public_id=$1
	)
	INSERT INTO football_player_stats ( player_id, player_position, matches, minutes_played, goals_scored, goals_conceded, clean_sheet, assists, yellow_card, red_card, created_at, updated_at)
	SELECT
		ri.playe_id, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP,
	FROM resolved_ids ri
	RETURNING *
`

func (q *Queries) AddFootballPlayerStats(ctx *gin.Context, playerPublicID uuid.UUID) (*models.FootballPlayerStats, error) {
	rows := q.db.QueryRowContext(ctx, addFootballPlayerStats, playerPublicID)
	var stats models.FootballPlayerStats
	err := rows.Scan(
		&stats.ID,
		&stats.PublicID,
		&stats.PlayerID,
		&stats.PlayerPosition,
		&stats.Matches,
		&stats.MinutesPlayed,
		&stats.GoalsScored,
		&stats.GoalsConceded,
		&stats.CleanSheet,
		&stats.Assists,
		&stats.YellowCards,
		&stats.RedCards,
		&stats.Average,
		&stats.CreatedAt,
		&stats.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

const getFootballPlayerStatsQuery = `
	SELECT JSON_BUILD_OBJECT(
		'id', fps.id, 'pubic_id', fps.public_id, 'player_id', fps.player_id, 'player_position', fps.player_position, 'matches', fps.matches, 'minutes_played', fps.minutes_played, 'goals_scored, fps.GoalsScored,
		'goals_conceded', fps.goals_conceded, 'clean_sheet', fps.clean_sheet, 'assists', fps.assists, 'yellow_card', fps.yellow_card, 'red_card', fps.red_card, 'avergae', fps.average,
		'created_at', fps.created_at, 'updated_at', fps.updated_at,
		'player', JSON_BUILD_OBJECT(
				'id',p.id,
				'public_id', p.public_id,
				'user_id',p.user_Id,
				'game_id' p.game_id.
				'name', p.player_name, 
				'slug', p.slug, 
				'short_name', p.short_name, 
				'country', p.country, 
				'positions', p.positions, 
				'media_url', p.media_url,
				'created_at', p.created_at,
				'updated_at', p.updated_at,
			)
	) FROM football_player_stats fps
	JOIN players p ON p.id = fps.player_id
	WHERE p.public_id = $1
`

func (q *Queries) GetFootballPlayerStats(ctx context.Context, playerPublicID uuid.UUID) (*map[string]interface{}, error) {
	row := q.db.QueryRowContext(ctx, getFootballPlayerStatsQuery, playerPublicID)
	var jsonByte []byte
	var stats map[string]interface{}
	err := row.Scan(&jsonByte)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan row: %w", err)
	}

	err = json.Unmarshal(jsonByte, &stats)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal JSON: %w", err)
	}

	return &stats, nil
}

const addAndUpdateFootballPlayerStats = `
WITH match_context AS (
	SELECT id AS match_id FROM matches WHERE public_id = $1
),
subs AS (
	SELECT fsp.*, fi.incident_time
	FROM football_substitutions_player fsp
	JOIN football_incidents fi ON fi.id = fsp.incident_id
	JOIN match_context mc ON fi.match_id = mc.match_id
),
all_players AS (
	SELECT fs.player_id, fs.is_substitute
	FROM football_squad fs 
	JOIN match_context mc ON fs.match_id = mc.match_id
),
incident_data AS (
	SELECT
		fip.player_id,
		fi.incident_type,
		COUNT(*) AS count
	FROM football_incident_player fip
	JOIN football_incidents fi ON fip.incident_id = fi.id
	JOIN match_context mc ON fi.match_id = mc.match_id
	GROUP BY fip.player_id, fi.incident_type
),
minutes_played AS (
	SELECT
		ap.player_id,
		CASE
			WHEN ap.is_substitute = FALSE AND s.player_out_id = ap.player_id THEN s.incident_time
			WHEN ap.is_substitute = TRUE AND s.player_in_id = ap.player_id THEN 90 - s.incident_time
			WHEN ap.is_substitute = FALSE THEN 90
			ELSE 0
		END AS minutes
	FROM all_players ap
	LEFT JOIN subs s ON ap.player_id = s.player_out_id OR ap.player_id = s.player_in_id
),
aggregated AS (
	SELECT
		mp.player_id AS player_id,
		COALESCE(SUM(mp.minutes), 0) AS minutes_played,
		COALESCE(SUM(CASE WHEN id.incident_type = 'goals' THEN id.count ELSE 0 END), 0) AS goals_scored,
		COALESCE(SUM(CASE WHEN id.incident_type = 'assists' THEN id.count ELSE 0 END), 0) AS assists,
		COALESCE(SUM(CASE WHEN id.incident_type = 'goals_conceded' THEN id.count ELSE 0 END), 0) AS goals_conceded,
		COALESCE(SUM(CASE WHEN id.incident_type = 'clean_sheet' THEN id.count ELSE 0 END), 0) AS clean_sheets,
		COALESCE(SUM(CASE WHEN id.incident_type = 'yellow_card' THEN id.count ELSE 0 END), 0) AS yellow_cards,
		COALESCE(SUM(CASE WHEN id.incident_type = 'red_card' THEN id.count ELSE 0 END), 0) AS red_cards
	FROM minutes_played mp
	LEFT JOIN incident_data id ON mp.player_id = id.player_id
	GROUP BY mp.player_id
),
updated AS (
	UPDATE football_player_stats fps
	SET
		matches = fps.matches + 1,
		minutes_played = fps.minutes_played + a.minutes_played,
		goals_scored = fps.goals_scored + a.goals_scored,
		goals_conceded = fps.goals_conceded + a.goals_conceded,
		clean_sheet = fps.clean_sheet + a.clean_sheets,
		assists = fps.assists + a.assists,
		yellow_cards = fps.yellow_cards + a.yellow_cards,
		red_cards = fps.red_cards + a.red_cards,
		updated_at = NOW()
	FROM aggregated a
	WHERE fps.player_id = a.player_id
	RETURNING fps.player_id
)
INSERT INTO football_player_stats (
	player_id, player_position, matches, minutes_played, goals_scored, goals_conceded, clean_sheet, assists, yellow_cards, red_cards, average, created_at, updated_at )
SELECT 
	a.player_id,
	'',
	1,
	a.minutes_played,
	a.goals_scored,
	a.goals_conceded,
	a.clean_sheets,
	a.assists,
	a.yellow_cards,
	a.red_cards,
	'0.0',
	CURRENT_TIMESTAMP,
	CURRENT_TIMESTAMP
FROM aggregated a
WHERE NOT EXISTS (
	SELECT 1 FROM updated u WHERE u.player_id = a.player_id
)
RETURNING *;
`

func (q *Queries) AddORUpdateFootballPlayerStats(ctx context.Context, mathchPublicID uuid.UUID) (*[]models.FootballPlayerStats, error) {
	var playerStats []models.FootballPlayerStats
	rows, err := q.db.QueryContext(ctx, addAndUpdateFootballPlayerStats, mathchPublicID)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var stats models.FootballPlayerStats
		err = rows.Scan(
			&stats.ID,
			&stats.PublicID,
			&stats.PlayerID,
			&stats.PlayerPosition,
			&stats.Matches,
			&stats.MinutesPlayed,
			&stats.GoalsScored,
			&stats.GoalsConceded,
			&stats.CleanSheet,
			&stats.Assists,
			&stats.YellowCards,
			&stats.RedCards,
			&stats.Average,
			&stats.CreatedAt,
			&stats.UpdatedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("Failed to scan row: %w", err)
		}

		playerStats = append(playerStats, stats)
	}

	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, nil
	// 	}
	// 	return nil, fmt.Errorf("Failed to scan the query data: ", err)
	// }

	return &playerStats, nil
}
