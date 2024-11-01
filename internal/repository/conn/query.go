package conn

import (
	"main/internal/repository/query"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func NewQueries(db *gorm.DB) (*query.Queries, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Errorf("get sql db, err: %+v", err)
	}

	return query.New(sqlDB), nil
}
