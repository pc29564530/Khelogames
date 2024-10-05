package database

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// Store provides all functions to execute db
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore create a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) BeginTx(ctx *gin.Context) (*sql.Tx, error) {
	return s.db.Begin()
}
