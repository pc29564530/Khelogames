package new_db

import (
	"context"
	"khelogames/new_db/models"
)

const createSignup = `INSERT INTO signup (
    mobile_number,
    otp
) VALUES (
    $1, $2
) RETURNING mobile_number, otp
`

func (q *Queries) CreateSignup(ctx context.Context, mobileNumber, otp string) (models.Signup, error) {
	var Signup models.Signup
	row := q.db.QueryRowContext(ctx, createSignup, mobileNumber, otp)
	err := row.Scan(&Signup.MobileNumber, &Signup.Otp)
	return Signup, err
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
