-- name: AddFootballSubstitution :one
INSERT INTO football_substitution (
    team_id,
    player_id,
    match_id,
    position
) VALUES ( $1, $2, $3, $4 )
RETURNING *;

-- name: GetFootballSubstitution :many
SELECT * FROM football_substitution
WHERE match_id=$1 AND team_id=$2;
