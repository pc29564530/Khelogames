package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const addCricketToss = `
INSERT INTO cricket_toss (
    match_id,
    toss_decision,
    toss_win
) VALUES ($1, $2, $3)
RETURNING id, match_id, toss_decision, toss_win
`

type AddCricketTossParams struct {
	MatchID      int64  `json:"match_id"`
	TossDecision string `json:"toss_decision"`
	TossWin      int64  `json:"toss_win"`
}

func (q *Queries) AddCricketToss(ctx context.Context, arg AddCricketTossParams) (models.CricketToss, error) {
	row := q.db.QueryRowContext(ctx, addCricketToss, arg.MatchID, arg.TossDecision, arg.TossWin)
	var i models.CricketToss
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TossDecision,
		&i.TossWin,
	)
	return i, err
}

const getCricketToss = `
SELECT id, match_id, toss_decision, toss_win FROM cricket_toss
WHERE match_id=$1
`

func (q *Queries) GetCricketToss(ctx context.Context, matchID int64) (models.CricketToss, error) {
	row := q.db.QueryRowContext(ctx, getCricketToss, matchID)
	var i models.CricketToss
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TossDecision,
		&i.TossWin,
	)
	return i, err
}
