-- name: AddAdmin :one
INSERT INTO content_admin (
    content_id,
    admin
) VALUES ( $1, $2 )
RETURNING *;

-- name: GetAdmin :many
SELECT * FROM content_admin
WHERE content_id=$1;

-- name: UpdateAdmin :one
UPDATE content_admin
SET admin=$1
WHERE content_id=$2 AND admin=$3
RETURNING *;

-- name: DeleteAdmin :one
DELETE FROM content_admin
WHERE content_id=$1 AND admin=$2
RETURNING *;