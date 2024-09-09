-- name: NewFootballScore :one
INSERT INTO football_score (
    match_id,
    team_id,
    first_half,
    second_half,
    goals
) VALUES ( $1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFootballScore :one
SELECT * FROM football_score
WHERE match_id=$1 AND team_id=$2;

-- name: UpdateFootballScore :one
UPDATE football_score
SET goals=$1
WHERE match_id=$2 AND team_id=$3
RETURNING *;

-- name: UpdateFirstHalfScore :one
UPDATE football_score
SET 
    first_half = COALESCE(first_half, 0) + $1,
    goals = COALESCE(goals, 0) + $1
WHERE 
    match_id = $2 AND team_id = $3
RETURNING *;

-- name: UpdateSecondHalfScore :one
UPDATE football_score
SET 
    second_half = COALESCE(second_half, 0) + $1,
    goals = COALESCE(goals, 0) + $1
WHERE 
    match_id = $2 AND team_id = $3
RETURNING *;