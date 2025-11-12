package database

import (
	"context"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const createCricketStanding = `
WITH tournamentID AS (
    SELECT * FROM tournaments WHERE public_id = $1
),
teamID AS (
    SELECT * FROM teams WHERE public_id = $3
)
INSERT INTO cricket_standing (
    tournament_id,
    group_id,
    team_id,
	matches,
    wins,
    loss,
    draw,
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
    $8
FROM tournamentID, teamID
RETURNING *;
`

func (q *Queries) CreateCricketStanding(ctx context.Context, tournamentPublicID uuid.UUID, groupID int32, teamPublicID uuid.UUID) (*models.CricketStanding, error) {
	row := q.db.QueryRowContext(ctx, createCricketStanding,
		tournamentPublicID,
		groupID,
		teamPublicID,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	var i models.CricketStanding
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
		&i.Points,
	)
	return &i, err
}

const getCricketStanding = `
	SELECT
	CASE
		WHEN EXISTS (
			SELECT 1
			FROM cricket_standing cs
			JOIN tournaments t ON t.id = cs.tournament_id
			WHERE t.public_id = $1
		) THEN 
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'id', cs.id,
					'public_id', cs.public_id,
					'tournament_id', cs.tournament_id,
					'group_id', CASE
						WHEN cs.group_id IS NOT NULL THEN cs.group_id
						ELSE NULL
					END,
					'team_id', cs.team_id,
					'matches', cs.matches,
					'wins', cs.wins,
					'loss', cs.loss,
					'draw', cs.draw,
					'points', cs.points,
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
	FROM cricket_standing cs
	LEFT JOIN groups g ON cs.group_id = g.id
	JOIN tournaments t ON t.id = cs.tournament_id
	JOIN teams tm ON cs.team_id = tm.id
	WHERE t.public_id = $1;
`

type GetCricketStandingR struct {
	StandingData interface{} `json:"standing_data"`
}

func (q *Queries) GetCricketStanding(ctx context.Context, tournamentPublicID uuid.UUID) (*GetCricketStandingR, error) {
	rows, err := q.db.QueryContext(ctx, getCricketStanding, tournamentPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var standings GetCricketStandingR

	if rows.Next() {
		if err := rows.Scan(&standings.StandingData); err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return &standings, nil
}

const updateCricketStanding = `
UPDATE cricket_standing AS ts
SET 
	matches = (
        SELECT COUNT(*)
        FROM matches ms
        LEFT JOIN cricket_score cs_home
            ON cs_home.match_id = ms.id
            AND cs_home.team_id = ms.home_team_id
        LEFT JOIN football_score cs_away
            ON cs_away.match_id = ms.id
            AND cs_away.team_id = ms.away_team_id
        WHERE (ms.home_team_id = ts.team_id OR ms.away_team_id = ts.team_id)
        AND ms.tournament_id = ts.tournament_id
        AND (LOWER(ms.stage) = 'group' OR LOWER(ms.stage) = 'league')
        AND (ms.status_code) = 'finished'
    ),
    wins = COALESCE((
        SELECT COUNT(*)
        FROM matches AS ms
        LEFT JOIN cricket_score AS cs_home ON ms.id = cs_home.match_id AND ms.home_team_id = ts.team_id
        LEFT JOIN cricket_score AS cs_away ON ms.id = cs_away.match_id AND ms.away_team_id = cs_away.team_id
        WHERE ((ms.home_team_id = ts.team_id AND cs_home.score > cs_away.score)
        OR (ms.away_team_id = ts.team_id AND cs_away.score > cs_home.score))
		AND (ms.status_code) = 'finished'
    ), 0),
    loss = COALESCE((
        SELECT COUNT(*)
        FROM matches AS ms
        LEFT JOIN cricket_score cs_home ON ms.id = cs_home.match_id AND ms.home_team_id = cs_home.team_id
        LEFT JOIN cricket_score cs_away ON ms.id = cs_away.match_id AND ms.away_team_id = cs_away.team_id
        WHERE ((ms.home_team_id = ts.team_id AND cs_home.score < cs_away.score)
        OR (ms.away_team_id = ts.team_id AND cs_away.score < cs_home.score))
		AND (ms.status_code) = 'finished'
    ), 0),
    draw = COALESCE((
        SELECT COUNT(*)
        FROM matches AS ms
        LEFT JOIN cricket_score AS cs_home ON ms.id = cs_home.match_id AND ms.home_team_id = ts.team_id
        LEFT JOIN cricket_score AS cs_away ON ms.id = cs_away.match_id AND ms.away_team_id = cs_away.team_id
        WHERE ((ms.home_team_id = ts.team_id AND cs_home.score = cs_away.score)
        OR (ms.away_team_id = ts.team_id AND cs_away.score = cs_home.score))
		AND (ms.status_code) = 'finished'
    ), 0),
    points = ((wins * 3) + draw)
WHERE ts.tournament_id = t.id
  AND ts.team_id = tm.id
  AND t.id = $1
  AND tm.id = $2
RETURNING *
`

func (q *Queries) UpdateCricketStanding(ctx context.Context, tournamentID, teamID int32) (models.CricketStanding, error) {
	row := q.db.QueryRowContext(ctx, updateCricketStanding, tournamentID, teamID)
	var i models.CricketStanding
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.TournamentID,
		&i.GroupID,
		&i.TeamID,
		&i.Wins,
		&i.Loss,
		&i.Draw,
		&i.Points,
	)
	return i, err
}
