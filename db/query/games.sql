-- name: GetGame :one
SELECT * FROM games
WHERE id=$1;

-- name: GetGames :many
SELECT * FROM games;
