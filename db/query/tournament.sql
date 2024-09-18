-- name: NewTournament :one
INSERT INTO tournaments (
    tournament_name,
    slug,
    sports,
    country,
    status_code,
    level,
    start_timestamp,
    game_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetTournaments :many
SELECT * FROM tournaments;

-- name: GetTournament :one
SELECT * FROM tournaments
WHERE id=$1;

-- name: UpdateTournamentDate :one
UPDATE tournaments
SET start_timestamp=$1
WHERE id=$2
RETURNING *;

-- name: GetTournamentsByLevel :many
SELECT * FROM tournaments
WHERE sports=$1 AND level=$2;

-- name: UpdateTournamentStatus :one
UPDATE tournaments
SET status_code=$1
WHERE id=$2
RETURNING *;

-- name: GetTournamentsBySport :many
SELECT 
    g.id, g.name, g.min_players,
    JSON_BUILD_OBJECT(
        'tournament', JSON_BUILD_OBJECT('id', t.id, 'name', t.name, 'slug', t.slug, 'country', t.country, 'status_code', t.status_code, 'level', t.level, 'start_timestamp', t.start_timestamp, 'game_id', t.game_id)
    ) AS tournament_data
FROM tournaments t
JOIN games AS g ON g.id = t.game_id
WHERE t.game_id=$1;