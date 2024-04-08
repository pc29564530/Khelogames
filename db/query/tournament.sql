-- name: CreateTournament :one
INSERT INTO "tournament" (
    tournament_name,
    sport_type,
    format,
    teams_joined
) VALUES ($1, $2, $3, $4 )
RETURNING *;

-- name: GetTournaments :many
SELECT * FROM tournament;

-- name: GetTournament :one
SELECT * FROM tournament
WHERE tournament_id=$1;
