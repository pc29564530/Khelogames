package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"
)

const getCricketScore = `
SELECT * FROM cricket_score
WHERE match_id=$1 AND team_id=$2
`

type GetCricketScoreParams struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (q *Queries) GetCricketScore(ctx context.Context, arg GetCricketScoreParams) (models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, getCricketScore, arg.MatchID, arg.TeamID)
	var i models.CricketScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.Inning,
		&i.Score,
		&i.Wickets,
		&i.Overs,
		&i.RunRate,
		&i.TargetRunRate,
		&i.FollowOn,
		&i.IsInningCompleted,
		&i.Declared,
	)
	return i, err
}

const getCricketScores = `
SELECT * FROM cricket_score
WHERE match_id=$1;
`

func (q *Queries) GetCricketScores(ctx context.Context, matchID int64) ([]models.CricketScore, error) {
	row, err := q.db.QueryContext(ctx, getCricketScore, matchID)
	var cricketScores []models.CricketScore

	if row.Next() {
		var i models.CricketScore
		err := row.Scan(
			&i.ID,
			&i.MatchID,
			&i.TeamID,
			&i.Inning,
			&i.Score,
			&i.Wickets,
			&i.Overs,
			&i.RunRate,
			&i.TargetRunRate,
			&i.FollowOn,
			&i.IsInningCompleted,
			&i.Declared,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("Not Able to scan rows: ", err)
		}
		cricketScores = append(cricketScores, i)

	}
	return cricketScores, err
}

const newCricketScore = `
INSERT INTO cricket_score (
    match_id,
    team_id,
    inning,
    score,
    wickets,
    overs,
    run_rate,
    target_run_rate,
	follow_on,
	is_inning_completed,
	declared
) VALUES ( $1, $2, $3, $4, $5, $6, CAST($7 AS numeric(5,2)), CAST($8 AS numeric(5,2)), follow_on,
is_inning_completed, declared)
RETURNING *
`

type NewCricketScoreParams struct {
	MatchID           int64  `json:"match_id"`
	TeamID            int64  `json:"team_id"`
	Inning            string `json:"inning"`
	Score             int32  `json:"score"`
	Wickets           int32  `json:"wickets"`
	Overs             int32  `json:"overs"`
	RunRate           string `json:"run_rate"`
	TargetRunRate     string `json:"target_run_rate"`
	FollowOn          bool   `json:"follow_on"`
	IsInningCompleted bool   `json:"is_inning_completed"`
	Declared          bool   `json:"declared"`
}

func (q *Queries) NewCricketScore(ctx context.Context, arg NewCricketScoreParams) (models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, newCricketScore,
		arg.MatchID,
		arg.TeamID,
		arg.Inning,
		arg.Score,
		arg.Wickets,
		arg.Overs,
		arg.RunRate,
		arg.TargetRunRate,
		arg.FollowOn,
		arg.IsInningCompleted,
		arg.Declared,
	)
	var i models.CricketScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.Inning,
		&i.Score,
		&i.Wickets,
		&i.Overs,
		&i.RunRate,
		&i.TargetRunRate,
		&i.FollowOn,
		&i.IsInningCompleted,
		&i.Declared,
	)
	return i, err
}

const updateCricketInnings = `
UPDATE cricket_score
SET inning=$1
WHERE match_id=$2 AND team_id=$3
RETURNING *
`

type UpdateCricketInningsParams struct {
	Inning  string `json:"inning"`
	MatchID int64  `json:"match_id"`
	TeamID  int64  `json:"team_id"`
}

func (q *Queries) UpdateCricketInnings(ctx context.Context, arg UpdateCricketInningsParams) (models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, updateCricketInnings, arg.Inning, arg.MatchID, arg.TeamID)
	var i models.CricketScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.Inning,
		&i.Score,
		&i.Wickets,
		&i.Overs,
		&i.RunRate,
		&i.TargetRunRate,
		&i.FollowOn,
		&i.IsInningCompleted,
		&i.Declared,
	)
	return i, err
}

const updateCricketOvers = `
UPDATE cricket_score cs
SET overs = (
        SELECT COALESCE(SUM(bl.balls_faced), 0)
        FROM bats bl
        WHERE bl.match_id = cs.match_id AND bl.team_id = cs.team_id
    )
WHERE cs.match_id = $1 AND cs.team_id = $2
RETURNING *;

`

type UpdateCricketOversParams struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (q *Queries) UpdateCricketOvers(ctx context.Context, arg UpdateCricketOversParams) (models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, updateCricketOvers, arg.MatchID, arg.TeamID)
	var i models.CricketScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.Inning,
		&i.Score,
		&i.Wickets,
		&i.Overs,
		&i.RunRate,
		&i.TargetRunRate,
		&i.FollowOn,
		&i.IsInningCompleted,
		&i.Declared,
	)
	return i, err
}

const updateCricketScore = `
UPDATE cricket_score cs
SET score = (
        SELECT SUM(bt.runs_scored) + SUM(bl.wide + bl.no_ball)
        FROM bats bt
		LEFT JOIN balls AS bl ON bt.match_id=bl.match_id AND bl.team_id =  $3 AND bl.bowling_status=true
        WHERE bt.match_id=cs.match_id AND bt.team_id=cs.team_id
        GROUP BY (bt.match_id, bt.team_id)
    )
WHERE cs.match_id=$1 AND cs.team_id=$2
RETURNING *;
`

type UpdateCricketScoreParams struct {
	MatchID       int64 `json:"match_id"`
	BattingTeamID int64 `json:"batting_team_id"`
	BowlingTeamID int64 `json:"bowling_team_id"`
}

func (q *Queries) UpdateCricketScore(ctx context.Context, arg UpdateCricketScoreParams) (models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, updateCricketScore, arg.MatchID, arg.BattingTeamID, arg.BowlingTeamID)
	var i models.CricketScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.Inning,
		&i.Score,
		&i.Wickets,
		&i.Overs,
		&i.RunRate,
		&i.TargetRunRate,
		&i.FollowOn,
		&i.IsInningCompleted,
		&i.Declared,
	)
	return i, err
}

const updateCricketWickets = `
UPDATE cricket_score cs
SET wickets = (
        SELECT COUNT(*)
        FROM wickets w
        WHERE w.match_id = cs.match_id AND w.team_id = cs.team_id
    )
WHERE cs.match_id = $1 AND cs.team_id = $2
RETURNING *;
`

type UpdateCricketWicketsParams struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (q *Queries) UpdateCricketWickets(ctx context.Context, arg UpdateCricketWicketsParams) (models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, updateCricketWickets, arg.MatchID, arg.TeamID)
	var i models.CricketScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.Inning,
		&i.Score,
		&i.Wickets,
		&i.Overs,
		&i.RunRate,
		&i.TargetRunRate,
		&i.FollowOn,
		&i.IsInningCompleted,
		&i.Declared,
	)
	return i, err
}

const updateCricketEndInnings = `
UPDATE cricket_score
SET is_inning_completed=true
WHERE match_id=$1 AND team_id=$2 AND inning=$3
RETURNING *
`

func (q *Queries) UpdateCricketEndInnings(ctx context.Context, matchID, teamID int64, inning string) (models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, updateCricketEndInnings, matchID, teamID, inning)
	var i models.CricketScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.Inning,
		&i.Score,
		&i.Wickets,
		&i.Overs,
		&i.RunRate,
		&i.TargetRunRate,
		&i.FollowOn,
		&i.IsInningCompleted,
		&i.Declared,
	)
	return i, err
}
