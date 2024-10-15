package database

import (
	"context"
	"database/sql"
	"khelogames/database/models"
	"log"
)

const createGoogleUser = `INSERT INTO users (
	username,
    gmail,
	role
) VALUES (
    $1, $2, $3
) RETURNIN *
`

func (q *Queries) CreateGoogleUser(ctx context.Context, username, gmail, role string) (models.User, error) {
	var User models.User
	row := q.db.QueryRowContext(ctx, createGoogleUser, username, gmail, role)
	err := row.Scan(&User.Username, &User.MobileNumber, &User.HashedPassword, &User.Role, &User.Gmail)
	if err != nil {
		log.Printf("Failed to create google user: %v", err)
		return models.User{}, err
	}
	return User, err
}

func (q *Queries) GetGoogleMailID(ctx context.Context, gmail string) (models.User, error) {
	var User models.User
	row := q.db.QueryRowContext(ctx, gmail)
	err := row.Scan(&User.Username, &User.Role, &User.Gmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil
		}
		return User, err
	}
	return User, nil
}
