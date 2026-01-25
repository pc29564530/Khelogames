package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"khelogames/database/models"
	"log"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const listMatchesByLocationQuery = `
    WITH nearby_locations AS (
        SELECT
            id,
            public_id,
            city,
            state,
            country,
            latitude,
            longitude,
            h3_index,
            (6371 * acos(
                LEAST(1.0, GREATEST(-1.0,
                    cos(radians($2::double precision)) *
                    cos(radians(latitude)) *
                    cos(radians(longitude) - radians($3::double precision)) +
                    sin(radians($2::double precision)) *
                    sin(radians(latitude))
                ))
            )) AS distance_km
        FROM locations
        WHERE h3_index = ANY($4)
    )
    SELECT json_build_object(
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
            'day_number', m.day_number,
            'sub_status', m.sub_status,
            'location_id', m.location_id,
            'location_locked', m.location_locked,
            'homeTeam', json_build_object(
                'id', ht.id,
                'public_id', ht.public_id,
                'user_id', ht.user_id,
                'name', ht.name,
                'slug', ht.slug,
                'short_name', ht.shortname,
                'media_url', ht.media_url,
                'gender', ht.gender,
                'national', ht.national,
                'country', ht.country,
                'type', ht.type,
                'player_count', ht.player_count,
                'game_id', ht.game_id,
                'location_id', ht.location_id
            ),

            'homeScore', CASE 
                WHEN g.name = 'football' THEN 
                    json_build_object(
                        'id', fs_home.id,
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
                'user_id', at.user_id,
                'name', at.name,
                'slug', at.slug,
                'short_name', at.shortname,
                'media_url', at.media_url,
                'gender', at.gender,
                'national', at.national,
                'country', at.country,
                'type', at.type,
                'player_count', at.player_count,
                'game_id', at.game_id,
                'location_id', at.location_id
            ),

            'awayScore', CASE 
                WHEN g.name = 'football' THEN 
                    json_build_object(
                        'id', fs_away.id,
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
                'user_id', t.user_id,
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
            ),

            'location', CASE
                WHEN loc.id IS NOT NULL THEN
                    json_build_object(
                        'id', loc.id,
                        'public_id', loc.public_id,
                        'city', loc.city,
                        'state', loc.state,
                        'country', loc.country,
                        'latitude', loc.latitude,
                        'longitude', loc.longitude,
                        'h3_index', loc.h3_index
                    )
                ELSE NULL
            END
        ) AS response
    FROM matches m
    INNER JOIN nearby_locations nl ON m.location_id = nl.id
    JOIN teams ht ON m.home_team_id = ht.id
    JOIN teams at ON m.away_team_id = at.id
    LEFT JOIN tournaments t ON m.tournament_id = t.id
    JOIN games g ON t.game_id = g.id
    LEFT JOIN locations loc ON m.location_id = loc.id

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
        ) as scores
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
        ) as scores
        FROM cricket_score cs
        WHERE cs.match_id = m.id AND cs.team_id = at.id
    ) AS cricket_away_scores ON true
    WHERE
        m.start_timestamp >= $1
        AND nl.distance_km <= $5
        AND m.game_id = $6
    ORDER BY nl.distance_km ASC, m.start_timestamp ASC
    LIMIT 50;
`

