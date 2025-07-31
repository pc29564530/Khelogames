package token

import (
	"time"

	"github.com/google/uuid"
)

// Make is a interface for mannaging a token
type Maker interface {
	//CreateToken create a token for specifc user
	CreateToken(publicID uuid.UUID, duration time.Duration) (string, *Payload, error)

	//verify token
	VerifyToken(token string) (*Payload, error)
}
