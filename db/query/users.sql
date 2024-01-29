-- name: CreateUser :one
INSERT INTO users (
  username,
  mobile_number,
  hashed_password,
  created_at
) VALUES (
  $1, $2, $3, CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUser :many
SELECT DISTINCT * FROM users
WHERE username = $1
ORDER BY username
LIMIT $2
OFFSET $3;