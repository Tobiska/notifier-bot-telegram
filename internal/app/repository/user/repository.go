package user

import (
	"context"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"os/user"
)

var (
	driver = goqu.Dialect("sqlite3")
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SaveUser(ctx context.Context, user user.User) error {
	driver.Select(goqu.L(""))
	return nil
}

func (r *Repository) FindUser(ctx context.Context, chatID int64) (*user.User, error) {
	return nil, nil
}
