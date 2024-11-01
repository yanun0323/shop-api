package conn

import (
	"context"
	"fmt"
	"main/config"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormDB(ctx context.Context, conf config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.MySQL.User,
		conf.MySQL.Password,
		conf.MySQL.Host,
		conf.MySQL.Port,
		conf.MySQL.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Errorf("connect gorm, err: %+v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Errorf("get sql db, err: %+v", err)
	}

	sqlDB.SetMaxIdleConns(conf.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MySQL.MaxOpenConns)

	return db, nil
}
