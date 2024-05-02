-- name: AddCricketTeamPlayerScore :one
INSERT INTO cricket_team_player_score (
    match_id,
    tournament_id,
    team_id,
    batting_or_bowling,
    position,
    player_id,
    runs_scored,
    balls_faced,
    fours,
    sixes,
    wickets_taken,
    overs_bowled,
    runs_conceded,
    wicket_taken_by,
    wicket_of
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15 ) RETURNING *;

-- name: GetCricketTeamPlayerScore :many
SELECT * FROM cricket_team_player_score
WHERE match_id=$1 AND tournament_id=$2 AND team_id=$3;

-- name: UpdateCricketTeamPlayerScoreBatting :one
UPDATE cricket_team_player_score
SET position=$1, runs_scored=$2, balls_faced=$3, fours=$4, sixes=$5, wicket_taken_by=$6
WHERE tournament_id=$7 AND match_id=$8 AND team_id=$9 AND player_id=$10
RETURNING *;

-- name: UpdateCricketTeamPlayerScoreBowling :one
UPDATE cricket_team_player_score
SET overs_bowled=$1, runs_conceded=$2, wickets_taken=$3
WHERE tournament_id=$4 AND match_id=$5 AND team_id=$6 AND player_id=$7
RETURNING *;