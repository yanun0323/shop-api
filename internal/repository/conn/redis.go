package conn

import (
	"context"
	"main/config"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, conf config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:           conf.Redis.Addr,
		DB:             conf.Redis.Database,
		Username:       conf.Redis.Username,
		Password:       conf.Redis.Password,
		MaxIdleConns:   conf.Redis.MaxIdleConns,
		MaxActiveConns: conf.Redis.MaxOpenConns,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, errors.Errorf("redis ping err: %+v", err)
	}

	return rdb, nil
}
