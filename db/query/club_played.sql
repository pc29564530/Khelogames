-- name: GetClubPlayedTournaments :many
SELECT * FROM club_played
WHERE club_id=$1;

-- name: GetClubPlayedTournament :one
SELECT * FROM club_played
WHERE (club_id=$1 AND tournament_id=$2);

