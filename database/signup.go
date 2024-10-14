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
