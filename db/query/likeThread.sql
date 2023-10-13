-- name: CreateLike :one
INSERT INTO like_thread (
    thread_id,
    username
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetLike :one
SELECT * FROM like_thread
WHERE username = $1 LIMIT $1;

-- name: CountLikeUser :one
SELECT COUNT(*) FROM like_thread
WHERE thread_id = $1;

-- name: UserListLike :many
SELECT * FROM like_thread
ORDER BY username;
