-- name: CreateFootballIncidents :one
INSERT INTO football_incidents (
    match_id,
    team_id,
    incident_type,
    incident_time,
    player_id,
    description
) VALUES ($1, $2, $3, $4, $5, $6
)
RETURNING *;


-- name: GetFootballIncidents :many
SELECT * FROM football_incidents
WHERE match_id=$1
ORDER BY created_at DESC;