-- name: CreateSessions :one
INSERT INTO sessions (
    id,
    username,
    refresh_token,
    user_agent,
    client_ip,
    expires_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetSessions :one
SELECT * FROM sessions
WHERE username = $1;

-- name: DeleteSessions :exec
DELETE FROM sessions
WHERE username = $1;