package database

import (
	"context"
	"khelogames/database/models"
)

const getFootballScore = `
SELECT id, match_id, team_id, first_half, second_half, goals FROM football_score
WHERE match_id=$1 AND team_id=$2
`

type GetFootballScoreParams struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (q *Queries) GetFootballScore(ctx context.Context, arg GetFootballScoreParams) (models.FootballScore, error) {
	row := q.db.QueryRowContext(ctx, getFootballScore, arg.MatchID, arg.TeamID)
	var i models.FootballScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.FirstHalf,
		&i.SecondHalf,
		&i.Goals,
	)
	return i, err
}

const newFootballScore = `
INSERT INTO football_score (
    match_id,
    team_id,
    first_half,
    second_half,
    goals
) VALUES ( $1, $2, $3, $4, $5)
RETURNING id, match_id, team_id, first_half, second_half, goals
`

type NewFootballScoreParams struct {
	MatchID    int64 `json:"match_id"`
	TeamID     int64 `json:"team_id"`
	FirstHalf  int32 `json:"first_half"`
	SecondHalf int32 `json:"second_half"`
	Goals      int64 `json:"goals"`
}

func (q *Queries) NewFootballScore(ctx context.Context, arg NewFootballScoreParams) (models.FootballScore, error) {
	row := q.db.QueryRowContext(ctx, newFootballScore,
		arg.MatchID,
		arg.TeamID,
		arg.FirstHalf,
		arg.SecondHalf,
		arg.Goals,
	)
	var i models.FootballScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.FirstHalf,
		&i.SecondHalf,
		&i.Goals,
	)
	return i, err
}

const updateFirstHalfScore = `
UPDATE football_score
SET 
    first_half = COALESCE(first_half, 0) + $1,
    goals = COALESCE(goals, 0) + $1
WHERE 
    match_id = $2 AND team_id = $3
RETURNING id, match_id, team_id, first_half, second_half, goals
`

type UpdateFirstHalfScoreParams struct {
	FirstHalf int32 `json:"first_half"`
	MatchID   int64 `json:"match_id"`
	TeamID    int64 `json:"team_id"`
}

func (q *Queries) UpdateFirstHalfScore(ctx context.Context, arg UpdateFirstHalfScoreParams) (models.FootballScore, error) {
	row := q.db.QueryRowContext(ctx, updateFirstHalfScore, arg.FirstHalf, arg.MatchID, arg.TeamID)
	var i models.FootballScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.FirstHalf,
		&i.SecondHalf,
		&i.Goals,
	)
	return i, err
}

const updateFootballScore = `
UPDATE football_score
SET goals=$1
WHERE match_id=$2 AND team_id=$3
RETURNING id, match_id, team_id, first_half, second_half, goals
`

type UpdateFootballScoreParams struct {
	Goals   int64 `json:"goals"`
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (q *Queries) UpdateFootballScore(ctx context.Context, arg UpdateFootballScoreParams) (models.FootballScore, error) {
	row := q.db.QueryRowContext(ctx, updateFootballScore, arg.Goals, arg.MatchID, arg.TeamID)
	var i models.FootballScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.FirstHalf,
		&i.SecondHalf,
		&i.Goals,
	)
	return i, err
}

const updateSecondHalfScore = `
UPDATE football_score
SET 
    second_half = COALESCE(second_half, 0) + $1,
    goals = COALESCE(goals, 0) + $1
WHERE 
    match_id = $2 AND team_id = $3
RETURNING id, match_id, team_id, first_half, second_half, goals
`

type UpdateSecondHalfScoreParams struct {
	SecondHalf int32 `json:"second_half"`
	MatchID    int64 `json:"match_id"`
	TeamID     int64 `json:"team_id"`
}

func (q *Queries) UpdateSecondHalfScore(ctx context.Context, arg UpdateSecondHalfScoreParams) (models.FootballScore, error) {
	row := q.db.QueryRowContext(ctx, updateSecondHalfScore, arg.SecondHalf, arg.MatchID, arg.TeamID)
	var i models.FootballScore
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.FirstHalf,
		&i.SecondHalf,
		&i.Goals,
	)
	return i, err
}
