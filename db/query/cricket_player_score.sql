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
    over_number,
    ball_number,
    runs,
    wickets
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: AddCricketWickets :one
INSERT INTO wickets (
    match_id,
    batsman_id,
    bowler_id,
    fielder_id,
    wicket_type
) VALUES ($1, $2, $3, $4, $5)
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
WHERE match_id=$1;


-- name: UpdateCricketRunsScored :one
UPDATE bats
SET runs_scored = runs_scored + $1,
    balls_faced = balls_faced + $2,
    fours = fours + $3,
    sixes = sixes + $4
WHERE match_id = $5 AND batsman_id = $6
RETURNING *;

-- name: UpdateCricketBowler :one
UPDATE balls
SET over_number = over_number + $1,
    ball_number = ball_number + $2,
    runs = runs + $3,
    wickets = wickets + $4
WHERE match_id = $5 AND bowler_id = $6
RETURNING *;
