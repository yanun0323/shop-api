package repository

import (
	"context"
	"fmt"
	"main/config"
	"main/internal/domain/repository"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type otpRepository struct {
	rdb        *redis.Client
	expiration time.Duration
}

func NewOTPRepository(conf config.Config, rdb *redis.Client) repository.OTPRepository {
	return &otpRepository{
		rdb:        rdb,
		expiration: conf.OTP.Expiration,
	}
}

func (otpRepository) cacheKey(email string) string {
	return fmt.Sprintf("OTP:CODE:%s", email)
}

func (repo *otpRepository) Store(ctx context.Context, email string, code string) error {
	_, err := repo.rdb.Set(ctx, repo.cacheKey(email), code, repo.expiration).Result()
	if err != nil {
		return errors.Errorf("set otp, err: %+v", err)
	}

	return nil
}

func (repo *otpRepository) Get(ctx context.Context, email string) (string, error) {
	code, err := repo.rdb.Get(ctx, repo.cacheKey(email)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}

		return "", errors.Errorf("get otp, err: %+v", err)
	}

	return code, nil
}
