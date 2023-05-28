-- name: CreateCommunity :one
INSERT INTO communities (
    communities_name,
    description,
    community_type
) VALUES (
    $1, $2, $3
) RETURNING *;