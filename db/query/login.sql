-- name: login :one
INSERT INTO login (
    username,
    password
) VALUES (
    $1, $2
) RETURNING *;