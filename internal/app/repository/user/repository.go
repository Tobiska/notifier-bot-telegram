package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"

	"notifier-bot-telegram/internal/app/models"
	sql2 "notifier-bot-telegram/pkg/storage/sql"
)

var (
	driver = goqu.Dialect("sqlite3")
)

type database interface {
	Begin() (*sql.Tx, error)
	WithTx(tx sql2.Transaction) *sql2.Storage
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type Repository struct {
	db database
}

func New(db database) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SaveUser(ctx context.Context, user models.User) error {
	query := driver.Insert(goqu.T("users")).Rows(
		goqu.Record{
			"chat_id":  user.ChatID,
			"user_id":  user.UserID,
			"status":   user.Status,
			"username": user.Username,
		},
	)

	sqlQuery, args, err := query.ToSQL()
	if err != nil {
		return fmt.Errorf("build sqlQuery query error: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, sqlQuery, args...); err != nil {
		return fmt.Errorf("execution error: %w", err)
	}
	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, chatID int64, status models.Status) error {
	query := driver.Update(goqu.T("users")).Where(goqu.Ex{"chat_id": chatID}).Set(
		goqu.Record{"status": status, "updated_at": time.Now().Local()},
	)

	sqlQuery, args, err := query.ToSQL()
	if err != nil {
		return fmt.Errorf("build sqlQuery query error: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, sqlQuery, args...); err != nil {
		return fmt.Errorf("execution error: %w", err)
	}
	return nil
}

func (r *Repository) FindUser(ctx context.Context, chatID int64) (*models.User, error) {
	query := driver.Select(
		goqu.C("chat_id"),
		goqu.C("user_id"),
		goqu.C("username"),
		goqu.C("status"),
		goqu.C("created_at"),
		goqu.C("updated_at"),
	).From(goqu.T("users")).Where(goqu.Ex{"chat_id": chatID})

	querySQL, args, err := query.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("build query error: %w", err)
	}

	u := &models.User{}

	if err := r.db.QueryRowContext(ctx, querySQL, args...).Scan(&u.ChatID, &u.UserID, &u.Username, &u.Status, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("execution error: %w", err)
	}
	return u, nil
}
