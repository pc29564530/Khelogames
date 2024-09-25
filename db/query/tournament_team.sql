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
    tt.tournament_id, JSON_BUILD_OBJECT('id', tm.id, 'name', tm.name, 'slug', tm.slug, 'short_name', tm.shortname, 'admin', tm.admin, 'media_url', tm.media_url, 'gender', tm.gender, 'national', tm.national, 'country', tm.country, 'type', tm.type, 'player_count', tm.player_count, 'game_id', tm.game_id) AS team_data
FROM tournament_team tt
JOIN teams AS tm ON tm.id = tt.team_id
WHERE tt.tournament_id=$1;

