-- name: CreateFriendsRequest :one
INSERT INTO friends_request (
    sender_username,
    reciever_username,
    status
) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFriendsRequest :one
SELECT * FROM friends_request
WHERE id = $1 LIMIT 1;

-- name: ListFriends :many
SELECT * FROM friends_request
ORDER BY id;

-- name: UpdateFriendsRequest :exec
UPDATE friends_request
SET status = 'accepted'
WHERE id = $1
RETURNING *;