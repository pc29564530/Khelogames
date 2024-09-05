-- name: NewMatch :one
INSERT INTO matches (
    tournament_id,
    away_team_id,
    home_team_id,
    start_timestamp,
    end_timestamp,
    type,
    status_code
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetMatch :one
SELECT * FROM matches
WHERE id=$1 AND tournament_id=$2;

-- name: GetMatchByMatchID :one
SELECT * FROM matches
WHERE id=$1;

-- name: GetMatches :many
SELECT * FROM matches
WHERE tournament_id=$1
ORDER BY id DESC;

-- name: UpdateMatchSchedule :one
UPDATE matches
SET start_timestamp=$1
WHERE id=$2
RETURNING *;

-- name: GetMatchesByTournamentID :many
SELECT * FROM matches
WHERE tournament_id=$1
ORDER BY id ASC;

-- name: UpdateMatchStatus :one
UPDATE matches
SET status_code=$1
WHERE id=$2
RETURNING *;

-- name: GetMatchByID :many
SELECT
    m.id, m.tournament_id, m.away_team_id, m.home_team_id, m.start_timestamp, m.end_timestamp, m.type, m.status_code,
    t1.name AS home_team_name, t1.slug AS home_team_slug, t1.shortName AS home_team_shortName, t1.media_url AS home_team_media_url, t1.gender AS home_team_gender, t1.country AS home_team_country, t1.national AS home_team_national, t1.type AS home_team_type,
    t2.name AS away_team_name, t2.slug AS away_team_slug, t2.shortName AS away_team_shortName, t2.media_url AS away_team_media_url, t2.gender AS away_team_gender, t2.country AS away_team_country, t2.national AS away_team_national, t2.type AS away_team_type
FROM matches m
JOIN teams t1 ON m.home_team_id=t1.id
JOIN teams t2 ON m.away_team_id=t2.id
WHERE m.tournament_id=$1;