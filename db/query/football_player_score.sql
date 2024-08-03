-- name: AddFootballPlayerScore :one
INSERT INTO goals (
    match_id,
    team_id,
    player_id,
    goal_time
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetFootballPlayerScore :many
SELECT * FROM goals
WHERE match_id=$1 AND player_id=$2;

-- name: CountGoalByPlayerTeam :one
SELECT COUNT(*) FROM goals
WHERE team_id=$1 AND  player_id=$2;