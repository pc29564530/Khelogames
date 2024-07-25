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

-- name: UpdateFirstHalf :one
UPDATE football_score
SET first_half=$1
WHERE match_id=$2 AND team_id=$3
RETURNING *;

-- name: UpdateSecondHald :one
UPDATE football_score
SET second_half=$1
WHERE match_id=$2 AND team_id=$3
RETURNING *;