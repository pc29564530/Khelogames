-- name: NewTeams :one
INSERT INTO teams (
    name,
    slug,
    shortName,
    admin,
    media_url,
    gender,
    national,
    country,
    type,
    sports
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: AddTeamPlayers :one
INSERT INTO team_players (
    team_id,
    player_id
) VALUES ($1, $2)
RETURNING *;

-- name: GetTeamPlayers :many
SELECT * FROM team_players
WHERE team_id=$1;


-- name: GetTeamsBySport :many
SELECT * FROM teams
WHERE sports=$1;

-- name: GetTeams :many
SELECT * FROM teams;

-- name: GetTeam :one
SELECT * FROM teams
WHERE id=$1;

-- name: UpdateMediaUrl :one
UPDATE teams
SET media_url=$1
WHERE id=$2
RETURNING *;

-- -- name: UpdateTeamsSport :one
-- UPDATE club
-- SET sport=$1
-- WHERE club_name=$2
-- RETURNING *;

-- name: UpdateTeamName :one
UPDATE teams
SET name=$1
WHERE id=$2
RETURNING *;

-- name: SearchTeam :many
SELECT id, name from teams
WHERE name LIKE $1;

-- name: GetTournamentsByTeam :many
SELECT t.id, t.tournament_name, t.sports FROM tournaments t
JOIN tournament_team tt ON t.id=tt.tournament_id
JOIN teams c ON tt.team_id=c.id
WHERE c.id=$1;

-- name: GetMatchByTeam :many
SELECT t.id AS tournament_id, t.tournament_name, tm.id AS match_id, tm.home_team_id, tm.away_team_id, c1.name AS home_team_name, c2.name AS away_team_name, tm.start_timestamp, t.sports
FROM matches tm
JOIN tournaments t ON tm.tournament_id = t.id
JOIN teams c1 ON tm.home_team_id = c1.id
JOIN teams c2 ON tm.away_team_id = c2.id
WHERE c1.id=$1 OR c2.id=$1
ORDER BY tm.id DESC, tm.start_timestamp DESC;
