-- name: AddGroupTeam :one
INSERT INTO "group_team" (
    group_id,
    tournament_id,
    team_id
) VALUES ($1, $2, $3) RETURNING *;

-- name: GetGroupTeam :many
SELECT * FROM "group_team"
WHERE tournament_id=$1;

-- name: GetTeamByGroup :many
SELECT * FROM "group_team"
WHERE (tournament_id=$1 AND group_id=$2);