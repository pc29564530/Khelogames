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