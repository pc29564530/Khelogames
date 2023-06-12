-- name: CreateFriends :one
INSERT INTO friends (
    owner,
    friend_username
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetFriends :one
SELECT * FROM friends
WHERE friend_username = $1;

-- name: GetAllFriends :many
SELECT * FROM friends
WHERE owner = $1
ORDER BY id;