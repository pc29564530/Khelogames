package database

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/database/models"
	"time"

	"github.com/google/uuid"
)

const createSessions = `
INSERT INTO sessions (
    user_id, refresh_token, user_agent, client_ip, created_at, expires_at
) VALUES (
    $1, $2, $3, $4, CURRENT_TIMESTAMP, $5
)
RETURNING *;
`

type CreateSessionsParams struct {
	UserID       int32     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	ExpiresTime  time.Time `json:"expire_time"`
}

func (q *Queries) CreateSessions(ctx context.Context, arg CreateSessionsParams) (*models.Session, error) {
	row := q.db.QueryRowContext(ctx, createSessions,
		arg.UserID,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.ExpiresTime,
	)
	var session models.Session
	err := row.Scan(
		&session.ID,
		&session.PublicID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.ClientIp,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no rows in result")
		}
		return nil, fmt.Errorf("Failed to scan: ", err)
	}
	return &session, nil
}

const deleteSessions = `
	DELETE FROM sessions s
	WHERE user_id IN (
		SELECT id FROM users WHERE public_id = $1
	);
`

func (q *Queries) DeleteSessions(ctx context.Context, publicID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteSessions, publicID)
	return err
}

const deleteSessionsByUserID = `
	DELETE FROM sessions s
	WHERE s.user_id = $1;
`

func (q *Queries) DeleteSessionsByUserID(ctx context.Context, userID int32) error {
	_, err := q.db.ExecContext(ctx, deleteSessionsByUserID, userID)
	return err
}

const getSessions = `
	SELECT * FROM sessions
	WHERE public_id = $1
`

func (q *Queries) GetSessions(ctx context.Context, publicID uuid.UUID) (*models.Session, error) {
	row := q.db.QueryRowContext(ctx, getSessions, publicID)
	var session models.Session
	err := row.Scan(
		&session.ID,
		&session.PublicID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.ClientIp,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to scan session: %w", err)
	}

	return &session, nil
}

func (store *Store) GetSessionByUserID(ctx context.Context, userID int32) (*models.Session, error) {
	const query = `
		SELECT * FROM sessions 
		WHERE user_id = $1`

	var session models.Session
	err := store.db.QueryRowContext(ctx, query, userID).Scan(
		&session.ID,
		&session.PublicID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.ClientIp,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to scan session: %w", err)
	}

	return &session, nil
}

func (store *Store) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	const query = `
		SELECT 
			* 
		FROM sessions 
		WHERE refresh_token = $1 AND expires_at > CURRENT_TIMESTAMP;`

	var session models.Session
	err := store.db.QueryRowContext(ctx, query, refreshToken).Scan(
		&session.ID,
		&session.PublicID,
		&session.UserID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.ClientIp,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to scan session: %w", err)
	}

	return &session, nil
}
