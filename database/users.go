package database

import (
	"context"
	"khelogames/database/models"
	"time"

	"github.com/google/uuid"
)

const getUser = `
SELECT * FROM users
WHERE public_id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, publicID uuid.UUID) (models.Users, error) {
	row := q.db.QueryRowContext(ctx, getUser, publicID)
	var users models.Users
	err := row.Scan(
		&users.ID,
		&users.PublicID,
		&users.FullName,
		&users.Username,
		&users.Email,
		&users.HashPassword,
		&users.IsVerified,
		&users.IsBanned,
		&users.GoogleID,
		&users.Role,
		&users.CreatedAt,
		&users.UpdatedAt,
	)
	return users, err
}

const listUser = `
SELECT * FROM users
WHERE public_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

func (q *Queries) ListUser(ctx context.Context, pageSize, offSet int32) ([]models.Users, error) {
	rows, err := q.db.QueryContext(ctx, listUser, pageSize, offSet)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Users
	for rows.Next() {
		var users models.Users
		if err := rows.Scan(
			&users.ID,
			&users.PublicID,
			&users.FullName,
			&users.Username,
			&users.Email,
			&users.HashPassword,
			&users.IsVerified,
			&users.IsBanned,
			&users.GoogleID,
			&users.Role,
			&users.CreatedAt,
			&users.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, users)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersByGmail = `
SELECT * FROM users
WHERE email = $1;
`

func (q *Queries) GetUsersByGmail(ctx context.Context, gmail string) (*models.Users, error) {
	row := q.db.QueryRowContext(ctx, getUsersByGmail, gmail)
	var users models.Users
	err := row.Scan(
		&users.ID,
		&users.PublicID,
		&users.FullName,
		&users.Username,
		&users.Email,
		&users.HashPassword,
		&users.IsVerified,
		&users.IsBanned,
		&users.GoogleID,
		&users.Role,
		&users.CreatedAt,
		&users.UpdatedAt,
	)
	return &users, err
}

const updateUsersFullName = `
	UPDATE users
	SET full_name = $2
	WHERE id = $1
	RETURNING *;
`

func (q *Queries) UpdateUser(ctx context.Context, userID int32, fullName string) (models.Users, error) {
	row := q.db.QueryRowContext(ctx, updateUsersFullName,
		userID,
		fullName,
	)
	var users models.Users
	err := row.Scan(
		&users.ID,
		&users.PublicID,
		&users.FullName,
		&users.Username,
		&users.Email,
		&users.HashPassword,
		&users.IsVerified,
		&users.IsBanned,
		&users.GoogleID,
		&users.Role,
		&users.CreatedAt,
		&users.UpdatedAt,
	)
	return users, err
}

type searchUserParam struct {
	ID        int64     `json:"id"`
	PublicID  uuid.UUID `json:"public_id"`
	UserID    int32     `json:"user_id"`
	UserName  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const searchUser = `
SELECT
	up.id AS id,
	up.public_id AS public_id,
	u.id AS user_id,
	u.username AS username,
	u.full_name AS full_name,
	COALESCE(up.bio, '') AS bio,
	COALESCE(up.avatar_url, '') AS avatar_url,
	u.created_at AS created_at,
	u.updated_at AS updated_at
FROM users u
LEFT JOIN user_profiles AS up ON u.id = up.user_id
WHERE u.full_name ILIKE $1
ORDER BY u.full_name
LIMIT 10
`

func (q *Queries) SearchUser(ctx context.Context, searchTerm string) ([]searchUserParam, error) {
	// Add wildcards for LIKE query
	pattern := "%" + searchTerm + "%"

	rows, err := q.db.QueryContext(ctx, searchUser, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []searchUserParam
	for rows.Next() {
		var i searchUserParam
		if err := rows.Scan(
			&i.ID,
			&i.PublicID,
			&i.UserID,
			&i.UserName,
			&i.FullName,
			&i.Bio,
			&i.AvatarURL,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
