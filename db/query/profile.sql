-- name: CreateProfile :one
INSERT INTO profile (
    owner,
    full_name,
    bio,
    avatar_url,
    cover_url,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, CURRENT_TIMESTAMP
) RETURNING *;

-- name: EditProfile :one
UPDATE profile
SET full_name=$1, avatar_url=$2, bio=$3, cover_url=$4
WHERE id=$5
RETURNING *;


-- name: GetProfile :one
SELECT * FROM profile
WHERE owner=$1;

-- name: UpdateAvatar :one
UPDATE profile
SET avatar_url=$1
WHERE owner=$2
RETURNING *;

-- name: UpdateCover :one
UPDATE profile
SET cover_url=$1
WHERE owner=$2
RETURNING *;

-- name: UpdateFullName :one
UPDATE profile
SET full_name=$1
WHERE owner=$2
RETURNING *;

-- name: UpdateBio :one
UPDATE profile
SET bio=$1
WHERE owner=$2
RETURNING *;