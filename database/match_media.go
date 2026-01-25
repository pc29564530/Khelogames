package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const createMatchMedia = `
	WITH resolve_ids AS (
		SELECT * FROM matches WHERE public_id = $2
	)
	INSERT INTO match_highlights (
		user_id,
		tournament_id,
		match_id,
		title,
		description,
		media_url,
		created_at,
		updated_at
	)
	SELECT
		$1,
		rsi.tournament_id,
		rsi.id,
		$3,
		$4,
		$5,
		NOW(),
		NOW()
	FROM resolve_ids rsi
	RETURNING *
`

func (q *Queries) CreateMatchMedia(ctx context.Context, userPublicID int32, matchPublicID uuid.UUID, title, description, mediaURL string) (*models.MatchHighlights, error) {

	var matchHighlights models.MatchHighlights
	rows := q.db.QueryRowContext(ctx, createMatchMedia, userPublicID, matchPublicID, title, description, mediaURL)
	err := rows.Scan(
		&matchHighlights.ID,
		&matchHighlights.PublicID,
		&matchHighlights.UserID,
		&matchHighlights.TournamentID,
		&matchHighlights.MatchID,
		&matchHighlights.Title,
		&matchHighlights.Description,
		&matchHighlights.MediaURL,
		&matchHighlights.CreatedAT,
		&matchHighlights.UpdatedAT,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}
	return &matchHighlights, err
}

const getMatchMedia = `
	SELECT mh.* FROM match_highlights mh
	LEFT JOIN matches m ON m.id = mh.match_id
	WHERE m.public_id = $1
`

func (q *Queries) GetMatchMedia(ctx context.Context, matchPublicID uuid.UUID) (*[]models.MatchHighlights, error) {
	rows, err := q.db.QueryContext(ctx, getMatchMedia, matchPublicID)
	if err != nil {
		return nil, fmt.Errorf("Failed to query row: ", err)
	}

	defer rows.Close()
	var matchHighlights []models.MatchHighlights
	for rows.Next() {
		var i models.MatchHighlights
		err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.UserID,
			&i.TournamentID,
			&i.MatchID,
			&i.MediaURL,
			&i.Title,
			&i.Description,
			&i.CreatedAT,
			&i.UpdatedAT,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan match highlights: ", err)
		}
		matchHighlights = append(matchHighlights, i)
	}
	return &matchHighlights, nil
}
