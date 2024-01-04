-- name: CreateFollowing :one
INSERT INTO follow (
    follower_owner,
    following_owner
) VALUES (
             $1, $2
) RETURNING *;

-- name: GetAllFollower :many
SELECT DISTINCT follower_owner FROM follow
WHERE following_owner = $1;

-- name: GetAllFollowing :many
SELECT DISTINCT following_owner FROM follow
WHERE follower_owner =  $1;

-- name: DeleteFollowing :one
DELETE FROM follow
WHERE following_owner = $1 RETURNING *;
