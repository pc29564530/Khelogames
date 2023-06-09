-- name: CreateFriends :one
INSERT INTO friends (
    friend_name,
    friend_username
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetFriends :one
SELECT * FROM friends
WHERE friend_username = $1;

-- name: GetListFriends :many
SELECT * FROM friends
ORDER BY id;