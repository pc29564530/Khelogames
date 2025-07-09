package database

import (
	"context"
	"khelogames/database/models"
)

const createUser = `
INSERT INTO users (
  username,
  mobile_number,
  role,
  gmail
) VALUES (
  $1, $2, $3, $4
) RETURNING *`

func (q *Queries) CreateUser(ctx context.Context, username string, mobileNumber string, role string, gmail string) (models.User, error) {
	var users models.User
	row := q.db.QueryRowContext(ctx, createUser,
		username,
		mobileNumber,
		role,
		gmail,
	)
	err := row.Scan(
		&users.Username,
		&users.MobileNumber,
		&users.Role,
		&users.Gmail,
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
			&user.Role,
			&user.Gmail,
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

const getUserByMobileNumber = `
SELECT * FROM users
WHERE mobile_number = $1
`

func (q *Queries) GetUserByMobileNumber(ctx context.Context, mobile_number string) (models.User, error) {
	row := q.db.QueryRowContext(ctx, getUserByMobileNumber, mobile_number)
	var users models.User
	err := row.Scan(
		&users.Username,
		&users.MobileNumber,
		&users.Role,
		&users.Gmail,
	)
	return users, err
}

const getUserByGmail = `
SELECT * FROM users
WHERE gmail = $1
`

func (q *Queries) GetUserByGmail(ctx context.Context, gmail string) (models.User, error) {
	row := q.db.QueryRowContext(ctx, getUserByGmail, gmail)
	var users models.User
	err := row.Scan(
		&users.Username,
		&users.MobileNumber,
		&users.Role,
		&users.Gmail,
	)
	return users, err
}

const getModifyUserByGmail = `
SELECT * FROM modify_user
WHERE email = $1;
`

func (q *Queries) GetModifyUserByGmail(ctx context.Context, gmail string) (*models.ModifyUser, error) {
	row := q.db.QueryRowContext(ctx, getModifyUserByGmail, gmail)
	var users models.ModifyUser
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
