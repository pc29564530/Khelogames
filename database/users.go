package database

import (
	"context"
	"khelogames/database/models"
)

const getUser = `
SELECT * FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (models.Users, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var users models.Users
	err := row.Scan(
		&users.ID,
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
WHERE username = $1
ORDER BY username
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
