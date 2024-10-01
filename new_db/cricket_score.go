package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const getCricketScore = `
SELECT id, match_id, team_id, inning, score, wickets, overs, run_rate, target_run_rate FROM cricket_score
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
	)
	return i, err
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
    target_run_rate
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, match_id, team_id, inning, score, wickets, overs, run_rate, target_run_rate
`

type NewCricketScoreParams struct {
	MatchID       int64  `json:"match_id"`
	TeamID        int64  `json:"team_id"`
	Inning        string `json:"inning"`
	Score         int32  `json:"score"`
	Wickets       int32  `json:"wickets"`
	Overs         int32  `json:"overs"`
	RunRate       string `json:"run_rate"`
	TargetRunRate string `json:"target_run_rate"`
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
	)
	return i, err
}

const updateCricketInnings = `
UPDATE cricket_score
SET inning=$1
WHERE match_id=$2 AND team_id=$3
RETURNING id, match_id, team_id, inning, score, wickets, overs, run_rate, target_run_rate
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
	)
	return i, err
}

const updateCricketOvers = `
UPDATE cricket_score cs
SET overs = (
        SELECT SUM(bl.balls_faced) FROM bats bl
        WHERE bl.match_id = cs.match_id AND bl.team_id = cs.team_id
        GROUP BY (bl.match_id, bl.team_id)
    )
WHERE cs.match_id=$1 AND cs.team_id=$2
RETURNING id, match_id, team_id, inning, score, wickets, overs, run_rate, target_run_rate
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
	)
	return i, err
}

const updateCricketScore = `
UPDATE cricket_score cs
SET score = (
        SELECT SUM(bt.runs_scored)
        FROM bats bt
        WHERE bt.match_id=cs.match_id AND bt.team_id=cs.team_id
        GROUP BY (bt.match_id, bt.team_id)
    )
WHERE cs.match_id=$1 AND cs.team_id=$2
RETURNING id, match_id, team_id, inning, score, wickets, overs, run_rate, target_run_rate
`

type UpdateCricketScoreParams struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (q *Queries) UpdateCricketScore(ctx context.Context, arg UpdateCricketScoreParams) (models.CricketScore, error) {
	row := q.db.QueryRowContext(ctx, updateCricketScore, arg.MatchID, arg.TeamID)
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
RETURNING id, match_id, team_id, inning, score, wickets, overs, run_rate, target_run_rate
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
	)
	return i, err
}
