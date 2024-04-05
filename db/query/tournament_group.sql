-- name: CreateTournamentGroup :one
INSERT INTO group_league (
    group_name,
    tournament_id,
    group_strength
) VALUES ( $1, $2, $3) RETURNING *;

-- name: GetTournamentGroups :many
SELECT * FROM group_league
WHERE tournament_id=$1;

-- name: GetTournamentGroup :one
SELECT * FROM group_league
WHERE (tournament_id=$1 AND group_id=$2);