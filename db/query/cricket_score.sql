-- name: NewCricketScore :one
INSERT INTO cricket_score (
    match_id,
    team_id,
    inning,
    score,
    wickets,
    overs,
    run_rate,
    target_run_rate
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetCricketScore :one
SELECT * FROM cricket_score
WHERE match_id=$1 AND team_id=$2;

-- name: UpdateCricketScore :one
UPDATE cricket_score cs
SET score = (
        SELECT SUM(bt.runs_scored)
        FROM bats bt
        WHERE bt.match_id=cs.match_id AND bt.team_id=cs.team_id
        GROUP BY (bt.match_id, bt.team_id)
    )
WHERE cs.match_id=$1 AND cs.team_id=$2
RETURNING *;

-- name: UpdateCricketWickets :one
UPDATE cricket_score cs
SET wickets = (
        SELECT COUNT(*)
        FROM wickets w
        WHERE w.match_id = cs.match_id AND w.team_id = cs.team_id
    )
WHERE cs.match_id = $1 AND cs.team_id = $2
RETURNING *;

-- name: UpdateCricketOvers :one
UPDATE cricket_score cs
SET overs = (
        SELECT SUM(bl.balls_faced) FROM bats bl
        WHERE bl.match_id = cs.match_id AND bl.team_id = cs.team_id
        GROUP BY (bl.match_id, bl.team_id)
    )
WHERE cs.match_id=$1 AND cs.team_id=$2
RETURNING *;

-- name: UpdateCricketInnings :one
UPDATE cricket_score
SET inning=$1
WHERE match_id=$2 AND team_id=$3
RETURNING *;
