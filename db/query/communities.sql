-- name: CreateCommunity :one
INSERT INTO communities (
    communities_name,
    description,
    community_type
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetCommunity :one
SELECT id,communities_name,description,community_type,created_at FROM communities
WHERE id = $1 LIMIT 1;

-- name: GetListCommunity :many
SELECT * FROM communities
ORDER BY id;