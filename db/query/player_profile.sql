-- name: AddPlayerProfile :one
INSERT INTO player_profile (
    player_name,
    player_avatar_url,
    player_bio,
    player_sport,
    player_playing_category,
    nation
) VALUES ( $1, $2, $3, $4, $5, $6 )
RETURNING *;

-- name: GetPlayerProfile :one
SELECT * FROM player_profile
WHERE id=$1;

-- name: UpdatePlayerProfileAvatar :one
UPDATE player_profile
SET player_avatar_url=$1
WHERE id=$2
RETURNING *;