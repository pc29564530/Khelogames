package transcation_setup

import (
	"context"
	"database/sql"
	"khelogames/database/models"
	"khelogames/logger"
	"khelogames/token"
	"khelogames/util"

	"github.com/google/uuid"
)

type AuthServer struct {
	store      *Store
	logger     *logger.Logger
	tokenMaker token.Maker
	config     util.Config
}

// DBTX represents a database transaction interface
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// Querier defines all database query methods
type Querier interface {
	CreateEmailSignUp(ctx context.Context, fullName, username, email, hashPassword string) (*models.Users, error)
	CreateProfile(ctx context.Context, arg CreateProfileParams) (*models.UserProfiles, error)
	CreateNewToken(ctx context.Context, userPublicID uuid.UUID, userID int32, s *AuthServer) (map[string]interface{}, error)
}

// Store defines the interface for database operations including transactions
type Store interface {
	Querier
	BeginTx(ctx context.Context) (*sql.Tx, error)

	// Transaction methods
	CreateEmailSignUpTx(ctx context.Context, fullName, username, email, hashPassword string) (CreateEmailSignUpTxResult, error)
}
