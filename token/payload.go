package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

//Payload contain the payload data of the token

type Payload struct {
	ID           uuid.UUID `json:"id"`
	PublicID     uuid.UUID `json:"public_id"`
	UserID       int32     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	IssueAt      time.Time `json:"issue_at"`
	ExpiredAt    time.Time `json:"expired_at"`
}

// new payload creates a new token paylaod with specific username and duration

func NewPayload(publicID uuid.UUID, userID int32, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		PublicID:  publicID,
		UserID:    userID,
		IssueAt:   time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

//valid checks is function payload is valid or not

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
