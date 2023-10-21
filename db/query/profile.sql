-- name: CreateProfile :one
INSERT INTO profile (
    owner,
    full_name,
    bio,
    following_owner,
    follower_owner,
    avatar_url
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetProfile :one
SELECT * FROM profile
WHERE owner=$1;

-- name: UpdateProfileAvatar :one
UPDATE profile
SET avatar_url=$1
WHERE id=$2
RETURNING *;

-- name: UpdateProfileFullName :one
UPDATE profile
SET full_name=$1
WHERE id=$2
RETURNING *;

-- name: UpdateProfileBio :one
UPDATE profile
SET bio=$1
WHERE id=$2
RETURNING *;