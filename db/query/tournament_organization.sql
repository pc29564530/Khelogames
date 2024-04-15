-- name: CreateTournamentOrganization :one
INSERT INTO tournament_organization (
    tournament_id,
    tournament_start,
    player_count,
    team_count,
    group_count,
    advanced_team
) VALUES 
($1, $2, $3, $4, $5, $6 )
RETURNING *;

-- name: GetTournamentOrganization :one
SELECT * FROM tournament_organization
WHERE tournament_id=$1;