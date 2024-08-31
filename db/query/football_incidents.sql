-- name: CreateFootballIncidents :one
INSERT INTO football_incidents (
    match_id,
    team_id,
    incident_type,
    incident_time,
    description
) VALUES ($1, $2, $3, $4, $5
) RETURNING *;

-- name: AddFootballIncidentPlayer :one
INSERT INTO football_incident_player (
    incident_id,
    player_id
) VALUES ($1, $2)
RETURNING *;

-- name: ADDFootballSubsPlayer :one
INSERT INTO football_substitutions_player (
    incident_id,
    player_in_id,
    player_out_id
) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFootballIncidentPlayer :one
SELECT * FROM football_incident_player
WHERE incident_id=$1;

-- name: GetFootballIncidentSubsPlayer :one
SELECT * FROM football_substitutions_player
WHERE incident_id=$1;

-- name: GetFootballIncidents :many
SELECT * FROM football_incidents
WHERE match_id=$1
ORDER BY created_at DESC;

-- name: GetFootballIncidentBySubstitution :many
SELECT (fi.id, fi.match_id, fi.team_id, fi.incident_type, fi.incident_time, fi.substitution_in_player_id, fi.substitution_out_player_id, fi.description, fi.created_at) FROM football_incidents AS fi
WHERE match_id=$1 AND incident_type="substitutions"
ORDER BY incident_time DESC;

-- name: GetFootballIncidentByGoal :many
SELECT (fi.id, fi.match_id, fi.team_id, fi.incident_type, fi.incident_time, fi.player_id, fi.assist_player_id, fi.description, fi.created_at) FROM football_incidents AS fi
WHERE match_id=$1 AND incident_type="goal"
ORDER BY incident_time DESC;