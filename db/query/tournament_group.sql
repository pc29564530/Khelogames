-- name: CreateTournamentGroup :one
INSERT INTO groups (
    name,
    tournament_id,
    strength
) VALUES ( $1, $2, $3) RETURNING *;

-- name: GetTournamentGroups :many
SELECT * FROM groups
WHERE tournament_id=$1;

-- name: GetTournamentGroup :one
SELECT * FROM groups
WHERE tournament_id=$1 AND id=$2;

-- name: CreateGroupTeams :one
INSERT INTO teams_group (
    group_id,
    team_id,
    tournament_id
) VALUES ( $1, $2, $3) RETURNING *;

-- name: GetGroupTeams :many
SELECT * FROM teams_group
WHERE tournament_id=$1 AND group_id=$2;