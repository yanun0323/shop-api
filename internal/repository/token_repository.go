package repository

import (
	"context"
	"fmt"
	"main/internal/domain/entity"
	"main/internal/domain/repository"
	"main/internal/repository/conn"
	"main/internal/repository/query"
	"time"

	"github.com/pkg/errors"

	"github.com/redis/go-redis/v9"
)

type tokenRepository struct {
	db  *conn.Dao
	rdb *redis.Client
}

func NewTokenRepository(db *conn.Dao, rdb *redis.Client) repository.TokenRepository {
	return &tokenRepository{
		db:  db,
		rdb: rdb,
	}
}

func (tokenRepository) cacheKey(userID int64, deviceID string) string {
	return fmt.Sprintf("ACCESS:TOKEN:%d:%s", userID, deviceID)
}

func (repo *tokenRepository) Exist(ctx context.Context, q repository.TokenQuery) (bool, error) {
	if q.TokenType == entity.TokenTypeAccessToken {
		key := repo.cacheKey(q.UserID, q.DeviceID)
		result, err := repo.rdb.Exists(ctx, key).Result()
		if err != nil {
			if conn.IsNotFoundError(err) {
				return false, nil
			}

			return false, errors.Errorf("get redis key (%s), err: %+v", key, err)
		}

		return result != 0, nil
	}

	count, err := repo.db.CountToken(ctx, query.CountTokenParams{
		UserID:    q.UserID,
		DeviceID:  q.DeviceID,
		ExpiredAt: time.Now().Unix(),
	})
	if err != nil {
		if conn.IsNotFoundError(err) {
			return false, nil
		}

		return false, errors.Errorf("count token, err: %+v", err)
	}

	return count != 0, nil
}

func (repo *tokenRepository) Get(ctx context.Context, q repository.TokenQuery) (string, error) {
	if q.TokenType == entity.TokenTypeAccessToken {
		key := repo.cacheKey(q.UserID, q.DeviceID)
		result, err := repo.rdb.Get(ctx, key).Result()
		if err != nil {
			if conn.IsNotFoundError(err) {
				return "", errors.Errorf("redis get, err: %+v", repository.ErrNotFound)
			}

			return "", errors.Errorf("get redis key (%s), err: %+v", key, err)
		}

		return result, nil
	}

	token, err := repo.db.GetToken(ctx, query.GetTokenParams{
		UserID:    q.UserID,
		DeviceID:  q.DeviceID,
		ExpiredAt: time.Now().Unix(),
	})
	if err != nil {
		if conn.IsNotFoundError(err) {
			return "", errors.Errorf("get token, err: %+v", repository.ErrNotFound)
		}

		return "", errors.Errorf("get token, err: %+v", err)
	}

	return token.RefreshToken, nil

}

func (repo *tokenRepository) Create(ctx context.Context, q repository.CreateTokenQuery) error {
	if q.TokenType == entity.TokenTypeAccessToken {
		key := repo.cacheKey(q.UserID, q.DeviceID)
		expiration := time.Until(time.Unix(q.ExpiredAt, 0))

		_, err := repo.rdb.Set(ctx, key, q.Token, expiration).Result()
		if err != nil {
			return errors.Errorf("set token, err: %+v", err)
		}

		return nil
	}

	err := repo.db.CreateToken(ctx, query.CreateTokenParams{
		UserID:       q.UserID,
		DeviceID:     q.DeviceID,
		RefreshToken: q.Token,
		ExpiredAt:    q.ExpiredAt,
	})
	if err != nil {
		if conn.IsDuplicateKeyError(err) {
			return errors.Errorf("create token, err: %+v", repository.ErrDuplicateKey)
		}

		return errors.Errorf("create token, err: %+v", err)
	}

	return nil
}

func (repo *tokenRepository) Delete(ctx context.Context, userID int64, deviceID string) error {

	_, err := repo.rdb.Del(ctx, repo.cacheKey(userID, deviceID)).Result()
	if err != nil {
		if conn.IsNotFoundError(err) {
			return nil
		}

		return errors.Errorf("delete token, err: %+v", err)
	}

	if err := repo.db.DeleteToken(ctx, query.DeleteTokenParams{
		UserID:   userID,
		DeviceID: deviceID,
	}); err != nil {
		if conn.IsNotFoundError(err) {
			return nil
		}

		return errors.Errorf("delete token, err: %+v", err)
	}

	return nil
}
