-- name: NewTournamentTeam :one
INSERT INTO tournament_team (
    tournament_id,
    team_id
) VALUES ( $1, $2 )
RETURNING *;

-- name: GetTournamentTeam :one
SELECT * FROM tournament_team
WHERE team_id=$1;

-- name: GetTournamentTeamsCount :one
SELECT COUNT(*) FROM tournament_team
WHERE tournament_id=$1;

-- name: GetTournamentTeams :many
SELECT 
    tm.id,
    tm.name,
    tm.admin,
    tm.slug,
    tm.shortname,
    tm.country,
    tm.media_url,
    tm.type,
    tm.gender,
    tm.national,
    tm.sports 
FROM tournament_team tt
LEFT JOIN teams AS tm ON tm.id == tt.team_id
WHERE tournament_id=$1;

