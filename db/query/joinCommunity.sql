-- name: AddJoinCommunity :one
INSERT INTO join_community (
    community_name,
    username
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetUserByCommunity :many
SELECT * FROM join_community
WHERE community_name=$1;

-- name: GetCommunityByUser :many
SELECT * FROM join_community
WHERE username=$1
ORDER BY id;