package token

import (
	"time"
)

// Make is a interface for mannaging a token
type Maker interface {
	//CreateToken create a token for specifc user
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	//verify token
	VerifyToken(token string) (*Payload, error)
}
