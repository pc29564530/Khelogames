package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const addCricketSquad = `
WITH matchID AS (
	SELECT * FROM matches WHERE public_id = $1
),
teamID AS (
	SELECT * FROM teams WHERE public_id = $2
),
playerID AS (
	SELECT * FROM players WHERE public_id = $3
)
INSERT INTO cricket_squad (
    match_id,
    team_id,
    player_id,
    role,
	on_bench,
	is_captain,
	created_at
)
SELECT 
	matchID.id,
	teamID.id,
	playerID.id,
	$4,
	$5,
	$6,
	CURRENT_TIMESTAMP
FROM matchID, teamID, playerID
WHERE NOT EXISTS (
	SELECT 1 FROM cricket_squad cs
	WHERE cs.match_id = matchID.id
		AND cs.team_id = teamID.id
		AND cs.player_id = playerID.id
)
RETURNING *;
`

func (q *Queries) AddCricketSquad(ctx context.Context, matchPublicID, teamPublicID, playerPublicID uuid.UUID, role string, OnBench, isCaptain bool) (*models.CricketSquad, error) {
	row := q.db.QueryRowContext(ctx, addCricketSquad, matchPublicID, teamPublicID, playerPublicID, role, OnBench, isCaptain)
	var i models.CricketSquad
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.PlayerID,
		&i.Role,
		&i.OnBench,
		&i.IsCaptain,
		&i.CreatedAT,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("player already exist")
		}
		return nil, err
	}
	return &i, err
}

const getCricketMatchSquad = `
	SELECT 
		JSON_AGG(
			JSON_BUILD_OBJECT(
				'id', cs.id, 'public_id', cs.public_id, 'match_id', cs.match_id, 'team_id', cs.team_id, 'player_id', cs.player_id, 'role', cs.role, 
				'on_bench', cs.on_bench, 'is_captain', cs.is_captain, 'created_at', cs.created_at,
				'player', JSON_BUILD_OBJECT(
					'id',pl.id,
					'public_id', pl.public_id,
					'user_id', pl.user_id,
					'name', pl.name, 
					'slug', pl.slug, 
					'short_name', pl.short_name, 
					'country', pl.country, 
					'positions', pl.positions, 
					'media_url', pl.media_url
				)
			)
		) AS teamSquad
	FROM cricket_squad cs
	JOIN players AS pl ON pl.id = cs.player_id
	JOIN teams AS t ON t.id = cs.team_id
	JOIN matches AS m ON m.id = cs.match_id
	WHERE m.public_id=$1 AND t.public_id=$2;
`

func (q *Queries) GetCricketMatchSquad(ctx context.Context, matchPublicID, teamPublicID uuid.UUID) ([]interface{}, error) {
	row := q.db.QueryRowContext(ctx, getCricketMatchSquad, matchPublicID, teamPublicID)
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
