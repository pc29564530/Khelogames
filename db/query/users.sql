-- name: CreateUser :one
INSERT INTO users (
  username,
  mobile_number
) VALUES (
  $1, $2
) RETURNING *;
