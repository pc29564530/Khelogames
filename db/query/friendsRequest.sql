-- name: CreateConnections :one
INSERT INTO friends_request (
    sender_username,
    reciever_username,
    status
) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetConnections :one
SELECT * FROM friends_request
WHERE sender_username = $1;
--
-- name: ListConnections :many
SELECT * FROM friends_request
WHERE sender_username = $1;

-- name: UpdateConnections :exec
UPDATE friends_request
SET status = 'accepted'
WHERE sender_username = $1
RETURNING *;