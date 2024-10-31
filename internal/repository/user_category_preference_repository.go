package repository

import (
	"context"

	"fmt"
	"main/config"
	"main/internal/domain/entity"
	"main/internal/domain/repository"
	"main/internal/repository/model"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userCategoryPreferenceRepository struct {
	db  *gorm.DB
	rdb *redis.Client

	recommendationExpiration time.Duration
}

func NewUserCategoryPreferenceRepository(conf config.Config, db *gorm.DB, rdb *redis.Client) repository.UserCategoryPreferenceRepository {
	return &userCategoryPreferenceRepository{
		db:                       db,
		rdb:                      rdb,
		recommendationExpiration: conf.Product.Recommendation.Expiration,
	}
}

func (userCategoryPreferenceRepository) cacheKey(userID int64) string {
	return fmt.Sprintf("USER:CATEGORY:PREFERENCE:ID:%d", userID)
}

func (repo *userCategoryPreferenceRepository) Get(ctx context.Context, userID int64) (entity.ProductCategory, error) {
	key := repo.cacheKey(userID)
	result, err := repo.rdb.Get(ctx, key).Int64()
	if err == nil {
		return entity.ProductCategory(result), nil
	}

	if !errors.Is(err, redis.Nil) {
		return 0, errors.Errorf("get redis key (%s), err: %+v", key, err)
	}

	categoryID, err := repo.getFromMySQL(ctx, userID)
	if err != nil {
		return 0, errors.Errorf("get from mysql, err: %+v", err)
	}

	err = repo.rdb.Set(ctx, key, categoryID, repo.recommendationExpiration).Err()
	if err != nil {
		return 0, errors.Errorf("set redis key (%s), err: %+v", key, err)
	}

	return categoryID, nil
}

func (repo *userCategoryPreferenceRepository) getFromMySQL(ctx context.Context, userID int64) (entity.ProductCategory, error) {
	var ucp *model.UserCategoryPreference

	if err := repo.db.WithContext(ctx).Table(model.UserCategoryPreference{}.TableName()).
		Where("user_id = ?", userID).
		Take(&ucp).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1, nil
		}

		return 0, errors.Errorf("get user category preference, err: %+v", err)
	}

	return ucp.CategoryID, nil
}
