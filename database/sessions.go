package database

import (
	"context"
	"khelogames/database/models"
	"time"

	"github.com/google/uuid"
)

const createSessions = `
INSERT INTO sessions (
    user_id,
    refresh_token,
    user_agent,
    client_ip,
	created_at,
    expires_at
) VALUES (
    $1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
) RETURNING *;
`

type CreateSessionsParams struct {
	UserID       int32     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (q *Queries) CreateSessions(ctx context.Context, arg CreateSessionsParams) (models.Session, error) {
	row := q.db.QueryRowContext(ctx, createSessions,
		arg.UserID,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.CreatedAt,
		arg.ExpiresAt,
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
	return session, err
}

const deleteSessions = `
	DELETE FROM sessions s
	JOIN users AS u ON u.id = s.user_id
	WHERE u.public_id = $1
`

func (q *Queries) DeleteSessions(ctx context.Context, publicID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteSessions, publicID)
	return err
}

const getSessions = `
	SELECT * FROM sessions
	WHERE public_id = $1
`

func (q *Queries) GetSessions(ctx context.Context, publicID uuid.UUID) (models.Session, error) {
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
	return session, err
}
