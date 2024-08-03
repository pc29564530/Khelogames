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
    t.tournament_name, t.sports,
    c.name
FROM 
    tournament_standing ts
JOIN 
    group_league tg ON ts.group_id = tg.group_id
JOIN 
    tournaments t ON ts.tournament_id = t.id
JOIN 
    teams c ON ts.team_id = c.id
WHERE 
    ts.tournament_id = $1
    AND tg.group_id = $2
    AND t.sports = $3;



-- name: UpdateTournamentStanding :one
UPDATE tournament_standing AS ts
SET 
    goal_for = COALESCE((
        SELECT SUM(CASE WHEN fs.goal_for IS NOT NULL THEN fs.goal_for ELSE 0 END)
        FROM football_matches_score AS fs
        WHERE fs.team_id = ts.team_id
    ), 0),
    goal_against = COALESCE((
        SELECT SUM(CASE WHEN fs.goal_against IS NOT NULL THEN fs.goal_against ELSE 0 END)
        FROM football_matches_score AS fs
        WHERE fs.team_id = ts.team_id
    ), 0),
    goal_difference = COALESCE((
        SELECT SUM(CASE WHEN fs.goal_for IS NOT NULL THEN fs.goal_for ELSE 0 END) -
               SUM(CASE WHEN fs.goal_against IS NOT NULL THEN fs.goal_against ELSE 0 END)
        FROM football_matches_score AS fs
        WHERE fs.team_id = ts.team_id
    ), 0),
    wins = COALESCE((
        SELECT COUNT(*)
        FROM football_matches_score AS fs
        WHERE fs.team_id = ts.team_id AND fs.goal_for > fs.goal_against
    ), 0),
    loss = COALESCE((
        SELECT COUNT(*)
        FROM football_matches_score AS fs
        WHERE fs.team_id = ts.team_id AND fs.goal_for < fs.goal_against
    ), 0),
    draw = COALESCE((
        SELECT COUNT(*)
        FROM football_matches_score AS fs
        WHERE fs.team_id = ts.team_id AND fs.goal_for = fs.goal_against
    ), 0),
    points = ((wins*3)+draw)
WHERE ts.tournament_id = $1 AND ts.team_id=$2
RETURNING *;
