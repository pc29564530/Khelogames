package transactions

import (
	"context"
	"database/sql"
	"fmt"
	"khelogames/api/shared"
	coreToken "khelogames/core/token"
	"khelogames/database"
	"khelogames/logger"
)

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db         *sql.DB
	tokenMaker *coreToken.Maker
	logger     *logger.Logger
	*database.Queries
	scoreBroadcaster shared.ScoreBroadcaster
}

// NewStore creates a new store
func NewStore(db *sql.DB, tokenMaker *coreToken.Maker, logger *logger.Logger, scoreBroadcaster shared.ScoreBroadcaster) *SQLStore {
	return &SQLStore{
		db:               db,
		tokenMaker:       tokenMaker,
		logger:           logger,
		Queries:          database.New(db),
		scoreBroadcaster: scoreBroadcaster,
	}
}

// BeginTx starts a database transaction
func (store *SQLStore) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return store.db.BeginTx(ctx, nil)
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*database.Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := database.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (s *SQLStore) SetScoreBroadcaster(broadcaster shared.ScoreBroadcaster) {
	s.scoreBroadcaster = broadcaster
}

// Ensure CricketServer implements the shared interfaces
// var _ shared.CricketScoreUpdater

var _ shared.ScoreBroadcaster
