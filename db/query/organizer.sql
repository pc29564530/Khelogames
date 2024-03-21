-- name: CreateOrganizer :one
INSERT INTO organizer (
    organizer_name,
    tournament_id
) VALUES ($1, $2) RETURNING *;

-- name: GetOrganizer :many
SELECT * FROM organizer
WHERE tournament_id=$1;
