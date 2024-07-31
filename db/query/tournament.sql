-- name: NewTournament :one
INSERT INTO tournaments (
    tournament_name,
    slug,
    sports,
    country,
    status_code,
    level,
    start_timestamp
    
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetTournaments :many
SELECT * FROM tournaments;

-- name: GetTournament :one
SELECT * FROM tournaments
WHERE id=$1;

-- name: GetTournamentsBySport :many
SELECT * FROM tournaments
WHERE sports=$1;

-- name: UpdateTournamentDate :one
UPDATE tournaments
SET start_timestamp=$1
WHERE id=$2
RETURNING *;

-- name: GetTournamentsByLevel :many
SELECT * FROM tournaments
WHERE sports=$1 AND level=$2;

-- name: UpdateTournamentStatus :one
UPDATE tournaments
SET status_code=$1
WHERE id=$2
RETURNING *;