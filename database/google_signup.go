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
) RETURNING *
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

// const getGoogleSignup = `SELECT mobile_number, otp FROM signup
// WHERE mobile_number = $1 LIMIT 1
// `

// func (q *Queries) GetGoogleSignup(ctx context.Context, mobileNumber string) (models.Signup, error) {
// 	var Signup models.Signup
// 	row := q.db.QueryRowContext(ctx, getGoogleSignup, mobileNumber)
// 	err := row.Scan(&Signup.MobileNumber, &Signup.Otp)
// 	return Signup, err
// }
