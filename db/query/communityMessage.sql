-- name: CreateCommunityMessage :one
INSERT INTO communitymessage(
    community_name,
    sender_username,
    content,
    sent_at
) VALUES ($1,$2, $3, CURRENT_TIMESTAMP )
RETURNING *;

-- name: CreateUploadMedia :one
INSERT INTO uploadmedia (
    media_url,
    media_type,
    sent_at
) VALUES ($1, $2, CURRENT_TIMESTAMP)
RETURNING id, media_url, media_type, sent_at;

-- name: CreateMessageMedia :one
INSERT INTO messagemedia (
    message_id,
    media_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetCommuntiyMessage :many
SELECT cm.*, um.media_url, um.media_type 
FROM communitymessage cm 
JOIN messagemedia mm ON mm.message_id = cm.id 
JOIN uploadmedia um ON mm.media_id = um.id;

-- name: GetCommunityByMessage :many
SELECT DISTINCT community_name FROM communitymessage;

