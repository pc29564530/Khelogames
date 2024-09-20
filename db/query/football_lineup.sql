-- name: AddFootballLineUp :one
INSERT INTO football_lineup (
    team_id,
    player_id,
    match_id,
    position
) VALUES ( $1, $2, $3, $4 )
RETURNING *;

-- name: GetFootballLineUp :many
SELECT * FROM football_lineup
WHERE match_id=$1 AND team_id=$2;

-- name: UpdateFootballSubsAndLineUp :one
WITH 
    sub AS (
        SELECT fs.position, fs.player_id 
        FROM football_substitution fs
        WHERE fs.id = $1
    ),
    lu AS (
        SELECT fl.position, fl.player_id 
        FROM football_lineup fl
        WHERE fl.id = $2
    ),
    update_sub AS (
        UPDATE football_substitution fs
        SET 
            position = lu.position, 
            player_id = lu.player_id
        FROM lu
        WHERE fs.id = $1
        RETURNING fs.*
    ),
    update_lu AS (
        UPDATE football_lineup fl
        SET 
            position = sub.position, 
            player_id = sub.player_id
        FROM sub
        WHERE fl.id = $2
        RETURNING fl.*
    )
SELECT 
    update_sub.*, 
    update_lu.*
FROM 
    update_sub, 
    update_lu;
