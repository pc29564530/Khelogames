-- name: AddCricketMatchToss :one
INSERT INTO cricket_match_team_toss (
    tournament_id,
    match_id,
    toss_won,
    bat_or_bowl
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCricketMatchToss :one
SELECT * FROM cricket_match_team_toss
WHERE tournament_id=$1 AND match_id=$2;