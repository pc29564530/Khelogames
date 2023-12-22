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
WHERE (sender_username=$1 AND receiver_username=$2 AND sent_at=%3)
ORDER BY id DESC;