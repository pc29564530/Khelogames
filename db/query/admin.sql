-- name: AddAdmin :one
INSERT INTO content_admin (
    content_id,
    admin
) VALUES ( $1, $2 )
RETURNING *;

-- name: GetAdmin :many
SELECT * FROM content_admin
WHERE content_id=$1;