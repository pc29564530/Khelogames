-- name: CreateComment :one
INSERT INTO comment (
    thread_id,
    owner,
	comment_text,
    created_at
) VALUES (
    $1, $2, $3, CURRENT_TIMESTAMP
)
RETURNING *;

-- name: GetAllComment :many
SELECT * FROM comment
WHERE thread_id=$1;

-- name: GetCommentByUser :many
SELECT * FROM comment
WHERE owner=$1;

-- name: DeleteComment :one
DELETE FROM comment
WHERE id=$1
RETURNING *;