package sql

import (
	"context"
	"database/sql"
)

type querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type txCreator interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type txer interface {
	Commit() error
	Rollback() error
}

type Database interface {
	querier
	txCreator
}

type Transaction interface {
	querier
	txer
}

type Storage struct {
	db Database
	tx Transaction
}

func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) Begin() (*sql.Tx, error) {
	if s.tx != nil {
		panic("in repository tx already exist")
	}
	return s.db.Begin()
}

func (s *Storage) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	if s.tx != nil {
		panic("in repository tx already exist")
	}
	return s.db.BeginTx(ctx, opts)
}

func (s *Storage) WithTx(tx Transaction) *Storage {
	return &Storage{
		db: s.db,
		tx: tx,
	}
}

func (s *Storage) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.getQuerier().ExecContext(ctx, query, args...)
}
func (s *Storage) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return s.getQuerier().PrepareContext(ctx, query)
}
func (s *Storage) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.getQuerier().QueryContext(ctx, query, args...)
}
func (s *Storage) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.getQuerier().QueryRowContext(ctx, query, args...)
}

func (s *Storage) Rollback() error {
	if s.tx == nil {
		panic("this repository has no transaction")
	}

	return s.tx.Rollback()
}

func (s *Storage) Commit() error {
	if s.tx == nil {
		panic("this repository has no transaction")
	}

	return s.tx.Commit()
}

func (s *Storage) getQuerier() querier {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}
