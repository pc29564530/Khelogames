-- name: CreateClub :one
INSERT INTO "club" (
    club_name,
    avatar_url,
    sport,
    owner,
    created_at
) VALUES (
    $1, $2, $3, $4, CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetClubs :many
SELECT * FROM "club";

-- name: GetClub :one
SELECT * FROM "club"
WHERE id=$1;

-- name: UpdateAvatarUrl :one
UPDATE club
SET avatar_url=$1
WHERE club_name=$2
RETURNING *;

-- name: UpdateClubSport :one
UPDATE club
SET sport=$1
WHERE club_name=$2
RETURNING *;

-- name: UpdateClubName :one
UPDATE club
SET club_name=$1
WHERE club_name=$2
RETURNING *;

-- name: SearchTeam :many
SELECT id, club_name from club
WHERE club_name LIKE $1;