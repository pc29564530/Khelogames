-- name: CreateTournamentStanding :one
INSERT INTO tournament_standing (
    tournament_id,
    group_id,
    team_id,
    wins,
    loss,
    draw,
    goal_for,
    goal_against,
    goal_difference,
    points
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10 ) RETURNING *;

-- name: GetTournamentStanding :many
SELECT 
    ts.standing_id, ts.tournament_id, ts.group_id, ts.team_id,
    ts.wins, ts.loss, ts.draw, ts.goal_for, ts.goal_against, ts.goal_difference, ts.points,
    t.tournament_name, t.sport_type,
    c.club_name, t.format
FROM 
    tournament_standing ts
JOIN 
    group_league tg ON ts.group_id = tg.group_id
JOIN 
    tournament t ON ts.tournament_id = t.tournament_id
JOIN 
    club c ON ts.team_id = c.id
WHERE 
    ts.tournament_id = $1
    AND tg.group_id = $2
    AND t.sport_type = $3;