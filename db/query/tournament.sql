-- name: CreateTournament :one
INSERT INTO "tournament" (
    tournament_name,
    sport_type,
    format,
    teams_joined
    start_on
    end_on
) VALUES ($1, $2, $3, $4 )
RETURNING *;

-- name: GetTournaments :many
SELECT * FROM tournament;

-- name: GetTournament :one
SELECT * FROM tournament
WHERE tournament_id=$1;

-- name: UpdateTeamsJoined :one
UPDATE tournament
SET teams_joined=$1
WHERE tournament_id=$2
RETURNING *;

-- name: GetTournamentsBySport :many
SELECT * FROM tournament
WHERE sport_type=$1;

-- name: UpdateTournamentDate :one
UPDATE tournament
SET start_on=$1 OR end_on=$2
WHERE tournament_id=$3
RETURNING *;
