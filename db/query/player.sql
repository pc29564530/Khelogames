-- name: NewPlayer :one
INSERT INTO players (
    username,
    slug,
    short_name,
    media_url,
    positions,
    sports,
    country,
    player_name,
    game_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetAllPlayer :many
SELECT * FROM players;

-- name: GetPlayer :one
SELECT * FROM players
WHERE id=$1;

-- name: SearchPlayer :many
SELECT * FROM players
WHERE player_name LIKE $1;

-- name: GetPlayersCountry :many
SELECT * FROM players
WHERE country=$1;

-- name: UpdatePlayerMedia :one
UPDATE players
SET media_url=$1
WHERE id=$2
RETURNING *;

-- name: UpdatePlayerPosition :one
UPDATE players
SET positions=$1
WHERE id=$2
RETURNING *;
