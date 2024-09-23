-- name: AddFootballPenaltyShootout :one
INSERT INTO penalty_shootout (
    match_id,
    team_id,
    player_id,
    scored
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetFootballPenaltyShootout :many
SELECT * FROM penalty_shootout
WHERE match_id=$1 AND team_id=$2
ORDER BY id DESC;

-- name: UpdateFootballPenaltyShootout :one
UPDATE penalty_shootout
SET scored = scored + $1
WHERE id=$2
RETURNING *;