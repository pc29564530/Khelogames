-- name: CreateNewMessage :one
INSERT INTO message (
    content,
    is_seen,
    sender_username,
    receiver_username,
    media_url,
    media_type,
    sent_at,
    is_deleted,
    deleted_at
) VALUES (
    $1,$2,$3,$4,$5,$6,CURRENT_TIMESTAMP, $7, $8
) RETURNING *;

-- name: GetMessageByReceiver :many
SELECT * FROM message
WHERE (sender_username=$1 AND receiver_username=$2) OR (receiver_username=$1 AND sender_username=$2)
ORDER BY id ASC;


-- name: GetUserByMessageSend :many
SELECT DISTINCT receiver_username
FROM message
WHERE sender_username = $1;

-- name: UpdateSoftDeleteMessage :one
UPDATE message
SET is_deleted = TRUE, deleted_at = NOW()
WHERE id = $1 AND sender_username = $2
RETURNING *;

-- name: ScheduledDeleteMessage :many
DELETE FROM message
WHERE is_deleted = TRUE AND deleted_at < NOW() - INTERVAL '30 days'
RETURNING *;


-- name: DeleteMessage :one
DELETE FROM message
WHERE sender_username=$1 and id=$2
RETURNING *;


-- name: UpdateDeletedMessage :one
UPDATE message
SET is_deleted=true AND deleted_at=NOW()
WHERE sender_username=$1 and id=$2
RETURNING *;

