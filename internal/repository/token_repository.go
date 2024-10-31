package repository

import (
	"context"
	"fmt"
	"main/internal/domain/entity"
	"main/internal/domain/repository"
	"main/internal/repository/model"
	"time"

	"github.com/pkg/errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tokenRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewTokenRepository(db *gorm.DB, rdb *redis.Client) repository.TokenRepository {
	return &tokenRepository{
		db:  db,
		rdb: rdb,
	}
}

func (tokenRepository) cacheKey(userID int64, deviceID string) string {
	return fmt.Sprintf("ACCESS:TOKEN:%d:%s", userID, deviceID)
}

func (repo *tokenRepository) Exist(ctx context.Context, query repository.TokenQuery) (bool, error) {
	if query.TokenType == entity.TokenTypeAccessToken {
		key := repo.cacheKey(query.UserID, query.DeviceID)
		result, err := repo.rdb.Exists(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return false, nil
			}

			return false, errors.Errorf("get redis key (%s), err: %+v", key, err)
		}

		return result != 0, nil
	}

	var count int64
	err := repo.db.WithContext(ctx).Table(model.Token{}.TableName()).
		Where("user_id = ? AND device_id = ? AND expired_at > ?", query.UserID, query.DeviceID, time.Now().Unix()).
		Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, errors.Errorf("count token, err: %+v", err)
	}

	return count != 0, nil
}

func (repo *tokenRepository) Get(ctx context.Context, query repository.TokenQuery) (string, error) {
	if query.TokenType == entity.TokenTypeAccessToken {
		key := repo.cacheKey(query.UserID, query.DeviceID)
		result, err := repo.rdb.Get(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return "", errors.Errorf("redis get, err: %+v", repository.ErrNotFound)
			}

			return "", errors.Errorf("get redis key (%s), err: %+v", key, err)
		}

		return result, nil
	}

	var token model.Token
	err := repo.db.WithContext(ctx).Table(token.TableName()).
		Where("user_id = ? AND device_id = ? AND expired_at > ?", query.UserID, query.DeviceID, time.Now().Unix()).
		Take(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.Errorf("get token, err: %+v", repository.ErrNotFound)
		}

		return "", errors.Errorf("get token, err: %+v", err)
	}

	return token.RefreshToken, nil

}

func (repo *tokenRepository) Create(ctx context.Context, query repository.CreateTokenQuery) error {
	if query.TokenType == entity.TokenTypeAccessToken {
		key := repo.cacheKey(query.UserID, query.DeviceID)
		expiration := time.Until(time.Unix(query.ExpiredAt, 0))

		_, err := repo.rdb.Set(ctx, key, query.Token, expiration).Result()
		if err != nil {
			return errors.Errorf("set token, err: %+v", err)
		}

		return nil
	}

	token := model.Token{
		UserID:       query.UserID,
		DeviceID:     query.DeviceID,
		RefreshToken: query.Token,
		ExpiredAt:    query.ExpiredAt,
	}
	err := repo.db.WithContext(ctx).Table(model.Token{}.TableName()).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "device_id"}},
			UpdateAll: true,
		}).Create(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.Errorf("create token, err: %+v", repository.ErrDuplicateKey)
		}

		return errors.Errorf("create token, err: %+v", err)
	}

	return nil
}

func (repo *tokenRepository) Delete(ctx context.Context, userID int64, deviceID string) error {

	_, err := repo.rdb.Del(ctx, repo.cacheKey(userID, deviceID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}

		return errors.Errorf("delete token, err: %+v", err)
	}

	err = repo.db.WithContext(ctx).Table(model.Token{}.TableName()).
		Where("user_id = ? AND device_id = ?", userID, deviceID).
		Delete(&model.Token{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return errors.Errorf("delete token, err: %+v", err)
	}

	return nil
}