func (q *Queries) ListMatchesByLocation(ctx context.Context, startTimestamp int32, latitude, longitude float64, h3Indices []string, maxDistanceKm float64, gameID int32) ([]map[string]interface{}, error) {
	log.Printf("ListMatchesByLocation - Params: startTimestamp=%d, lat=%f, lng=%f, h3Count=%d, maxDist=%f, gameID=%d",
		startTimestamp, latitude, longitude, len(h3Indices), maxDistanceKm, gameID)

	rows, err := q.db.QueryContext(ctx, listMatchesByLocationQuery, startTimestamp, latitude, longitude, pq.StringArray(h3Indices), maxDistanceKm, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute ListMatchesByLocation query: %w", err)
	}
	defer rows.Close()

	var matches []map[string]interface{}
	for rows.Next() {
		var jsonByte []byte
		err := rows.Scan(&jsonByte)
		if err != nil {
			log.Printf("Failed to scan match with location: %v", err)
			continue
		}

		var m map[string]interface{}
		err = json.Unmarshal(jsonByte, &m)
		if err != nil {
			log.Printf("Failed to unmarshal match JSON: %v", err)
			continue
		}

		matches = append(matches, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error in ListMatchesByLocation: %w", err)
	}

	return matches, nil
}

const listMatchesQuery = `
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
            'day_number', m.day_number,
            'sub_status', m.sub_status,
            'location_id', m.location_id,
            'location_locked', m.location_locked,
            'homeTeam', json_build_object(
                'id', ht.id,
                'public_id', ht.public_id,
                'user_id', ht.user_id,
                'name', ht.name,
                'slug', ht.slug,
                'short_name', ht.shortname,
                'media_url', ht.media_url,
                'gender', ht.gender,
                'national', ht.national,
                'country', ht.country,
                'type', ht.type,
                'player_count', ht.player_count,
                'game_id', ht.game_id,
                'location_id', ht.location_id
            ),

            'homeScore', CASE 
                WHEN g.name = 'football' THEN 
                    json_build_object(
                        'id', fs_home.id,
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
                'user_id', at.user_id,
                'name', at.name,
                'slug', at.slug,
                'short_name', at.shortname,
                'media_url', at.media_url,
                'gender', at.gender,
                'national', at.national,
                'country', at.country,
                'type', at.type,
                'player_count', at.player_count,
                'game_id', at.game_id,
                'location_id', at.location_id
            ),

            'awayScore', CASE 
                WHEN g.name = 'football' THEN 
                    json_build_object(
                        'id', fs_away.id,
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
                'user_id', t.user_id,
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
            ),

            'location', CASE
                WHEN loc.id IS NOT NULL THEN
                    json_build_object(
                        'id', loc.id,
                        'public_id', loc.public_id,
                        'city', loc.city,
                        'state', loc.state,
                        'country', loc.country,
                        'latitude', loc.latitude,
                        'longitude', loc.longitude,
                        'h3_index', loc.h3_index
                    )
                ELSE NULL
            END
        ) AS response

    FROM matches m
    JOIN teams ht ON m.home_team_id = ht.id
    JOIN teams at ON m.away_team_id = at.id
    LEFT JOIN tournaments t ON m.tournament_id = t.id
    JOIN games g ON t.game_id = g.id
    LEFT JOIN locations loc ON m.location_id = loc.id

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
        ) as scores
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
        ) as scores
        FROM cricket_score cs
        WHERE cs.match_id = m.id AND cs.team_id = at.id
    ) AS cricket_away_scores ON true
    WHERE m.start_timestamp >= $1
    AND t.game_id = $2
    ORDER BY m.start_timestamp;
`

const getMatchByPublicIdQuery = `
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
        'day_number', m.day_number,
        'sub_status', m.sub_status,
        'location_id', m.location_id,
        'location_locked', m.location_locked,

        'homeTeam', json_build_object(
            'id', ht.id,
            'public_id', ht.public_id,
            'user_id', ht.user_id,
            'name', ht.name,
            'slug', ht.slug,
            'short_name', ht.shortname,
            'media_url', ht.media_url,
            'gender', ht.gender,
            'national', ht.national,
            'country', ht.country,
            'type', ht.type,
            'player_count', ht.player_count,
            'game_id', ht.game_id,
            'location_id', ht.location_id
        ),

        'homeScore', CASE 
            WHEN g.name = 'football' THEN 
                json_build_object(
                    'id', fs_home.id,
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
            'user_id', at.user_id,
            'name', at.name,
            'slug', at.slug,
            'media_url', at.media_url,
            'gender', at.gender,
            'national', at.national,
            'country', at.country,
            'type', at.type,
            'player_count', at.player_count,
            'game_id', at.game_id,
            'location_id', at.location_id
        ),

        'awayScore', CASE 
            WHEN g.name = 'football' THEN 
                json_build_object(
                    'id', fs_away.id,
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
            'user_id', t.user_id,
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
            'declared', cs.declared,
            'inning_status', cs.inning_status
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
            'declared', cs.declared,
            'inning_status', cs.inning_status
        ) ORDER BY cs.inning_number
    ) AS scores
    FROM cricket_score cs
    WHERE cs.match_id = m.id AND cs.team_id = at.id
) AS cricket_away_scores ON true

WHERE m.public_id = $1 AND t.game_id = $2;
`

const getMatchModelByPublicIdQuery = `
    SELECT *
    FROM matches
    WHERE public_id = $1;
`

// ListMatches retrieves all matches with enriched data based on date and game filters
func (q *Queries) ListMatches(ctx context.Context, startDate int32, gameID int64) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, listMatchesQuery, startDate, gameID)
	if err != nil {
		log.Printf("Failed to execute ListMatches query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var matches []map[string]interface{}
	for rows.Next() {
		var jsonByte []byte
		if err := rows.Scan(&jsonByte); err != nil {
			log.Printf("Failed to scan match row: %v", err)
			continue
		}

		var match map[string]interface{}
		if err := json.Unmarshal(jsonByte, &match); err != nil {
			log.Printf("Failed to unmarshal match JSON: %v", err)
			continue
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Row iteration error in ListMatches: %v", err)
		return nil, err
	}

	return matches, nil
}

// GetMatchByPublicId retrieves a single match with all related data by public ID
func (q *Queries) GetMatchByPublicId(ctx context.Context, publicId uuid.UUID, gameID int64) (map[string]interface{}, error) {
	var jsonByte []byte
	row := q.db.QueryRowContext(ctx, getMatchByPublicIdQuery, publicId, gameID)
	if err := row.Scan(&jsonByte); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Match not found with public_id: %s, game_id: %d", publicId, gameID)
			return nil, nil
		}
		log.Printf("Failed to scan match detail: %v", err)
		return nil, err
	}

	var match map[string]interface{}
	if err := json.Unmarshal(jsonByte, &match); err != nil {
		log.Printf("Failed to unmarshal match detail JSON: %v", err)
		return nil, err
	}

	return match, nil
}

