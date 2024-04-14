-- name: CreateClub :one
INSERT INTO "club" (
    club_name,
    avatar_url,
    sport,
    owner,
    created_at
) VALUES (
    $1, $2, $3, $4, CURRENT_TIMESTAMP
) RETURNING *;

-- name: GetClubsBySport :many
SELECT * FROM "club"
WHERE sport=$1;

-- name: GetClubs :many
SELECT * FROM "club";

-- name: GetClub :one
SELECT * FROM "club"
WHERE id=$1;

-- name: UpdateAvatarUrl :one
UPDATE club
SET avatar_url=$1
WHERE club_name=$2
RETURNING *;

-- name: UpdateClubSport :one
UPDATE club
SET sport=$1
WHERE club_name=$2
RETURNING *;

-- name: UpdateClubName :one
UPDATE club
SET club_name=$1
WHERE club_name=$2
RETURNING *;

-- name: SearchTeam :many
SELECT id, club_name from club
WHERE club_name LIKE $1;

-- name: GetTournamentsByClub :many
SELECT t.tournament_id, t.tournament_name, t.format, t.sport_type , t.tournament_id from tournament t
JOIN tournament_team tt ON t.tournament_id=tt.tournament_id
JOIN club c ON tt.team_id=c.id
WHERE c.club_name=$1;

-- name: GetMatchByClubName :many
SELECT t.tournament_id, t.tournament_name, c1.club_name AS team1_name, c2.club_name AS team2_name, tm.start_at, tm.date_on
FROM tournament_match tm
JOIN tournament t ON tm.tournament_id = t.tournament_id
JOIN club c1 ON tm.team1_id = c1.id
JOIN club c2 ON tm.team2_id = c2.id
WHERE c1.club_name=$1 OR c2.club_name=$1
ORDER BY tm.match_id AND tm.start_at DESC;
