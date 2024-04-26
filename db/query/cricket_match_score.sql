-- name: CreateCricketMatchScore :one
INSERT INTO cricket_match_score (
    match_id,
    tournament_id,
    team_id,
    score,
    wickets,
    overs,
    extras,
    innings
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetCricketMatchScore :one
SELECT * FROM cricket_match_score
WHERE match_id=$1 AND team_id=$2;

-- name: UpdateCricketMatchScore :one
UPDATE cricket_match_score
SET score=$1, wickets=$2, extras=$3, innings=$4
WHERE match_id=$5 AND team_id=$6
RETURNING *;
