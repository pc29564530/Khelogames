package database

import (
	"context"
	"encoding/json"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const addCricketToss = `
WITH matchID AS (
	SELECT * FROM matches WHERE public_id=$1
),
teamID AS (
	SELECT * FROM teams WHERE public_id = $3
)
INSERT INTO cricket_toss (
    match_id,
    toss_decision,
    toss_win
) 
SELECT
	matchID.id,
	$2,
	teamID.id
FROM matchID, teamID
RETURNING *
`

func (q *Queries) AddCricketToss(ctx context.Context, matchPublicID uuid.UUID, tossDecision string, tossWinPublicID uuid.UUID) (models.CricketToss, error) {
	row := q.db.QueryRowContext(ctx, addCricketToss, matchPublicID, tossDecision, tossWinPublicID)
	var i models.CricketToss
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.TossDecision,
		&i.TossWin,
	)
	return i, err
}

const getCricketToss = `
SELECT 
	JOIN_BUILD_OBJECT(
		'id', ct.id, 'public_id', ct.public_id, 'match_id', ct.match_id, 'toss_decision', ct.toss_decision, 'toss_win', ct.toss_win,
		'toss_won_team', JSON_BUILD_OBJECT('id', tm.id, 'public_id', tm.public_id, 'user_id', tm.user_id, 'name', tm.name, 'slug', tm.slug, 'short_name', tm.shortname, 'admin', tm.admin, 'media_url', tm.media_url, 'gender', tm.gender, 'national', tm.national, 'country', tm.country, 'type', tm.type, 'player_count', tm.player_count, 'game_id', tm.game_id)
	)
FROM cricket_toss ct
JOIN matches m ON m.id = ct.match_id
WHERE m.public_id=$1
`

func (q *Queries) GetCricketToss(ctx context.Context, matchPublicID uuid.UUID) (map[string]interface{}, error) {
	row := q.db.QueryRowContext(ctx, getCricketToss, matchPublicID)
	var jsonByte []byte
	err := row.Scan(&jsonByte)
	if err != nil {
		return nil, fmt.Errorf("Failed to scan: ", err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(jsonByte, &data)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal: ", err)
	}
	return data, nil
}
