package details

import (
	"context"
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"notifier-bot-telegram/internal/app/models"
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

func (r *Repository) Upsert(ctx context.Context, chatID int64, detail models.Detail) error {

}
