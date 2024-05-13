-- name: AddFootballMatchScore :one
INSERT INTO football_matches_score (
    match_id,
    tournament_id,
    team_id,
    goal_for,
    goal_against,
    goal_score_time
) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
RETURNING *;

-- name: GetFootballMatchScore :one
SELECT * FROM football_matches_score
WHERE match_id=$1 AND team_id=$2 AND tournament_id=$3;


-- name: UpdateFootballMatchScore :one
UPDATE football_matches_score
SET goal_for=$1, goal_against=$2
WHERE match_id=$3 AND team_id=$4 AND tournament_id=$5
RETURNING *;


-- name: AddFootballGoalByPlayer :one
INSERT INTO football_team_player_score (
    match_id,
    team_id,
    player_id,
    tournament_id,
    goal_score_time
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;