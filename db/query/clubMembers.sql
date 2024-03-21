-- name: AddClubMember :one
INSERT INTO "club_member" (
    club_name,
    club_member,
    joined_at
) VALUES ($1, $2, CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetClubMember :many
SELECT * FROM "club_member"
WHERE club_name=$1
ORDER BY id ASC;