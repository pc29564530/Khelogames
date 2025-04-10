package database

import (
	"context"
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
	on_bench
) VALUES ( $1, $2, $3, $4, $5 )
RETURNING id, team_id, player_id, match_id, role, on_bench, created_at
`

func (q *Queries) AddCricketSquad(ctx context.Context, matchID, teamID, playerID int64, role string, OnBench bool) (models.CricketSquad, error) {
	row := q.db.QueryRowContext(ctx, addCricketSquad, matchID, teamID, playerID, role, OnBench)
	var i models.CricketSquad
	err := row.Scan(
		&i.ID,
		&i.MatchID,
		&i.TeamID,
		&i.PlayerID,
		&i.Role,
		&i.OnBench,
		&i.CreatedAT,
	)
	return i, err
}

const getCricketMatchSquad = `
	SELECT 
		JSON_BUILD_OBJECT(
			'id', cs.id, 'match_id', cs.match_id, 'team_id', cs.team_id, 'player_id', cs.player_id, 'role', cs.role, 
			'on_bench', cs.on_bench, 'created_at', cs.created_at,
			'player', JSON_BUILD_OBJECT('id',pl.id,'username',pl.username, 'name', pl.player_name, 'slug', pl.slug, 'short_name',pl.short_name, 'country', pl.country, 'positions', pl.positions, 'media_url', pl.media_url)
		) AS teamSquad
	FROM cricket_squad as cs
	LEFT JOIN players AS pl ON pl.id = cs.player_id
	WHERE match_id=$1 AND team_id=$2;
`

func (q *Queries) GetCricketMatchSquad(ctx context.Context, matchID, teamID int64) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getCricketMatchSquad, matchID, teamID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query: ", err)
	}

	defer rows.Close()

	var teamSquads []map[string]interface{}

	if rows.Next() {
		var jsonData []byte
		if err := rows.Scan(&jsonData); err != nil {
			return nil, fmt.Errorf("Failed to scan data: ", err)
		}
		var squad map[string]interface{}
		err := json.Unmarshal(jsonData, &squad)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal the json data: ", err)
		}
		teamSquads = append(teamSquads, squad)
	}

	return teamSquads, nil
}
