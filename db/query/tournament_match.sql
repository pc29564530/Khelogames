-- name: CreateMatch :one
INSERT INTO tournament_match (
    organizer_id,
    tournament_id,
    team1_id,
    team2_id,
    date_on,
    start_time,
    stage,
    sports,
    end_time
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,$8, $9
) RETURNING *;

-- name: GetMatch :one
SELECT * FROM tournament_match
WHERE match_id=$1 AND tournament_id=$2;

-- name: GetTournamentMatch :many
SELECT * FROM tournament_match
WHERE (tournament_id=$1 AND sports=$2)
ORDER BY match_id ASC;

-- name: UpdateMatchSchedule :one
UPDATE tournament_match
SET date_on=$1
WHERE match_id=$2
RETURNING *;


-- name: UpdateMatchScheduleTime :one
UPDATE tournament_match
SET start_time=$1 OR end_time=$2
WHERE match_id=$3
RETURNING *;

-- name: GetMatchesByTournamentID :many
SELECT * FROM tournament_match
WHERE tournament_id=$1
ORDER BY match_id ASC;


