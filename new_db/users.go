package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const createUser = `
INSERT INTO users (
  username,
  mobile_number,
  hashed_password,
  role
) VALUES (
  $1, $2, $3, $4
) RETURNING username, mobile_number, hashed_password, role
`

func (q *Queries) CreateUser(ctx context.Context, username string, mobileNumber string, hashedPassword string, role string) (models.User, error) {
	var users models.User
	row := q.db.QueryRowContext(ctx, createUser,
		username,
		mobileNumber,
		hashedPassword,
		role,
	)
	err := row.Scan(
		&users.Username,
		&users.MobileNumber,
		&users.HashedPassword,
		&users.Role,
	)
	return users, err
}

const getUser = `
SELECT username, mobile_number, hashed_password, role FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (models.User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var users models.User
	err := row.Scan(
		&users.Username,
		&users.MobileNumber,
		&users.HashedPassword,
		&users.Role,
	)
	return users, err
}

const listUser = `
SELECT DISTINCT username, mobile_number, hashed_password, role FROM users
WHERE username = $1
ORDER BY username
LIMIT $2
OFFSET $3
`

func (q *Queries) ListUser(ctx context.Context, username string) ([]models.User, error) {
	rows, err := q.db.QueryContext(ctx, listUser, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.Username,
			&user.MobileNumber,
			&user.HashedPassword,
			&user.Role,
		); err != nil {
			return nil, err
		}
		items = append(items, user)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
