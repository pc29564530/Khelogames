-- name: NewCricketScore :one
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
RETURNING *;

-- name: GetCricketScore :one
SELECT * FROM cricket_score
WHERE match_id=$1 AND team_id=$2;

-- name: UpdateCricketScore :one
UPDATE cricket_score
SET score=$1
WHERE match_id=$2 AND team_id=$3
RETURNING *;

-- name: UpdateCricketWickets :one
UPDATE cricket_score
SET wickets=$1
WHERE match_id=$2 AND team_id=$3
RETURNING *;

-- name: UpdateCricketOvers :one
UPDATE cricket_score
SET overs=$1
WHERE match_id=$2 AND team_id=$3
RETURNING *;
