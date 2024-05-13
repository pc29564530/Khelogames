-- name: AddClubMember :one
INSERT INTO "club_member" (
    club_id,
    player_id
) VALUES ($1, $2
) RETURNING *;

-- name: GetClubMember :many
SELECT * FROM "club_member"
WHERE club_id=$1
ORDER BY id ASC;