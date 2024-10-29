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
		DISTINCT ON (COALESCE(fs.group_id, fs.id))
		fs.id, fs.tournament_id, fs.group_id, fs.team_id, fs.matches, fs.wins, fs.loss, fs.draw, fs.goal_for, fs.goal_against, fs.goal_difference, fs.points,
		JSON_BUILD_OBJECT(
			'tournament', JSON_BUILD_OBJECT('id', t.id, 'name', t.name, 'slug', t.slug, 'country', t.country, 'status_code', t.status_code, 'level', t.level, 'start_timestamp', t.start_timestamp, 'game_id', t.game_id),
			'group', CASE WHEN g.id IS NOT NULL THEN JSON_BUILD_OBJECT('id', g.id, 'name', g.name) ELSE NULL END,
			'teams', JSON_BUILD_OBJECT('id', tm.id, 'name', tm.name, 'slug', tm.slug, 'short_name', tm.shortname, 'admin', tm.admin, 'media_url', tm.media_url, 'gender', tm.gender, 'national', tm.national, 'country', tm.country, 'type', tm.type, 'player_count', tm.player_count, 'game_id', tm.game_id)
		) AS standing_data
	FROM cricket_standing fs
	LEFT JOIN groups g ON fs.group_id = g.id
	JOIN tournaments t ON t.id = fs.tournament_id
	JOIN teams tm ON fs.team_id = tm.id
	WHERE fs.tournament_id = $1
	ORDER BY COALESCE(fs.group_id, fs.id), fs.id
`

type GetCricketStandingR struct {
	ID           int64       `json:"id"`
	TournamentID int64       `json:"tournament_id"`
	GroupID      *int64      `json:"group_id"`
	TeamID       int64       `json:"team_id"`
	Matches      *int64      `json:"json:"matches"`
	Wins         *int64      `json:"wins"`
	Loss         *int64      `json:"loss"`
	Draw         *int64      `json:"draw"`
	Points       *int64      `json:"points"`
	StandingData interface{} `json:"standing_data"`
}

func (q *Queries) GetCricketStanding(ctx context.Context, tournamentId int64) ([]GetCricketStandingR, error) {
	rows, err := q.db.QueryContext(ctx, getCricketStanding, tournamentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var standings []GetCricketStandingR
	for rows.Next() {
		var i GetCricketStandingR
		if err := rows.Scan(&i.ID, &i.TournamentID, &i.GroupID, &i.TeamID, &i.Matches, &i.Wins, &i.Loss, &i.Draw, &i.Points, &i.StandingData); err != nil {
			return nil, err
		}
		standings = append(standings, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return standings, nil
}

const updateCricketStanding = `
UPDATE Cricket_standing AS ts
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
