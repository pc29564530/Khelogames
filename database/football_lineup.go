package database

import (
	"context"
	"khelogames/database/models"
)

const addFootballLineUp = `
INSERT INTO football_lineup (
    team_id,
    player_id,
    match_id,
    position
) VALUES ( $1, $2, $3, $4 )
RETURNING id, team_id, player_id, match_id, position
`

type AddFootballLineUpParams struct {
	TeamID   int64  `json:"team_id"`
	PlayerID int64  `json:"player_id"`
	MatchID  int64  `json:"match_id"`
	Position string `json:"position"`
}

func (q *Queries) AddFootballLineUp(ctx context.Context, arg AddFootballLineUpParams) (models.FootballLineup, error) {
	row := q.db.QueryRowContext(ctx, addFootballLineUp,
		arg.TeamID,
		arg.PlayerID,
		arg.MatchID,
		arg.Position,
	)
	var i models.FootballLineup
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.PlayerID,
		&i.MatchID,
		&i.Position,
	)
	return i, err
}

const getFootballLineUp = `
SELECT id, team_id, player_id, match_id, position FROM football_lineup
WHERE match_id=$1 AND team_id=$2
`

type GetFootballLineUpParams struct {
	MatchID int64 `json:"match_id"`
	TeamID  int64 `json:"team_id"`
}

func (q *Queries) GetFootballLineUp(ctx context.Context, arg GetFootballLineUpParams) ([]models.FootballLineup, error) {
	rows, err := q.db.QueryContext(ctx, getFootballLineUp, arg.MatchID, arg.TeamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.FootballLineup
	for rows.Next() {
		var i models.FootballLineup
		if err := rows.Scan(
			&i.ID,
			&i.TeamID,
			&i.PlayerID,
			&i.MatchID,
			&i.Position,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateFootballSubsAndLineUp = `
WITH 
    sub AS (
        SELECT fs.position, fs.player_id 
        FROM football_substitution fs
        WHERE fs.id = $1
    ),
    lu AS (
        SELECT fl.position, fl.player_id 
        FROM football_lineup fl
        WHERE fl.id = $2
    ),
    update_sub AS (
        UPDATE football_substitution fs
        SET 
            position = lu.position, 
            player_id = lu.player_id
        FROM lu
        WHERE fs.id = $1
        RETURNING fs.id, fs.team_id, fs.player_id, fs.match_id, fs.position
    ),
    update_lu AS (
        UPDATE football_lineup fl
        SET 
            position = sub.position, 
            player_id = sub.player_id
        FROM sub
        WHERE fl.id = $2
        RETURNING fl.id, fl.team_id, fl.player_id, fl.match_id, fl.position
    )
SELECT 
    update_sub.id, update_sub.team_id, update_sub.player_id, update_sub.match_id, update_sub.position, 
    update_lu.id, update_lu.team_id, update_lu.player_id, update_lu.match_id, update_lu.position
FROM 
    update_sub, 
    update_lu
`

type UpdateFootballSubsAndLineUpParams struct {
	ID   int64 `json:"id"`
	ID_2 int64 `json:"id_2"`
}

type UpdateFootballSubsAndLineUpRow struct {
	ID         int64  `json:"id"`
	TeamID     int64  `json:"team_id"`
	PlayerID   int64  `json:"player_id"`
	MatchID    int64  `json:"match_id"`
	Position   string `json:"position"`
	ID_2       int64  `json:"id_2"`
	TeamID_2   int64  `json:"team_id_2"`
	PlayerID_2 int64  `json:"player_id_2"`
	MatchID_2  int64  `json:"match_id_2"`
	Position_2 string `json:"position_2"`
}

func (q *Queries) UpdateFootballSubsAndLineUp(ctx context.Context, arg UpdateFootballSubsAndLineUpParams) (UpdateFootballSubsAndLineUpRow, error) {
	row := q.db.QueryRowContext(ctx, updateFootballSubsAndLineUp, arg.ID, arg.ID_2)
	var i UpdateFootballSubsAndLineUpRow
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.PlayerID,
		&i.MatchID,
		&i.Position,
		&i.ID_2,
		&i.TeamID_2,
		&i.PlayerID_2,
		&i.MatchID_2,
		&i.Position_2,
	)
	return i, err
}
