-- name: AddCricketBatsScore :one
INSERT INTO bats (
    batsman_id,
    match_id,
    team_id,
    position,
    runs_scored,
    balls_faced,
    fours,
    sixes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;


-- name: AddCricketBall :one
INSERT INTO balls (
    match_id,
    team_id,
    bowler_id,
    ball,
    runs,
    wickets,
    wide,
    no_ball
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: AddCricketWickets :one
INSERT INTO wickets (
    match_id,
    team_id,
    batsman_id,
    bowler_id,
    wickets_number,
    wicket_type
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCricketPlayerScore :one
SELECT * FROM bats
WHERE match_id=$1 AND batsman_id=$2 LIMIT 1;


-- name: GetCricketPlayersScore :many
SELECT * FROM bats
WHERE match_id=$1 AND team_id=$2
ORDER BY position;

-- name: GetCricketBall :one
SELECT * FROM balls
WHERE match_id=$1 AND bowler_id=$2 LIMIT 1;


-- name: GetCricketBalls :many
SELECT * FROM balls
WHERE match_id=$1 AND team_id=$2;

-- name: GetCricketWicket :one
SELECT * FROM wickets
WHERE match_id=$1 AND batsman_id=$2 LIMIT 1;

-- name: GetCricketWickets :many
SELECT * FROM wickets
WHERE match_id=$1 AND team_id=$2;


-- name: UpdateCricketRunsScored :one
UPDATE bats
SET runs_scored = runs_scored + $1,
    balls_faced = balls_faced + $2,
    fours = fours + $3,
    sixes = sixes + $4
WHERE match_id = $5 AND batsman_id = $6 AND team_id=$7
RETURNING *;

-- name: UpdateCricketBowler :one
UPDATE balls
SET 
    ball = $1,
    runs = runs + $2,
    wickets = wickets + $3,
    wide = wide + $4,
    no_ball = no_ball + $5
WHERE match_id = $6 AND bowler_id = $7 AND team_id=$8
RETURNING *;