// GetMatchModelByPublicId retrieves raw match model data by public ID
func (q *Queries) GetMatchModelByPublicId(ctx context.Context, public_id uuid.UUID) (*models.Match, error) {
	var match models.Match

	row := q.db.QueryRowContext(ctx, getMatchModelByPublicIdQuery, public_id)
	if err := row.Scan(
		&match.ID,
		&match.PublicID,
		&match.TournamentID,
		&match.AwayTeamID,
		&match.HomeTeamID,
		&match.StartTimestamp,
		&match.EndTimestamp,
		&match.Type,
		&match.StatusCode,
		&match.Result,
		&match.Stage,
		&match.KnockoutLevelID,
		&match.MatchFormat,
		&match.DayNumber,
		&match.SubStatus,
		&match.LocationID,
		&match.LocationLocked,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Match model not found with public_id: %s", public_id)
			return nil, nil
		}
		log.Printf("Failed to scan match model: %v", err)
		return nil, err
	}

	return &match, nil
}

const getLiveMatches = `
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
            'day_number', m.day_number,
            'sub_status', m.sub_status,
            'location_id', m.location_id,
            'location_locked', m.location_locked,

            'homeTeam', json_build_object(
                'id', ht.id,
                'public_id', ht.public_id,
                'user_id', ht.user_id,
                'name', ht.name,
                'slug', ht.slug,
                'short_name', ht.shortname,
                'media_url', ht.media_url,
                'gender', ht.gender,
                'national', ht.national,
                'country', ht.country,
                'type', ht.type,
                'player_count', ht.player_count,
                'game_id', ht.game_id,
                'location_id', ht.location_id
            ),

            'homeScore', CASE 
                WHEN g.name = 'football' THEN 
                    json_build_object(
                        'id', fs_home.id,
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
                'user_id', at.user_id,
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
                'user_id', t.user_id,
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
            ),

            'location', CASE
                WHEN loc.id IS NOT NULL THEN
                    json_build_object(
                        'id', loc.id,
                        'public_id', loc.public_id,
                        'city', loc.city,
                        'state', loc.state,
                        'country', loc.country,
                        'latitude', loc.latitude,
                        'longitude', loc.longitude,
                        'h3_index', loc.h3_index
                    )
                ELSE NULL
            END
        ) AS response

    FROM matches m
    JOIN teams ht ON m.home_team_id = ht.id
    JOIN teams at ON m.away_team_id = at.id
    LEFT JOIN tournaments t ON m.tournament_id = t.id
    JOIN games g ON t.game_id = g.id
    LEFT JOIN locations loc ON m.location_id = loc.id

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
        ) as scores
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
        ) as scores
        FROM cricket_score cs
        WHERE cs.match_id = m.id AND cs.team_id = at.id
    ) AS cricket_away_scores ON true
    WHERE m.status_code='in_progress'
    AND t.game_id = $1
    ORDER BY m.start_timestamp;
`

func (q *Queries) GetLiveMatches(ctx context.Context, game int64) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getLiveMatches, game)
	if err != nil {
		log.Printf("Failed to query: %v", err)
		return nil, err
	}

	defer rows.Close()

	var liveMatches []map[string]interface{}
	for rows.Next() {
		var liveMatch map[string]interface{}
		var jsonByte []byte
		err := rows.Scan(&jsonByte)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		err = json.Unmarshal(jsonByte, &liveMatch)
		if err != nil {
			log.Printf("Failed to unmarshal: %v", err)
			return nil, err
		}
		liveMatches = append(liveMatches, liveMatch)
	}
	return liveMatches, nil
}

// Update match location
const updateMatchLocationQuery = `
	UPDATE matches
	SET location_id = $2
	WHERE public_id = $1
	RETURNING *;
`

func (q *Queries) UpdateMatchLocation(ctx context.Context, eventPublicID uuid.UUID, locationID int32) (*models.Match, error) {
	var i models.Match
	rows := q.db.QueryRowContext(ctx, updateMatchLocationQuery, eventPublicID, locationID)
	err := rows.Scan(
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
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}
	return &i, err
}

const udpateMatchLocationLockedQuery = `
    UPDATE matches
    SET location_locked = true
    WHERE id = $1
    RETURNING *;
`

func (q *Queries) UpdateMatchLocationLocked(ctx context.Context, matchID int64) (*models.Match, error) {
	var i models.Match
	rows := q.db.QueryRowContext(ctx, udpateMatchLocationLockedQuery, matchID)
	err := rows.Scan(
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
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}
	return &i, err
}
