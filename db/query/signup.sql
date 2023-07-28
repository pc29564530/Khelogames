-- name: CreateSignup :one
INSERT INTO signup (
    mobile_number,
    otp
) VALUES (
    $1, $2
) RETURNING *;


-- name: GetSignup :one
SELECT * FROM signup
WHERE mobile_number = $1 LIMIT 1;

-- name: DeleteSignup :one
DELETE FROM signup
WHERE mobile_number = $1 RETURNING *;

