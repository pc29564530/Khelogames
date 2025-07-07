package database

import (
	"context"
	"khelogames/database/models"
	"log"
)

const createMobileSignup = `INSERT INTO signup (
    mobile_number,
    otp
) VALUES (
    $1, $2
) RETURNING *;
`

func (q *Queries) CreateMobileSignup(ctx context.Context, mobileNumber, otp string) (models.Signup, error) {

	var Signup models.Signup
	row := q.db.QueryRowContext(ctx, createMobileSignup, mobileNumber, otp)
	err := row.Scan(&Signup.MobileNumber, &Signup.Otp)
	if err != nil {
		log.Printf("Failed to create mobile signup: %v", err)
		return Signup, err
	}
	return Signup, nil
}

const deleteSignup = `
DELETE FROM signup
WHERE mobile_number = $1 RETURNING mobile_number, otp
`

func (q *Queries) DeleteSignup(ctx context.Context, mobileNumber string) (models.Signup, error) {
	var Signup models.Signup
	row := q.db.QueryRowContext(ctx, deleteSignup, mobileNumber)
	err := row.Scan(&Signup.MobileNumber, &Signup.Otp)
	return Signup, err
}

const getSignup = `SELECT mobile_number, otp FROM signup
WHERE mobile_number = $1 LIMIT 1
`

func (q *Queries) GetSignup(ctx context.Context, mobileNumber string) (models.Signup, error) {
	var Signup models.Signup
	row := q.db.QueryRowContext(ctx, getSignup, mobileNumber)
	err := row.Scan(&Signup.MobileNumber, &Signup.Otp)
	return Signup, err
}

const createGoogleSignUp = `
	INSERT INTO modify_user (
		full_name, username, email, hash_password, is_verified, is_banned, google_id, role, created_at, updated_at
	) VALUES (
	 $1, $2, $3, $4, false, false, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP 
	) RETURNING *;
`

func (q *Queries) CreateGoogleSignUp(ctx context.Context, username, email, fullName, googleId string) (models.ModifyUser, error) {
	var users models.ModifyUser
	row := q.db.QueryRowContext(ctx, createGoogleSignUp, username, email, fullName, googleId)
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
	if err != nil {
		log.Printf("Failed to create google signup: %v", err)
	}
	return users, nil
}

const createEmailSignUp = `
	INSERT INTO modify_user (
		full_name, username, email, hash_password, is_verified, is_banned, google_id, role, created_at, updated_at
	) VALUES (
	 $1, $2, $3, $4, false, false, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP 
	) RETURNING *;
`

func (q *Queries) CreateEmailSignUp(ctx context.Context, username, email, fullName, googleId string) (models.ModifyUser, error) {
	var users models.ModifyUser
	row := q.db.QueryRowContext(ctx, createGoogleSignUp, username, email, fullName, googleId)
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
	if err != nil {
		log.Printf("Failed to create gmail signup: %v", err)
	}
	return users, nil
}
