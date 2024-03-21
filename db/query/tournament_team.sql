-- name: AddTeam :one
INSERT INTO tournament_team (
    tournament_id,
    team_id
) VALUES ( $1, $2 )
RETURNING *;

-- name: GetTeam :one
SELECT * FROM tournament_team
WHERE team_id=$1;

-- name: GetTeams :many
SELECT c.* FROM tournament_team tt
JOIN club c ON c.id = tt.team_id
WHERE tournament_id=$1;

-- name: GetTeamsCount :one
SELECT COUNT(*) FROM tournament_team
WHERE tournament_id=$1;
