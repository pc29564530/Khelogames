-- name: CreateFootballIncidents :one
INSERT INTO football_incidents (
    match_id,
    team_id,
    periods,
    incident_type,
    incident_time,
    description,
    penalty_shootout_scored
) VALUES ($1, $2, $3, $4, $5, $6, $7
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

-- name: GetFootballIncidentWithPlayer :many
SELECT
    fi.id, fi.match_id, fi.team_id, fi.periods, fi.incident_type, fi.incident_time, fi.description, fi.penalty_shootout_scored,
    CASE
        WHEN fi.incident_type='substitutions' THEN
            JSON_BUILD_OBJECT(
                'player_in', JSON_BUILD_OBJECT('id',player_in.id,'username',player_in.username, 'name', player_in.player_name, 'slug', player_in.slug, 'short_name',player_in.short_name, 'country', player_in.country, 'positions', player_in.positions, 'media_url', player_in.media_url ),
                'player_out', JSON_BUILD_OBJECT('id',player_out.id,'username',player_out.username, 'name', player_out.player_name, 'slug', player_out.slug, 'short_name',player_out.short_name, 'country', player_out.country, 'positions', player_out.positions, 'media_url', player_out.media_url)
            )
        ELSE
            JSON_BUILD_OBJECT(
                'player', JSON_BUILD_OBJECT('id',player_incident.id,'username',player_incident.username, 'name', player_incident.player_name, 'slug', player_incident.slug, 'short_name',player_incident.short_name, 'country', player_incident.country, 'positions', player_incident.positions, 'media_url', player_incident.media_url)
            )
    END AS players
FROM football_incidents fi
JOIN football_incident_player AS fip ON fip.incident_id=fi.id
JOIN players AS player_incident ON player_incident.id = fip.player_id
LEFT JOIN football_substitutions_player AS fis ON fis.incident_id=fi.id
LEFT JOIN players AS player_in ON player_in.id = fis.player_in_id
LEFT JOIN players AS player_out ON player_out.id = fis.player_out_id
WHERE fi.match_id =  $1
ORDER BY incident_time DESC;

-- name: GetFootballScoreByIncidentTime :many
SELECT SUM ( CASE WHEN team_id=$1 AND incident_type='goal' THEN 1 ELSE 0 END )
FROM football_incidents
WHERE match_id = $2 AND id <= $3;

-- name: GetFootballShootoutScoreByTeam :many
SELECT SUM ( CASE WHEN team_id=$1 AND incident_type='penalty_shootout' AND penalty_shootout_scored='t' THEN 1 ELSE 0 END )
FROM football_incidents
WHERE match_id=$2 AND id <= $3;