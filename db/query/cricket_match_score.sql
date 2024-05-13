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

-- name: UpdateCricketMatchRunsScore :one
UPDATE cricket_match_score
SET score=$1
WHERE match_id=$2 AND team_id=$3
RETURNING *;

-- name: UpdateCricketMatchWickets :one
UPDATE cricket_match_score
SET wickets=$1
WHERE match_id=$2 AND team_id=$3
RETURNING *;

-- name: UpdateCricketMatchExtras :one
UPDATE cricket_match_score
SET extras=$1
WHERE match_id=$2 AND team_id=$3
RETURNING *;

-- name: UpdateCricketMatchInnings :one
UPDATE cricket_match_score
SET innings=$1
WHERE match_id=$2 AND team_id=$3
RETURNING *;
