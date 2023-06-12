-- name: CreateThread :one
INSERT INTO threads (
    username,
    communities_name,
    title,
    content
) VALUES (
             $1, $2, $3, $4
) RETURNING *;

-- name: GetThread :one
SELECT * FROM threads
WHERE id = $1 LIMIT 1;

-- name: GetAllThreadsByCommunities :many
SELECT * FROM threads
WHERE communities_name = $1
ORDER BY id=$1;

-- name: GetAllThreads :many
SELECT * FROM threads
ORDER BY id=$1;