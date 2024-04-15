-- name: AddFootballMatchScore :one
INSERT INTO football_matches_score (
    match_id,
    tournament_id,
    team_id,
    goal_score,
    goal_score_time
) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
RETURNING *;

-- name: GetFootballMatchScore :one
SELECT * FROM football_matches_score
WHERE match_id=$1 AND team_id=$2;