package conn

import (
	"database/sql"
	"main/internal/repository/query"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Dao struct {
	*query.Queries

	db *sql.DB
}

func NewDao(db gorm.DB) (*Dao, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Errorf("get sql db, err: %+v", err)
	}

	return &Dao{
		Queries: query.New(sqlDB),
		db:      sqlDB,
	}, nil
}

// WithTx invokes the given query function within a transaction.
//
// If the query function returns an error, the transaction will be rolled back.
// Otherwise, the transaction will be committed.
func (dao *Dao) WithTx(fn func(q *query.Queries) error) error {
	tx, err := dao.db.Begin()
	if err != nil {
		return errors.Errorf("begin transaction, err: %+v", err)
	}
	defer tx.Rollback()

	if err := fn(dao.Queries); err != nil {
		return errors.Errorf("transaction, err: %+v", err)
	}

	if err := tx.Commit(); err != nil {
		return errors.Errorf("commit transaction, err: %+v", err)
	}

	return nil
}
