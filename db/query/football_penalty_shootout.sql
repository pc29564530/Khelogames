-- name: AddFootballPenaltyShootout :one
INSERT INTO penalties (
    match_id,
    team_id,
    player_id,
    scored
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetFootballPenaltyShootout :many
SELECT * FROM penalties
WHERE match_id=$1 AND team_id=$2
ORDER BY id DESC;

-- name: UpdateFootballPenaltyShootout :one
UPDATE penalties
SET scored = scored + $1
WHERE id=$2
RETURNING *;