package database

import (
	"context"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const createFootballStanding = `
WITH tournamentID AS (
    SELECT * FROM tournaments WHERE public_id = $1
),
teamID AS (
    SELECT * FROM teams WHERE public_id = $3
)
INSERT INTO football_standing (
    tournament_id,
    group_id,
    team_id,
	matches,
    wins,
    loss,
    draw,
    goal_for,
    goal_against,
    goal_difference,
    points
)
SELECT 
    tournamentID.id,
    $2,
    teamID.id,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11  
FROM tournamentID, teamID
RETURNING *;
`

func (q *Queries) CreateFootballStanding(ctx context.Context, tournamentPublicID uuid.UUID, groupID int32, teamPublicID uuid.UUID) (models.FootballStanding, error) {

	row := q.db.QueryRowContext(ctx, createFootballStanding,
		tournamentPublicID,
		groupID,
		teamPublicID,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	var i models.FootballStanding
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.TournamentID,
		&i.GroupID,
		&i.TeamID,
		&i.Matches,
		&i.Wins,
		&i.Loss,
		&i.Draw,
		&i.GoalFor,
		&i.GoalAgainst,
		&i.GoalDifference,
		&i.Points,
	)
	return i, err
}

const getFootballStanding = `
	SELECT
    CASE
        WHEN EXISTS (
            SELECT 1
            FROM football_standing fs
            JOIN tournaments t ON t.id = fs.tournament_id
            WHERE t.public_id = $1
        ) THEN 
            JSON_AGG(
                JSON_BUILD_OBJECT(
                    'id', fs.id,
                    'public_id', fs.public_id,
                    'tournament_id', fs.tournament_id,
                    'group_id', CASE
                        WHEN fs.group_id IS NOT NULL THEN fs.group_id
                        ELSE NULL
                    END,
                    'team_id', fs.team_id,
                    'matches', fs.matches,
                    'wins', fs.wins,
                    'loss', fs.loss,
                    'draw', fs.draw,
                    'goal_for', fs.goal_for,
                    'goal_against', fs.goal_against,
                    'goal_difference', fs.goal_difference,
                    'points', fs.points,
                    'tournament', JSON_BUILD_OBJECT(
                        'id', t.id,
                        'public_id', t.public_id,
                        'user_id', t.user_id,
                        'name', t.name,
                        'slug', t.slug,
                        'country', t.country,
                        'status', t.status,
                        'level', t.level,
                        'start_timestamp', t.start_timestamp,
                        'game_id', t.game_id
                    ),
                    'group', CASE 
                        WHEN g.id IS NOT NULL THEN JSON_BUILD_OBJECT(
                            'id', g.id,
                            'name', g.name
                        ) 
                        ELSE NULL 
                    END,
                    'teams', JSON_BUILD_OBJECT(
                        'id', tm.id,
                        'public_id', tm.public_id,
                        'user_id', tm.user_id,
                        'name', tm.name,
                        'slug', tm.slug,
                        'short_name', tm.shortname,
                        'media_url', tm.media_url,
                        'gender', tm.gender,
                        'national', tm.national,
                        'country', tm.country,
                        'type', tm.type,
                        'player_count', tm.player_count,
                        'game_id', tm.game_id
                    )
                )
            )
			ELSE NULL
		END AS standing_data
	FROM football_standing fs
	LEFT JOIN groups g ON fs.group_id = g.id
	JOIN tournaments t ON t.id = fs.tournament_id
	JOIN teams tm ON fs.team_id = tm.id
	WHERE t.public_id = $1;
`

type GetFootballStandingR struct {
	StandingData interface{} `json:"standing_data"`
}

func (q *Queries) GetFootballStanding(ctx context.Context, tournamentPublicID uuid.UUID) (*GetFootballStandingR, error) {
	rows, err := q.db.QueryContext(ctx, getFootballStanding, tournamentPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var standings GetFootballStandingR

	if rows.Next() {
		if err := rows.Scan(&standings.StandingData); err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return &standings, nil
}

const updateFootballStanding = `
UPDATE football_standing AS ts
SET 
    matches = (
        SELECT COUNT(*)
        FROM matches ms
        LEFT JOIN football_score fs_home 
            ON fs_home.match_id = ms.id 
            AND fs_home.team_id = ms.home_team_id
        LEFT JOIN football_score fs_away 
            ON fs_away.match_id = ms.id 
            AND fs_away.team_id = ms.away_team_id
        WHERE (ms.home_team_id = ts.team_id OR ms.away_team_id = ts.team_id)
          AND ms.tournament_id = ts.tournament_id
          AND (LOWER(ms.stage) = 'group' OR LOWER(ms.stage) = 'league')
          AND ms.status_code = 'finished'
    ),

    goal_for = COALESCE((
        SELECT SUM(CASE
            WHEN ms.home_team_id = ts.team_id THEN fs_home.goals
            WHEN ms.away_team_id = ts.team_id THEN fs_away.goals
            ELSE 0
        END)
        FROM matches ms
        LEFT JOIN football_score fs_home 
            ON fs_home.match_id = ms.id 
            AND fs_home.team_id = ms.home_team_id
        LEFT JOIN football_score fs_away 
            ON fs_away.match_id = ms.id 
            AND fs_away.team_id = ms.away_team_id
        WHERE (ms.home_team_id = ts.team_id OR ms.away_team_id = ts.team_id)
          AND ms.tournament_id = ts.tournament_id
          AND (LOWER(ms.stage) = 'group' OR LOWER(ms.stage) = 'league')
          AND ms.status_code = 'finished'
    ), 0),

    goal_against = COALESCE((
        SELECT SUM(CASE
            WHEN ms.home_team_id = ts.team_id THEN fs_away.goals
            WHEN ms.away_team_id = ts.team_id THEN fs_home.goals
            ELSE 0
        END)
        FROM matches ms
        LEFT JOIN football_score fs_home 
            ON fs_home.match_id = ms.id 
            AND fs_home.team_id = ms.home_team_id
        LEFT JOIN football_score fs_away 
            ON fs_away.match_id = ms.id 
            AND fs_away.team_id = ms.away_team_id
        WHERE (ms.home_team_id = ts.team_id OR ms.away_team_id = ts.team_id)
          AND ms.tournament_id = ts.tournament_id
          AND (LOWER(ms.stage) = 'group' OR LOWER(ms.stage) = 'league')
          AND ms.status_code = 'finished'   -- ✅ Added missing AND
    ), 0),

    goal_difference = COALESCE(goal_for, 0) - COALESCE(goal_against, 0),

    wins = COALESCE((
        SELECT COUNT(*)
        FROM matches ms
        LEFT JOIN football_score fs_home 
            ON fs_home.match_id = ms.id 
            AND fs_home.team_id = ms.home_team_id
        LEFT JOIN football_score fs_away 
            ON fs_away.match_id = ms.id 
            AND fs_away.team_id = ms.away_team_id
        WHERE
        (
            (ms.home_team_id = ts.team_id AND fs_home.goals > fs_away.goals)
            OR
            (ms.away_team_id = ts.team_id AND fs_away.goals > fs_home.goals)
        )
        AND ms.tournament_id = ts.tournament_id
        AND (LOWER(ms.stage) = 'group' OR LOWER(ms.stage) = 'league')
        AND ms.status_code = 'finished'
    ), 0),

    loss = COALESCE((
        SELECT COUNT(*)
        FROM matches ms
        LEFT JOIN football_score fs_home 
            ON fs_home.match_id = ms.id 
            AND fs_home.team_id = ms.home_team_id
        LEFT JOIN football_score fs_away 
            ON fs_away.match_id = ms.id 
            AND fs_away.team_id = ms.away_team_id
        WHERE
        (
            (ms.home_team_id = ts.team_id AND fs_home.goals < fs_away.goals)
            OR
            (ms.away_team_id = ts.team_id AND fs_away.goals < fs_home.goals)
        )
        AND ms.tournament_id = ts.tournament_id
        AND (LOWER(ms.stage) = 'group' OR LOWER(ms.stage) = 'league')
        AND ms.status_code = 'finished'
    ), 0),

    draw = COALESCE((
        SELECT COUNT(*)
        FROM matches ms
        LEFT JOIN football_score fs_home 
            ON fs_home.match_id = ms.id 
            AND fs_home.team_id = ms.home_team_id
        LEFT JOIN football_score fs_away 
            ON fs_away.match_id = ms.id 
            AND fs_away.team_id = ms.away_team_id
        WHERE
        (
            (ms.home_team_id = ts.team_id AND fs_home.goals = fs_away.goals)
            OR
            (ms.away_team_id = ts.team_id AND fs_away.goals = fs_home.goals)
        )
        AND ms.tournament_id = ts.tournament_id
        AND (LOWER(ms.stage) = 'group' OR LOWER(ms.stage) = 'league')
        AND ms.status_code = 'finished'
    ), 0),

    points = ((wins * 3) + draw)

FROM tournaments t, teams tm
WHERE ts.tournament_id = t.id
  AND ts.team_id = tm.id
  AND t.id = $1
  AND tm.id = $2
RETURNING ts.*;
`

func (q *Queries) UpdateFootballStanding(ctx context.Context, tournamentID, teamID int64) (*models.FootballStanding, error) {
	row := q.db.QueryRowContext(ctx, updateFootballStanding, tournamentID, teamID)
	var i models.FootballStanding
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.TournamentID,
		&i.GroupID,
		&i.TeamID,
		&i.Matches,
		&i.Wins,
		&i.Loss,
		&i.Draw,
		&i.GoalFor,
		&i.GoalAgainst,
		&i.GoalDifference,
		&i.Points,
	)
	return &i, err
}
