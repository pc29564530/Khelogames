-- name: CreateFootballStatistics :one
INSERT INTO football_statistics (
    match_id,
    team_id,
    shots_on_target,
    total_shots,
    corner_kicks,
    fouls,
    goalkeeper_saves,
    free_kicks,
    yellow_cards,
    red_cards
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetFootballStatistics :one
SELECT * FROM football_statistics
WHERE match_id=$1 AND team_id=$2;

-- name: UpdateFootballStatistics :one
UPDATE football_statistics
SET 
    shots_on_target = shots_on_target + $1,
    total_shots = total_shots + $2,
    corner_kicks = corner_kicks + $3,
    fouls = fouls + $4,
    goalkeeper_saves = goalkeeper_saves + $5,
    free_kicks = free_kicks + $6,
    yellow_cards = yellow_cards + $7,
    red_cards = red_cards + $8
WHERE match_id = $9 AND team_id = $10
RETURNING *;
