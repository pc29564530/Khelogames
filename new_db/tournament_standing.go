package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const createTournamentStanding = `
INSERT INTO tournament_standing (
    tournament_id,
    group_id,
    team_id,
    wins,
    loss,
    draw,
    goal_for,
    goal_against,
    goal_difference,
    points
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ) RETURNING standing_id, tournament_id, group_id, team_id, wins, loss, draw, goal_for, goal_against, goal_difference, points
`

type CreateTournamentStandingParams struct {
	TournamentID   int64 `json:"tournament_id"`
	GroupID        int64 `json:"group_id"`
	TeamID         int64 `json:"team_id"`
	Wins           int64 `json:"wins"`
	Loss           int64 `json:"loss"`
	Draw           int64 `json:"draw"`
	GoalFor        int64 `json:"goal_for"`
	GoalAgainst    int64 `json:"goal_against"`
	GoalDifference int64 `json:"goal_difference"`
	Points         int64 `json:"points"`
}

func (q *Queries) CreateTournamentStanding(ctx context.Context, arg CreateTournamentStandingParams) (models.TournamentStanding, error) {
	row := q.db.QueryRowContext(ctx, createTournamentStanding,
		arg.TournamentID,
		arg.GroupID,
		arg.TeamID,
		arg.Wins,
		arg.Loss,
		arg.Draw,
		arg.GoalFor,
		arg.GoalAgainst,
		arg.GoalDifference,
		arg.Points,
	)
	var i models.TournamentStanding
	err := row.Scan(
		&i.StandingID,
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

const getTournamentStanding = `
SELECT 
    ts.standing_id, ts.tournament_id, ts.group_id, ts.team_id,
    ts.wins, ts.loss, ts.draw, ts.goal_for, ts.goal_against, ts.goal_difference, ts.points,
    t.tournament_name, t.sports,
    c.name
FROM 
    tournament_standing ts
JOIN 
    groups tg ON ts.group_id = tg.id
JOIN 
    tournaments t ON ts.tournament_id = t.id
JOIN 
    teams c ON ts.team_id = c.id
WHERE 
    ts.tournament_id = $1
    AND ts.group_id = $2
    AND t.sports = $3
`

type GetTournamentStandingParams struct {
	TournamentID int64  `json:"tournament_id"`
	GroupID      int64  `json:"group_id"`
	Sports       string `json:"sports"`
}

type GetTournamentStandingRow struct {
	StandingID     int64  `json:"standing_id"`
	TournamentID   int64  `json:"tournament_id"`
	GroupID        int64  `json:"group_id"`
	TeamID         int64  `json:"team_id"`
	Wins           int64  `json:"wins"`
	Loss           int64  `json:"loss"`
	Draw           int64  `json:"draw"`
	GoalFor        int64  `json:"goal_for"`
	GoalAgainst    int64  `json:"goal_against"`
	GoalDifference int64  `json:"goal_difference"`
	Points         int64  `json:"points"`
	TournamentName string `json:"tournament_name"`
	Sports         string `json:"sports"`
	Name           string `json:"name"`
}

func (q *Queries) GetTournamentStanding(ctx context.Context, arg GetTournamentStandingParams) ([]GetTournamentStandingRow, error) {
	rows, err := q.db.QueryContext(ctx, getTournamentStanding, arg.TournamentID, arg.GroupID, arg.Sports)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTournamentStandingRow
	for rows.Next() {
		var i GetTournamentStandingRow
		if err := rows.Scan(
			&i.StandingID,
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
			&i.TournamentName,
			&i.Sports,
			&i.Name,
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

const updateTournamentStanding = `
UPDATE tournament_standing AS ts
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

type UpdateTournamentStandingParams struct {
	TournamentID int64 `json:"tournament_id"`
	TeamID       int64 `json:"team_id"`
}

func (q *Queries) UpdateTournamentStanding(ctx context.Context, arg UpdateTournamentStandingParams) (models.TournamentStanding, error) {
	row := q.db.QueryRowContext(ctx, updateTournamentStanding, arg.TournamentID, arg.TeamID)
	var i models.TournamentStanding
	err := row.Scan(
		&i.StandingID,
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
