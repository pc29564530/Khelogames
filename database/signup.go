package database

import (
	"context"
	"khelogames/database/models"
	"log"
)

const createGoogleSignUp = `
	INSERT INTO users (
		full_name, username, email, hash_password, is_verified, is_banned, google_id, role, created_at, updated_at
	) VALUES (
	 $1, $2, $3, null, false, false, $4, 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP 
	) RETURNING *;
`

func (q *Queries) CreateGoogleSignUp(ctx context.Context, fullName, username, email, googleId string) (*models.Users, error) {
	var users models.Users
	row := q.db.QueryRowContext(ctx, createGoogleSignUp, fullName, username, email, googleId)
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
	if err != nil {
		log.Printf("Failed to create google signup: %v", err)
	}
	return &users, nil
}

const createEmailSignUp = `
	INSERT INTO users (
		full_name, username, email, hash_password, is_verified, is_banned, google_id, role, created_at, updated_at
	) VALUES (
	 $1, $2, $3, $4, false, false, null, 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
	) RETURNING *;
`

func (q *Queries) CreateEmailSignUp(ctx context.Context, fullName, username, email, hashPassword string) (*models.Users, error) {
	var users models.Users
	row := q.db.QueryRowContext(ctx, createEmailSignUp, fullName, username, email, hashPassword)
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
	if err != nil {
		log.Printf("Failed to create gmail signup: %v", err)
	}
	return &users, nil
}
