-- name: CreateTournamentOrganizer :one
INSERT INTO tournament_organizer (
    organizer_id,
    tournament_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetTournamentOrganizer :many
SELECT * FROM tournament_organizer
WHERE tournament_id=$1;