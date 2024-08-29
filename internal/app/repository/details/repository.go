package details

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"

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

func (r *Repository) Begin() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *Repository) WithTx(tx *sql.Tx) *Repository {
	return &Repository{
		db: r.db.WithTx(tx),
	}
}

func (r *Repository) PartialUpdate(ctx context.Context, detailID int64, updateDetail UpdateModel) error {
	rec := goqu.Record{
		"updated_at": "NOW()",
	}

	if updateDetail.Name != nil {
		rec["name"] = *updateDetail.Name
	}

	if updateDetail.SoftDeadline != nil {
		rec["soft_deadline_at"] = *updateDetail.SoftDeadline
	}

	if updateDetail.HardDeadline != nil {
		rec["hard_deadline_at"] = *updateDetail.HardDeadline
	}

	query := driver.Update(goqu.T("details")).Set(rec).Where(goqu.Ex{"id": detailID})

	querySQL, args, err := query.ToSQL()
	if err != nil {
		return fmt.Errorf("build query error: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, querySQL, args...); err != nil {
		return fmt.Errorf("execute error: %w", err)
	}

	return nil
}

func (r *Repository) Search(ctx context.Context, filter Filter) ([]models.Detail, error) {
	query := driver.Select(
		goqu.C("id"),
		goqu.C("name"),
		goqu.C("chat_id"),
		goqu.C("soft_deadline_at"),
		goqu.C("hard_deadline_at"),
		goqu.C("created_at"),
		goqu.C("updated_at"),
	).From(goqu.T("details")).Prepared(true)

	if filter.Name != nil {
		query = query.Where(goqu.Ex{"name": *filter.Name})
	}

	if filter.ChatID != nil {
		query = query.Where(goqu.Ex{"chat_id": *filter.ChatID})
	}

	if filter.IsFilled != nil {
		if *filter.IsFilled {
			query = query.Where(goqu.C("chat_id").IsNotNull(),
				goqu.C("name").IsNotNull(),
				goqu.C("soft_deadline_at").IsNotNull(),
				goqu.C("hard_deadline_at").IsNotNull(),
			)
		} else {
			query = query.Where(goqu.Or(goqu.C("chat_id").IsNull(),
				goqu.C("name").IsNull(),
				goqu.C("soft_deadline_at").IsNull(),
				goqu.C("hard_deadline_at").IsNull(),
			))
		}
	}

	querySql, args, err := query.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("error while build sql: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, querySql, args...)
	if err != nil {
		return nil, fmt.Errorf("error exec: %w", err)
	}
	defer rows.Close()

	var details []models.Detail
	for rows.Next() {
		detail := models.Detail{}
		if err := rows.Scan(&detail.ID, &detail.ChatID, &detail.Name, &detail.SoftDeadline, &detail.HardDeadline); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		details = append(details, detail)
	}

	return details, nil
}

func (r *Repository) Save(ctx context.Context, detail models.Detail) error {
	rec := goqu.Record{
		"chat_id":          detail.ChatID,
		"name":             detail.Name,
		"soft_deadline_at": detail.SoftDeadline,
		"hard_deadline_at": detail.HardDeadline,
		"updated_at":       "NOW()",
	}

	query := driver.Insert(goqu.T("details")).Rows(rec).Prepared(true)

	querySQL, args, err := query.ToSQL()
	if err != nil {
		return fmt.Errorf("build query error: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, querySQL, args...); err != nil {
		return fmt.Errorf("execute error: %w", err)
	}
	return nil
}
