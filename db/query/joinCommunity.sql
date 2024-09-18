-- name: AddJoinCommunity :one
INSERT INTO join_community (
    community_name,
    username
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetUserByCommunity :many
SELECT DISTINCT username FROM join_community
WHERE community_name=$1;

-- name: GetCommunityByUser :many
SELECT * FROM join_community
WHERE username=$1
ORDER BY id;

-- name: RemoveUserFromCommunity :one
DELETE FROM join_community
WHERE id=$1 AND username=$2
RETURNING *;


-- name: InActiveUserFromCommunity :one
UPDATE join_community
SET is_active = FALSE
WHERE username = $1 AND id = $2
RETURNING *;
