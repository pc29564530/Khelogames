package database

import (
	"context"
	"encoding/json"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const addFootballSquad = `
WITH matchID AS (
	SELECT * FROM matches WHERE id = $1
),
teamID AS (
	SELECT * FROM teams WHERE id = $2
),
playerID AS (
	SELECT * FROM players WHERE id = $3
)
INSERT INTO football_squad (
    match_id,
    team_id,
    player_id,
	is_substitute,
)
SELECT 
	matchID.id,
	teamID.id,
	playerID.id,
	$4
FROM matchID, teamID, playerID
RETURNING *;
`

func (q *Queries) AddFootballSquad(ctx context.Context, matchPublicID, teamPublicID, playerPublicID uuid.UUID, IsSubstitute bool) (models.FootballSquad, error) {
	row := q.db.QueryRowContext(ctx, addFootballSquad, matchPublicID, teamPublicID, playerPublicID, IsSubstitute)
	var i models.FootballSquad
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TeamID,
		&i.PlayerID,
		&i.Position,
		&i.IsSubstitute,
		&i.Role,
		&i.CreatedAT,
	)
	return i, err
}

const getFootballMatchSquad = `
	SELECT
		JSON_BUILD_OBJECT(
			'id', cs.id, 'public_id', cs.public_id, 'match_id', cs.match_id, 'team_id', cs.team_id, 'player_id', cs.player_id, 'position', cs.position, 'is_substitute', cs.is_substitute,  'role', cs.role, 'created_at', cs.created_at,
			'player', JSON_BUILD_OBJECT(
				'id',pl.id,
				'public_id', pl.public_id,
				'user_id',pl.user_Id,
				'name', pl.player_name, 
				'slug', pl.slug, 
				'short_name', pl.short_name, 
				'country', pl.country, 
				'positions', pl.positions, 
				'media_url', pl.media_url
			)
		)
	FROM football_squad as cs
	JOIN players AS pl ON pl.id = cs.player_id
	JOIN teams AS t ON t.id = cs.team_id
	JOIN matches AS m ON m.id = cs.match_id
	WHERE m.public_id=$1 AND m.team_id=$2;
`

func (q *Queries) GetFootballMatchSquad(ctx context.Context, matchPublicID, teamPublicID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getFootballMatchSquad, matchPublicID, teamPublicID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var squads []map[string]interface{}

	for rows.Next() {
		var jsonData []byte
		if err := rows.Scan(&jsonData); err != nil {
			return nil, fmt.Errorf("Failed to scan: %w", err)
		}
		var squad map[string]interface{}
		if err := json.Unmarshal(jsonData, &squad); err != nil {
			return nil, fmt.Errorf("Failed to unmarshal: %w", err)
		}

		squads = append(squads, squad)
	}

	return squads, nil
}
