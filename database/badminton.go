package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const addBadmintonScore = `
	INSERT INTO badminton_score (
		match_id,
		set_number
	)
	VALUES (
		$1, $2
	) RETURNING *;
`

func (q *Queries) AddBadmintonScore(ctx context.Context, matchID int32, setNumber int) (*models.BadmintonScore, error) {
	row := q.db.QueryRowContext(ctx, addBadmintonScore, matchID, setNumber)
	var i *models.BadmintonScore
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.SetNumber,
		&i.HomeScore,
		&i.AwayScore,
		&i.SetStatus,
		&i.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}
	return i, nil
}

const updateBadmintonScore = `
	UDPATE badminton_score bs
	SET
		home_score = CASE
			WHEN m.home_team_id = (SELECT id FROM teams WHERE public_id = $2) THEN COALESCE(bs.home_score, 0) + 1
			ELSE bs.home_score
		END,
		away_score = CASE
			WHEN m.away_team_id = (SELECT id FROM teams WHERE public_id = $2) THEN COALESCE(bs.away_score, 0) + 1
		END
	FROM matches m
	WHERE bs.match_id = m.id
	AND m.public_id = $1
	AND bs.set_number = $3
	RETURNING *;
`

func (q *Queries) UpdateBadmintonScore(ctx context.Context, matchPublicID, teamPublicID uuid.UUID, setNumber int) (*models.BadmintonScore, error) {
	row := q.db.QueryRowContext(ctx, addBadmintonScore, matchPublicID, teamPublicID, setNumber)
	var i *models.BadmintonScore
	err := row.Scan(
		&i.ID,
		&i.PublicID,
		&i.MatchID,
		&i.SetNumber,
		&i.HomeScore,
		&i.AwayScore,
		&i.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}
	return i, nil
}

const getBadmintonScore = `
	SELECT
		JSON_BUILD_OBJECT(
			'public_id', bs.public_id, 'match_public_id', m.public_id, 'set_number', bs.set_number, 'home_score', bs.home_score, 'away_score', bs.away_score, 'set_status', bs.set_status,
		)
	FROM badminton_score bs
	LEFT JOIN matches AS m ON m.id=bs.match_id
	WHERE m.public_id=$1
	ORDER BY bs.set_number;
`

func (q *Queries) GetBadmintonScore(ctx context.Context, matchPublicID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getBadmintonScore, matchPublicID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}
	var scores []map[string]interface{}
	for rows.Next() {
		var score map[string]interface{}
		var jsonByte []byte
		err := rows.Scan(&jsonByte)
		if err := rows.Scan(&jsonByte); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		err = json.Unmarshal(jsonByte, &score)
		if err != nil {
			return nil, fmt.Errorf("Failed to unmarshal: ", err)
		}
		scores = append(scores, score)
	}
	return scores, nil
}
