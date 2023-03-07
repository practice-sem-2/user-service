package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db    *sqlx.DB
	scope Scope
	UserStorage
}

type Scope interface {
	sqlx.QueryerContext
	sqlx.ExecerContext
	sqlx.Execer
	sqlx.Queryer
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db:          db,
		scope:       db,
		UserStorage: NewUserStorage(db),
	}
}

func (s *Storage) Atomic(ctx context.Context, fn func(store *Storage) error) (err error) {
	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback caused by error: \"%v\" failed: %v", err, rbErr)
			}
		} else {
			err = tx.Commit()
		}
	}()

	storage := Storage{s.db, tx, s.UserStorage}
	err = fn(&storage)
	return nil
}
