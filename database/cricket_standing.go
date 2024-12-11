package database

import (
	"context"
	"khelogames/database/models"
)

const createCricketStanding = `
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
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11 ) RETURNING id, tournament_id, group_id, team_id, matches, wins, loss, draw, goal_for, goal_against, goal_difference, points
`

func (q *Queries) CreateCricketStanding(ctx context.Context, tournamentID, groupID, teamID int64) (models.CricketStanding, error) {
	row := q.db.QueryRowContext(ctx, createCricketStanding,
		tournamentID,
		groupID,
		teamID,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	var i models.CricketStanding
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.GroupID,
		&i.TeamID,
		&i.Matches,
		&i.Wins,
		&i.Loss,
		&i.Draw,
		&i.Points,
	)
	return i, err
}

const getCricketStanding = `
	SELECT
		CASE
			WHEN EXISTS (
				SELECT 1 FROM cricket_standing fs
				WHERE cs.tournament_id=$1
			) THEN
				JSON_AGG(
					JSON_BUILD_OBJECT(
						cs.id, cs.tournament_id, cs.group_id, cs.team_id, cs.matches, cs.wins, cs.loss, cs.draw, cs.points,
						JSON_BUILD_OBJECT(
							'tournament', JSON_BUILD_OBJECT('id', t.id, 'name', t.name, 'slug', t.slug, 'country', t.country, 'status_code', t.status_code, 'level', t.level, 'start_timestamp', t.start_timestamp, 'game_id', t.game_id),
							'group', CASE WHEN g.id IS NOT NULL THEN JSON_BUILD_OBJECT('id', g.id, 'name', g.name) ELSE NULL END,
							'teams', JSON_BUILD_OBJECT('id', tm.id, 'name', tm.name, 'slug', tm.slug, 'short_name', tm.shortname, 'admin', tm.admin, 'media_url', tm.media_url, 'gender', tm.gender, 'national', tm.national, 'country', tm.country, 'type', tm.type, 'player_count', tm.player_count, 'game_id', tm.game_id)
						)
					)
				) ELSE NULL
		END AS standing_data
	FROM cricket_standing fs
	LEFT JOIN groups g ON cs.group_id = g.id
	JOIN tournaments t ON t.id = cs.tournament_id
	JOIN teams tm ON cs.team_id = tm.id
	WHERE cs.tournament_id = $1
`

type GetCricketStandingR struct {
	StandingData interface{} `json:"standing_data"`
}

func (q *Queries) GetCricketStanding(ctx context.Context, tournamentId int64) (*GetCricketStandingR, error) {
	rows, err := q.db.QueryContext(ctx, getCricketStanding, tournamentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var standings GetCricketStandingR

	if rows.Next() {
		if err := rows.Scan(&standings.StandingData); err != nil {
			return nil, err
		}
	}

	// for rows.Next() {
	// 	var i GetCricketStandingR
	// 	if err := rows.Scan(&i.StandingData); err != nil {
	// 		return nil, err
	// 	}
	// 	standings = append(standings, i)
	// }
	// if err := rows.Close(); err != nil {
	// 	return nil, err
	// }
	if err := rows.Err(); err != nil {
		return nil, err
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
WHERE ts.tournament_id = $1 AND ts.team_id = $2 AND (m.stage = 'group' OR m.stage = 'league')
RETURNING standing_id, tournament_id, group_id, team_id, wins, loss, draw, goal_for, goal_against, goal_difference, points
`

func (q *Queries) UpdateCricketStanding(ctx context.Context, tournamentID, teamID int64) (models.CricketStanding, error) {
	row := q.db.QueryRowContext(ctx, updateCricketStanding, tournamentID, teamID)
	var i models.CricketStanding
	err := row.Scan(
		&i.ID,
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
