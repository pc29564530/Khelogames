-- name: CreateCommunity :one
INSERT INTO communities (
    owner,
    communities_name,
    description,
    community_type
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetCommunity :one
SELECT * FROM communities
WHERE id = $1 LIMIT 1;

-- name: GetAllCommunities :many
SELECT * FROM communities
WHERE owner=$1;
