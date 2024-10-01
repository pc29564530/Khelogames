package new_db

import (
	"context"
	"khelogames/new_db/models"
	"time"

	"github.com/google/uuid"
)

const createSessions = `
INSERT INTO sessions (
    id,
    username,
    refresh_token,
    user_agent,
    client_ip,
    expires_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP
) RETURNING id, username, refresh_token, user_agent, client_ip, expires_at, created_at
`

type CreateSessionsParams struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (q *Queries) CreateSessions(ctx context.Context, arg CreateSessionsParams) (models.Session, error) {
	row := q.db.QueryRowContext(ctx, createSessions,
		arg.Username,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.ExpiresAt,
	)
	var session models.Session
	err := row.Scan(
		&session.ID,
		&session.Username,
		&session.RefreshToken,
		&session.UserAgent,
		&session.ClientIp,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	return session, err
}

const deleteSessions = `
DELETE FROM sessions
WHERE username = $1
`

func (q *Queries) DeleteSessions(ctx context.Context, username string) error {
	_, err := q.db.ExecContext(ctx, deleteSessions, username)
	return err
}

const getSessions = `
SELECT id, username, refresh_token, user_agent, client_ip, expires_at, created_at FROM sessions
WHERE username = $1
`

func (q *Queries) GetSessions(ctx context.Context, username string) (models.Session, error) {
	row := q.db.QueryRowContext(ctx, getSessions, username)
	var session models.Session
	err := row.Scan(
		&session.ID,
		&session.Username,
		&session.RefreshToken,
		&session.UserAgent,
		&session.ClientIp,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	return session, err
}
