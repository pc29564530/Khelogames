package database

import (
	"context"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const addCricketToss = `
WITH matchID AS (
	SELECT * FROM matches WHERE public_id=$1
),
teamID AS (
	SELECT * FROM teams WHERE public_id = $3
)
INSERT INTO cricket_toss (
    match_id,
    toss_decision,
    toss_win
) 
SELECT
	matchID.id,
	$2,
	teamID.id
FROM matchID, teamID
RETURNING *
`

func (q *Queries) AddCricketToss(ctx context.Context, matchPublicID uuid.UUID, tossDecision string, tossWinPublicID uuid.UUID) (models.CricketToss, error) {
	row := q.db.QueryRowContext(ctx, addCricketToss, matchPublicID, tossDecision, tossWinPublicID)
	var i models.CricketToss
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TossDecision,
		&i.TossWin,
	)
	return i, err
}

const getCricketToss = `
SELECT * FROM cricket_toss ct
JOIN matches m ON m.id = ct.match_id
WHERE m.public_id=$1
`

func (q *Queries) GetCricketToss(ctx context.Context, matchPublicID uuid.UUID) (models.CricketToss, error) {
	row := q.db.QueryRowContext(ctx, getCricketToss, matchPublicID)
	var i models.CricketToss
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TossDecision,
		&i.TossWin,
	)
	return i, err
}
