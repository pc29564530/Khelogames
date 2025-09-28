package transaction

import (
	"context"
	"database/sql"
	"khelogames/internal/domain/repository"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*repository.Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: repository.New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *Store) ExecTx(ctx context.Context, fn func(*repository.Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := repository.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}
