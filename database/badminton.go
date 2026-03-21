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
	var i models.BadmintonScore
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
	return &i, nil
}

const updateBadmintonScore = `
	UPDATE badminton_score bs
	SET
		home_score = CASE
			WHEN m.home_team_id = (SELECT id FROM teams WHERE public_id = $2)
			THEN COALESCE(bs.home_score, 0) + 1
			ELSE bs.home_score
		END,
		away_score = CASE
			WHEN m.away_team_id = (SELECT id FROM teams WHERE public_id = $2)
			THEN COALESCE(bs.away_score, 0) + 1
			ELSE bs.away_score
		END
	FROM matches m
	WHERE bs.match_id = m.id
	AND m.public_id = $1
	AND bs.set_number = $3
	RETURNING bs.*;
`

func (q *Queries) UpdateBadmintonScore(ctx context.Context, matchPublicID, teamPublicID uuid.UUID, setNumber int) (*models.BadmintonScore, error) {
	row := q.db.QueryRowContext(ctx, updateBadmintonScore, matchPublicID, teamPublicID, setNumber)
	var i models.BadmintonScore
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
		return nil, fmt.Errorf("Failed to scan: ", err)
	}
	return &i, nil
}

const getBadmintonMatchSetsScore = `
	SELECT
		JSON_BUILD_OBJECT(
			'public_id', bs.public_id, 'match_public_id', m.public_id, 'set_number', bs.set_number, 'home_score', bs.home_score, 'away_score', bs.away_score, 'set_status', bs.set_status
		)
	FROM badminton_score bs
	LEFT JOIN matches AS m ON m.id=bs.match_id
	WHERE m.public_id=$1
	ORDER BY bs.set_number;
`

func (q *Queries) GetBadmintonMatchSetsScore(ctx context.Context, matchPublicID uuid.UUID) ([]map[string]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getBadmintonMatchSetsScore, matchPublicID)
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

const getBadmintonMatchScore = `
	SELECT
    CASE 
        WHEN m.status_code = 'in_progress'
        THEN COUNT(*) FILTER (
            WHERE bs.home_score IS NOT NULL
			AND bs.away_score IS NOT NULL
            AND bs.home_score > bs.away_score
        )
		
		WHEN m.status_code = 'finished'
		THEN COUNT(*) FILTER (
			WHERE bs.set_status = 'finished'
			AND bs.home_score > bs.away_score
		)

		ELSE NULL
    END AS home_sets_won,

    CASE
        WHEN m.status_code = 'in_progress'
        THEN COUNT(*) FILTER (
            WHERE bs.home_score IS NOT NULL
			AND bs.away_score IS NOT NULL
            AND bs.away_score > bs.home_score
        )
		
		WHEN m.status_code = 'finished'
		THEN COUNT(*) FILTER (
			WHERE bs.set_status = 'finished'
			AND bs.away_score > bs.home_score
		)

		ELSE NULL
    END AS away_sets_won
	
	FROM badminton_score bs
	JOIN matches m ON m.id = bs.match_id
	WHERE m.public_id = $1
	GROUP BY m.status_code;
`

type setsWon struct {
	HomeSetsWon *int `json:"home_sets_won"`
	AwaySetsWon *int `json:"away_sets_won"`
}

func (q *Queries) GetBadmintonMatchScore(ctx context.Context, matchPublicID uuid.UUID) (*setsWon, error) {
	rows := q.db.QueryRowContext(ctx, getBadmintonMatchScore, matchPublicID)
	var result setsWon
	err := rows.Scan(&result.HomeSetsWon, &result.AwaySetsWon)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}

	return &result, nil
}

const getBadmintonMatchSetScore = `
	SELECT bs.*
	FROM badminton_score bs
	LEFT JOIN matches AS m ON m.id=bs.match_id
	WHERE m.public_id=$1 AND bs.set_number = $2
	ORDER BY bs.set_number;
`

func (q *Queries) GetBadmintonMatchSetScore(ctx context.Context, matchPublicID uuid.UUID, setNumber int) (*models.BadmintonScore, error) {
	var i models.BadmintonScore
	rows := q.db.QueryRowContext(ctx, getBadmintonMatchSetScore, matchPublicID, setNumber)
	err := rows.Scan(
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

	return &i, nil
}

const updateBadmintonSetStatus = `
	UPDATE badminton_score bs
	SET
		set_status = $3
	FROM matches m
	WHERE bs.match_id = m.id
	AND m.public_id = $1
	AND bs.set_number = $2
	RETURNING bs.*;
`

func (q *Queries) UpdateBadmintonSetStatus(ctx context.Context, matchPublicID uuid.UUID, setNumber int, setStatus string) (*models.BadmintonScore, error) {
	row := q.db.QueryRowContext(ctx, updateBadmintonSetStatus, matchPublicID, setNumber, setStatus)
	var i models.BadmintonScore
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
		return nil, fmt.Errorf("Failed to scan: ", err)
	}
	return &i, nil
}
