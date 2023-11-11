-- name: CreateProfile :one
INSERT INTO profile (
    owner,
    full_name,
    bio,
    avatar_url,
    cover_url
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: EditProfile :one
UPDATE profile
SET full_name=$1, avatar_url=$2, bio=$3, cover_url=$4
WHERE id=$5
RETURNING *;


-- name: GetProfile :one
SELECT * FROM profile
WHERE owner=$1;