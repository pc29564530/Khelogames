package database

import (
	"context"
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
SELECT * FROM football_score fs
JOIN matches m ON m.id = fs.match_id
JOIN teams t ON t.id = fs.team_id
WHERE m.public_id=$1 AND t.public_id=$2
`

func (q *Queries) GetFootballScore(ctx context.Context, matchPublicID, teamPublicID uuid.UUID) (models.FootballScore, error) {
	row := q.db.QueryRowContext(ctx, getFootballScore, matchPublicID, teamPublicID)
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
	return i, err
}

const newFootballScore = `
WITH resolve_ids AS (
	SELECT m.id AS match_id, t.id AS team_id FROM matches m, teams t
	WHERE m.public_id = $1 AND t.public_id = $2
)
INSERT INTO football_score (
    match_id,
    team_id,
    first_half,
    second_half,
    goals,
	penalty_shootout
)
SELECT ri.match_id, ri.team_id, $3, $4, $5 FROM resolve_ids ri
RETURNING *;
`

type NewFootballScoreParams struct {
	MatchPublicID   uuid.UUID `json:"match_public_id"`
	TeamPublicID    uuid.UUID `json:"team_public_id"`
	FirstHalf       int32     `json:"first_half"`
	SecondHalf      int32     `json:"second_half"`
	Goals           int       `json:"goals"`
	PenaltyShootOut int       `json:"penalty_shootout"`
}

func (q *Queries) NewFootballScore(ctx context.Context, arg NewFootballScoreParams) (models.FootballScore, error) {
	row := q.db.QueryRowContext(ctx, newFootballScore,
		arg.MatchPublicID,
		arg.TeamPublicID,
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
	return i, err
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
	return i, err
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
		&i.MatchID,
		&i.TeamID,
		&i.FirstHalf,
		&i.SecondHalf,
		&i.Goals,
		&i.PenaltyShootOut,
	)
	return i, err
}

const updatePenaltyShootoutScore = `
UPDATE football_score fs
SET 
    penalty_shootout = COALESCE(penalty_shootout, 0) + 1,
FROM matches m, teams t
WHERE 
   m.id = fs.match_id AND t.id = fs.team_id AND m.public_id = $1 AND t.public_id = $2
RETURNING *
`

func (q *Queries) UpdatePenaltyShootoutScore(ctx context.Context, matchPublicID, teamPublicID uuid.UUID) (models.FootballScore, error) {
	row := q.db.QueryRowContext(ctx, updatePenaltyShootoutScore, matchPublicID, teamPublicID)
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
	return i, err
}
