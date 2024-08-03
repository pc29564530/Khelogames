-- name: NewTournamentTeam :one
INSERT INTO tournament_team (
    tournament_id,
    team_id
) VALUES ( $1, $2 )
RETURNING *;

-- name: GetTournamentTeam :one
SELECT * FROM tournament_team
WHERE team_id=$1;

-- name: GetTournamentTeams :many
SELECT * FROM tournament_team
WHERE tournament_id=$1;

-- name: GetTournamentTeamsCount :one
SELECT COUNT(*) FROM tournament_team
WHERE tournament_id=$1;

