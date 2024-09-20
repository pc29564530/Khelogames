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
    shots_on_target = CASE WHEN $1 IS NOT NULL THEN shots_on_target + $1 ELSE 0 END,
    total_shots = CASE WHEN $2 IS NOT NULL THEN total_shots + $2 ELSE 0 END,
    corner_kicks = CASE WHEN $3 IS NOT NULL THEN corner_kicks + $3 ELSE 0 END,
    fouls = CASE WHEN $4 IS NOT NULL THEN fouls + $4 ELSE fouls END,
    goalkeeper_saves = CASE WHEN $5 IS NOT NULL THEN goalkeeper_saves + $5 ELSE 0 END,
    free_kicks = CASE WHEN $6 IS NOT NULL THEN free_kicks + $6 ELSE 0 END,
    yellow_cards = CASE WHEN $7 IS NOT NULL THEN yellow_cards + $7 ELSE 0 END,
    red_cards = CASE WHEN $8 IS NOT NULL THEN red_cards + $8 ELSE 0 END
WHERE match_id = $9 AND team_id = $10
RETURNING *;
