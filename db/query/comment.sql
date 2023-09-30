-- name: CreateComment :one
INSERT INTO comment (
    thread_id,
    owner,
	comment_text
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetAllComment :many
SELECT * FROM comment
WHERE thread_id=$1;