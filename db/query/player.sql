-- name: NewPlayer :one
INSERT INTO players (
    name,
    slug,
    short_name,
    media_url,
    positions,
    sports,
    country
) VALUES ($1, $2, $3, $4, $5, $6, $7 ) RETURNING *;

-- name: GetPlayer :one
SELECT * FROM players
WHERE id=$1;

-- name: SearchPlayer :many
SELECT * FROM players
WHERE name LIKE $1;

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
