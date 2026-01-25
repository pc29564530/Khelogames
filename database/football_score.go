package database

import (
	"context"
	"database/sql"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const GetFootballScoreByMatchID = `
SELECT * FROM football_score fs
JOIN matches m ON fs.match_id = m.id
WHERE m.public_id=$1
`

func (q *Queries) GetFootballScoreByMatchID(ctx context.Context, matchPublicID uuid.UUID) ([]models.FootballScore, error) {
	rows, err := q.db.QueryContext(ctx, GetFootballScoreByMatchID, matchPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.FootballScore
	for rows.Next() {
		var i models.FootballScore
		err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.MatchID,
			&i.TeamID,
			&i.FirstHalf,
			&i.SecondHalf,
			&i.Goals,
			&i.PenaltyShootOut,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}
	return items, nil
}

const getFootballScore = `
SELECT * FROM football_score
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
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.FirstHalf,
		&i.SecondHalf,
		&i.Goals,
		&i.PenaltyShootOut,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.FootballScore{}, nil
		}
		return models.FootballScore{}, err
	}
	return i, nil
}

const newFootballScore = `
INSERT INTO football_score (
    match_id,
    team_id,
    first_half,
    second_half,
    goals,
	penalty_shootout
) VALUES ( $1, $2, $3, $4, $5, $6)
RETURNING *
`

type NewFootballScoreParams struct {
	MatchID         int32 `json:"match_id"`
	TeamID          int32 `json:"team_id"`
	FirstHalf       int   `json:"first_half"`
	SecondHalf      int   `json:"second_half"`
	Goals           int   `json:"goals"`
	PenaltyShootOut *int  `json:"penalty_shootout"`
}

func (q *Queries) NewFootballScore(ctx context.Context, arg NewFootballScoreParams) (models.FootballScore, error) {
	row := q.db.QueryRowContext(ctx, newFootballScore,
		arg.MatchID,
		arg.TeamID,
		arg.FirstHalf,
		arg.SecondHalf,
		arg.Goals,
		arg.PenaltyShootOut,
	)
	var i models.FootballScore
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.FirstHalf,
		&i.SecondHalf,
		&i.Goals,
		&i.PenaltyShootOut,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.FootballScore{}, nil
		}
		return models.FootballScore{}, err
	}
	return i, nil
}

const updateFirstHalfScore = `
UPDATE football_score
SET 
    first_half = COALESCE(first_half, 0) + $1,
    goals = COALESCE(goals, 0) + $1
WHERE 
    match_id = $2 AND team_id = $3
RETURNING *
`

type UpdateFirstHalfScoreParams struct {
	FirstHalf int   `json:"first_half"`
	MatchID   int32 `json:"match_id"`
	TeamID    int32 `json:"team_id"`
}

func (q *Queries) UpdateFirstHalfScore(ctx context.Context, arg UpdateFirstHalfScoreParams) (models.FootballScore, error) {
	row := q.db.QueryRowContext(ctx, updateFirstHalfScore, arg.FirstHalf, arg.MatchID, arg.TeamID)
	var i models.FootballScore
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.FirstHalf,
		&i.SecondHalf,
		&i.Goals,
		&i.PenaltyShootOut,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.FootballScore{}, nil
		}
		return models.FootballScore{}, err
	}
	return i, nil
}

const updateSecondHalfScore = `
UPDATE football_score
SET 
    second_half = COALESCE(second_half, 0) + $1,
    goals = COALESCE(goals, 0) + $1
WHERE 
    match_id = $2 AND team_id = $3
RETURNING *
`

type UpdateSecondHalfScoreParams struct {
	SecondHalf int32 `json:"second_half"`
	MatchID    int32 `json:"match_id"`
	TeamID     int32 `json:"team_id"`
}

func (q *Queries) UpdateSecondHalfScore(ctx context.Context, arg UpdateSecondHalfScoreParams) (models.FootballScore, error) {
	row := q.db.QueryRowContext(ctx, updateSecondHalfScore, arg.SecondHalf, arg.MatchID, arg.TeamID)
	var i models.FootballScore
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.FirstHalf,
		&i.SecondHalf,
		&i.Goals,
		&i.PenaltyShootOut,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.FootballScore{}, nil
		}
		return models.FootballScore{}, err
	}
	return i, nil
}

const updatePenaltyShootoutScore = `
UPDATE football_score
SET 
    penalty_shootout = COALESCE(penalty_shootout, 0) + 1,
WHERE 
    match_id = $1 AND team_id = $2
RETURNING *
`

func (q *Queries) UpdatePenaltyShootoutScore(ctx context.Context, matchID, teamID int32) (models.FootballScore, error) {
	row := q.db.QueryRowContext(ctx, updatePenaltyShootoutScore, matchID, teamID)
	var i models.FootballScore
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.FirstHalf,
		&i.SecondHalf,
		&i.Goals,
		&i.PenaltyShootOut,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.FootballScore{}, nil
		}
		return models.FootballScore{}, err
	}
	return i, nil
}
