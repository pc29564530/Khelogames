package database

import (
	"context"
	"encoding/json"
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
    tournamentID.tournament_id,
    $2,
    teamID.team_id,
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
		JSON_BUILD_OBJECT(
			'id', cs.id,
			'public_id', cs.public_id,
			'tournament_id', cs.tournament_id, 
			'group_id', CASE
                        WHEN cs.group_id IS NOT NULL THEN cs.group_id
                        ELSE NULL
                    END, 
			'team_id', cs.team_id, 
			'matches', COALESCE(cs.matches,0), 
			'wins', COALESCE(cs.wins,0), 
			'loss', COALESCE(cs.loss,0), 
			'draw', COALESCE(cs.draw,0), 
			'point', COALESCE(cs.points,0),
			'details', JSON_BUILD_OBJECT(
				'tournament', JSON_BUILD_OBJECT('id', t.id, 'public_id', t.public_id 'user_id', t.user_id, 'name', t.name, 'slug', t.slug, 'country', t.country, 'status_code', t.status_code, 'level', t.level, 'start_timestamp', t.start_timestamp, 'game_id', t.game_id, 'group_count', t.group_count, 'max_group_team', t.max_group_teams, 'stage', t.stage, 'has_knockout', t.has_knockout),
				'group', CASE WHEN g.id IS NOT NULL THEN JSON_BUILD_OBJECT('id', g.id, 'name', g.name) ELSE NULL END,
				'teams', JSON_BUILD_OBJECT('id', tm.id, 'public_id', tm.public_id, 'user_id', tm.user_id, 'name', tm.name, 'slug', tm.slug, 'short_name', tm.shortname, 'media_url', tm.media_url, 'gender', tm.gender, 'national', tm.national, 'country', tm.country, 'type', tm.type, 'player_count', tm.player_count, 'game_id', tm.game_id)
			)
		) AS standing_data
	FROM cricket_standing cs
	LEFT JOIN groups g ON cs.group_id = g.id
	JOIN tournaments t ON t.id = cs.tournament_id
	JOIN teams tm ON cs.team_id = tm.id
	WHERE t.public_id = $1
`

type GetCricketStandingR struct {
	StandingData interface{} `json:"standing_data"`
}

func (q *Queries) GetCricketStanding(ctx context.Context, tournamentPublicID uuid.UUID) (*[]GetCricketStandingR, error) {
	rows, err := q.db.QueryContext(ctx, getCricketStanding, tournamentPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var standings []GetCricketStandingR
	if rows.Next() {
		var standing GetCricketStandingR
		var jsonData []byte
		if err := rows.Scan(
			&jsonData,
		); err != nil {
			return nil, err
		}

		err = json.Unmarshal(jsonData, &standing.StandingData)
		if err != nil {
			return nil, err
		}
		standings = append(standings, standing)
	}
	return &standings, nil
}

const updateCricketStanding = `
UPDATE Cricket_standing AS ts
SET 
    score = COALESCE((
        SELECT SUM(CASE 
            WHEN ms.home_team_id = ts.team_id THEN cs.score
            WHEN ms.away_team_id = ts.team_id THEN cs.score
            ELSE 0
        END)
        FROM cricket_score AS cs
        JOIN matches AS ms ON cs.match_id = ms.id
        WHERE cs.team_id = ts.team_id
    ), 0),
    wickets = COALESCE((
        SELECT SUM(CASE 
            WHEN ms.home_team_id = ts.team_id THEN (
                SELECT SUM(cs.goals) 
                FROM cricket_score AS cs
                WHERE cs.match_id = ms.id AND cs.team_id = ms.away_team_id
            )
            WHEN ms.away_team_id = ts.team_id THEN (
                SELECT SUM(cs.score) 
                FROM cricket_score AS cs
                WHERE cs.match_id = ms.id AND cs.team_id = ms.home_team_id
            )
        END)
        FROM matches AS ms
        WHERE ms.home_team_id = ts.team_id OR ms.away_team_id = ts.team_id
    ), 0),
    wins = COALESCE((
        SELECT COUNT(*)
        FROM matches AS ms
        LEFT JOIN cricket_score AS cs.home ON ms.id = cs.home.match_id AND ms.home_team_id = ts.team_id
        LEFT JOIN cricket_score AS cs.away ON ms.id = cs.away.match_id AND ms.away_team_id = cs.away.team_id
        WHERE (ms.home_team_id = ts.team_id AND cs.home.goals > cs.away.goals)
        OR (ms.away_team_id = ts.team_id AND cs.away.goals > cs.home.goals)
    ), 0),
    loss = COALESCE((
        SELECT COUNT(*)
        FROM matches AS ms
        LEFT JOIN cricket_score cs.home ON ms.id = cs.home.match_id AND ms.home_team_id = cs.home.team_id
        LEFT JOIN cricket_score cs.away ON ms.id = cs.away.match_id AND ms.away_team_id = cs.away.team_id
        WHERE (ms.home_team_id = ts.team_id AND cs.home.goals < cs.away.goals)
        OR (ms.away_team_id = ts.team_id AND cs.away.goals < cs.home.goals)
    ), 0),
    draw = COALESCE((
        SELECT COUNT(*)
        FROM matches AS ms
        LEFT JOIN cricket_score AS cs.home ON ms.id = cs.home.match_id AND ms.home_team_id = ts.team_id
        LEFT JOIN cricket_score AS cs.away ON ms.id = cs.away.match_id AND ms.away_team_id = cs.away.team_id
        WHERE (ms.home_team_id = ts.team_id AND cs.home.goals = cs.away.goals)
        OR (ms.away_team_id = ts.team_id AND cs.away.goals = cs.home.goals)
    ), 0),
    points = ((wins * 3) + draw)
WHERE ts.tournament_id = t.id
  AND ts.team_id = tm.id
  AND t.public_id = $1
  AND tm.public_id = $2
RETURNING *
`

func (q *Queries) UpdateCricketStanding(ctx context.Context, tournamentPublicID, teamPublicID uuid.UUID) (models.CricketStanding, error) {
	row := q.db.QueryRowContext(ctx, updateCricketStanding, tournamentPublicID, teamPublicID)
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
