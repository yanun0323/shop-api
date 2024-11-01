package conn

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	_mysqlDuplicateKeyErrorCode = 1062
)

func IsDuplicateKeyError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	var mysqlError *mysql.MySQLError
	if errors.As(err, &mysqlError) {
		return mysqlError.Number == _mysqlDuplicateKeyErrorCode
	}
	return false
}

func IsNotFoundError(err error) bool {
	if errors.Is(err, redis.Nil) {
		return true
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}

	// var mysqlError *mysql.MySQLError
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}
	return false
}
