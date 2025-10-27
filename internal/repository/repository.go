package repository

import (
	"database/sql"
	"errors"
	"marketplace/internal/errs"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}
func (r *Repository) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}
func (r *Repository) translateError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errs.ErrNotfound
	default:
		return err
	}
}
