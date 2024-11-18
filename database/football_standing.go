package database

import (
	"context"
	"khelogames/database/models"
)

const createFootballStanding = `
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

func (q *Queries) CreateFootballStanding(ctx context.Context, tournamentID, groupID, teamID int64) (models.FootballStanding, error) {
	row := q.db.QueryRowContext(ctx, createFootballStanding,
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
	var i models.FootballStanding
	err := row.Scan(
		&i.ID,
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
            WHERE tournament_id = $1
        ) THEN 
            JSON_AGG(
                JSON_BUILD_OBJECT(
                    'id', fs.id,
                    'tournament_id', fs.tournament_id,
                    'group_id', fs.group_id,
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
                        'name', t.name,
                        'slug', t.slug,
                        'country', t.country,
                        'status_code', t.status_code,
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
                        'name', tm.name,
                        'slug', tm.slug,
                        'short_name', tm.shortname,
                        'admin', tm.admin,
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
	WHERE fs.tournament_id = $1;
`

type GetFootballStandingR struct {
	StandingData interface{} `json:"standing_data"`
}

func (q *Queries) GetFootballStanding(ctx context.Context, tournamentId int64) (*GetFootballStandingR, error) {
	rows, err := q.db.QueryContext(ctx, getFootballStanding, tournamentId)
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
	// for rows.Next() {
	// 	var i GetFootballStandingR
	// 	if err := rows.Scan(&i.StandingData); err != nil {
	// 		return nil, err
	// 	}
	// 	standings = append(standings, i)
	// }
	// if err := rows.Close(); err != nil {
	// 	return nil, err
	// }
	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }
	return &standings, nil
}

const updateFootballStanding = `
UPDATE football_standing AS ts
SET 
    goal_for = COALESCE((
        SELECT SUM(CASE 
            WHEN ms.home_team_id = ts.team_id THEN fs.goals
            WHEN ms.away_team_id = ts.team_id THEN fs.goals
            ELSE 0
        END)
        FROM football_score AS fs
        JOIN matches AS ms ON fs.match_id = ms.id
        WHERE fs.team_id = ts.team_id
    ), 0),
    goal_against = COALESCE((
        SELECT SUM(CASE 
            WHEN ms.home_team_id = ts.team_id THEN (
                SELECT SUM(fs.goals) 
                FROM football_score AS fs 
                WHERE fs.match_id = ms.id AND fs.team_id = ms.away_team_id
            )
            WHEN ms.away_team_id = ts.team_id THEN (
                SELECT SUM(fs2.goals) 
                FROM football_score AS fs2
                WHERE fs2.match_id = ms.id AND fs2.team_id = ms.home_team_id
            )
        END)
        FROM matches AS ms
        WHERE ms.home_team_id = ts.team_id OR ms.away_team_id = ts.team_id
    ), 0),
    goal_difference = COALESCE(goal_for, 0) - COALESCE(goal_against, 0),
    wins = COALESCE((
        SELECT COUNT(*)
        FROM matches AS ms
        LEFT JOIN football_score AS fs_home ON ms.id = fs_home.match_id AND ms.home_team_id = ts.team_id
        LEFT JOIN football_score AS fs_away ON ms.id = fs_away.match_id AND ms.away_team_id = fs_away.team_id
        WHERE (ms.home_team_id = ts.team_id AND fs_home.goals > fs_away.goals)
        OR (ms.away_team_id = ts.team_id AND fs_away.goals > fs_home.goals)
    ), 0),
    loss = COALESCE((
        SELECT COUNT(*)
        FROM matches AS ms
        LEFT JOIN football_score fs_home ON ms.id = fs_home.match_id AND ms.home_team_id = fs_home.team_id
        LEFT JOIN football_score fs_away ON ms.id = fs_away.match_id AND ms.away_team_id = fs_away.team_id
        WHERE (ms.home_team_id = ts.team_id AND fs_home.goals < fs_away.goals)
        OR (ms.away_team_id = ts.team_id AND fs_away.goals < fs_home.goals)
    ), 0),
    draw = COALESCE((
        SELECT COUNT(*)
        FROM matches AS ms
        LEFT JOIN football_score AS fs_home ON ms.id = fs_home.match_id AND ms.home_team_id = ts.team_id
        LEFT JOIN football_score AS fs_away ON ms.id = fs_away.match_id AND ms.away_team_id = fs_away.team_id
        WHERE (ms.home_team_id = ts.team_id AND fs_home.goals = fs_away.goals)
        OR (ms.away_team_id = ts.team_id AND fs_away.goals = fs_home.goals)
    ), 0),
    points = ((wins * 3) + draw)
WHERE ts.tournament_id = $1 AND ts.team_id = $2
RETURNING standing_id, tournament_id, group_id, team_id, wins, loss, draw, goal_for, goal_against, goal_difference, points
`

func (q *Queries) UpdateFootballStanding(ctx context.Context, tournamentID, teamID int64) (models.FootballStanding, error) {
	row := q.db.QueryRowContext(ctx, updateFootballStanding, tournamentID, teamID)
	var i models.FootballStanding
	err := row.Scan(
		&i.ID,
		&i.TournamentID,
		&i.GroupID,
		&i.TeamID,
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
