-- name: NewMatch :one
INSERT INTO matches (
    tournament_id,
    away_team_id,
    home_team_id,
    start_timestamp,
    end_timestamp,
    type,
    status_code
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetMatch :one
SELECT * FROM matches
WHERE id=$1 AND tournament_id=$2;

-- name: GetMatches :many
SELECT * FROM matches
WHERE tournament_id=$1
ORDER BY id DESC;

-- name: UpdateMatchSchedule :one
UPDATE matches
SET start_timestamp=$1
WHERE id=$2
RETURNING *;

-- name: GetMatchesByTournamentID :many
SELECT * FROM matches
WHERE tournament_id=$1
ORDER BY id ASC;