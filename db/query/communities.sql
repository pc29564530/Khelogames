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
ORDER BY id;


-- name: GetCommunitiesMember :many
SELECT users.username FROM users
JOIN communities ON users.username = communities.owner
WHERE communities.communities_name=$1;