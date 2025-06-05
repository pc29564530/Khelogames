package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"khelogames/database/models"
)

const addCricketSquad = `
INSERT INTO cricket_squad (
    match_id,
    team_id,
    player_id,
    role,
	on_bench,
	created_at,
	is_captain
) VALUES ( $1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6 )
RETURNING *;
`

func (q *Queries) AddCricketSquad(ctx context.Context, matchID, teamID, playerID int64, role string, OnBench, isCaptain bool) (models.CricketSquad, error) {
	row := q.db.QueryRowContext(ctx, addCricketSquad, matchID, teamID, playerID, role, OnBench, isCaptain)
	var i models.CricketSquad
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.PlayerID,
		&i.Role,
		&i.OnBench,
		&i.CreatedAT,
		&i.IsCaptain,
	)
	return i, err
}

const getCricketMatchSquad = `
	SELECT 
		JSON_AGG(
			JSON_BUILD_OBJECT(
				'id', cs.id, 'match_id', cs.match_id, 'team_id', cs.team_id, 'player_id', cs.player_id, 'role', cs.role, 
				'on_bench', cs.on_bench, 'is_captain', cs.is_captain, 'created_at', cs.created_at,
				'player', JSON_BUILD_OBJECT(
					'id',pl.id,
					'username',pl.username,
					'name', pl.player_name, 
					'slug', pl.slug, 
					'short_name', pl.short_name, 
					'country', pl.country, 
					'positions', pl.positions, 
					'media_url', pl.media_url
				)
			)
		) AS teamSquad
	FROM cricket_squad as cs
	JOIN players AS pl ON pl.id = cs.player_id
	WHERE match_id=$1 AND team_id=$2;
`

func (q *Queries) GetCricketMatchSquad(ctx context.Context, matchID, teamID int64) ([]interface{}, error) {
	row := q.db.QueryRowContext(ctx, getCricketMatchSquad, matchID, teamID)
	var jsonData []byte
	err := row.Scan(&jsonData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}

	if len(jsonData) == 0 {
		return nil, nil
	}

	var teamSquads []interface{}
	if err := json.Unmarshal(jsonData, &teamSquads); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal: %w", err)
	}

	return teamSquads, nil
}
