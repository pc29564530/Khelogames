-- name: SearchByCommunityName :one
SELECT communities_name FROM communities
WHERE communities_name ILIKE $1;

-- name: SearchCommunityByCommunityType :many
SELECT communities_name FROM communities
WHERE community_type ILIKE $1
ORDER BY id;

-- name: SearchByFullName :many
SELECT full_name, owner FROM profile
WHERE full_name ILIKE $1
ORDER BY id;