package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"

	"github.com/google/uuid"
)

const addMatchUserRole = `
	WITH match_id AS (
		SELECT * FROM matches WHERE public_id = $1
	),
	user_id AS (
		SELECT * FROM profile WHERE public_id = $2
	),
	INSERT INTO match_user_roles (
		match_id,
		user_id,
		role,
		assigned_by,
	)
	SELECT (
		 match_id.id,
		 user_id.id,
		 $3,
		 $4,
	)
	FROM match_id
	JOIN user_id ON TRUE
	RETURNING *;
`

func (q *Queries) AddMatchUserRole(ctx context.Context, matchID, userID int32, role string, assignedBy int32) (*models.MatchUserRoles, error) {
	row := q.db.QueryRowContext(ctx, addMatchUserRole, matchID, userID, role, assignedBy)
	var i models.MatchUserRoles
	err := row.Scan(&i)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}
	return &i, nil
}

const getMatchUserRole = `
	SELECT * FROM match_user_roles mu
	JOIN matches AS m ON m.id = mu.match_id
	WHERE m.public_id = $1 AND mu.user_id = $2 AND mu.is_active = true
	RETURNING *;
`

func (q *Queries) GetMatchUserRole(ctx context.Context, matchPublicID uuid.UUID, userID int32) (*models.MatchUserRoles, error) {
	row := q.db.QueryRowContext(ctx, getMatchUserRole, matchPublicID, userID)
	var i models.MatchUserRoles
	err := row.Scan(&i)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Failed to scan: %w", err)
	}
	return &i, nil
}
