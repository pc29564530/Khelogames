-- name: CreateNewMessage :one
INSERT INTO message (
    content,
    is_seen,
    sender_username,
    receiver_username
) VALUES (
    $1,$2,$3,$4
) RETURNING *;

-- name: GetMessageByReceiver :many
SELECT * FROM message
WHERE (sender_username=$1 AND receiver_username=$2) OR (receiver_username=$1 AND sender_username=$2)
ORDER BY id ASC;