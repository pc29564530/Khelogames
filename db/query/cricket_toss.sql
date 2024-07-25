-- name: AddCricketToss :one
INSERT INTO cricket_toss (
    match_id,
    toss_decision,
    toss_win
) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCricketToss :one
SELECT * FROM cricket_toss
WHERE match_id=$1;