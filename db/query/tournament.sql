-- name: CreateTournament :one
INSERT INTO "tournament" (
    tournament_name,
    sport_type,
    format,
    teams_joined
) VALUES ($1, $2, $3, $4 )
RETURNING *;

-- name: GetTournaments :many
SELECT * FROM tournament;

-- name: GetTournament :one
SELECT * FROM tournament
WHERE tournament_id=$1;
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 6140c65 (added functionality for tournament teams count)

-- name: UpdateTeamsJoined :one
UPDATE tournament
SET teams_joined=$1
WHERE tournament_id=$2
<<<<<<< HEAD
RETURNING *;
=======
>>>>>>> 5d54c26 (second rebase)
=======
RETURNING *;
>>>>>>> 6140c65 (added functionality for tournament teams count)
