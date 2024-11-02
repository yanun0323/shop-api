package repository

import (
	"context"
	"fmt"
	"main/config"
	"main/internal/domain/entity"
	"main/internal/domain/repository"
	"main/internal/repository/conn"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/redis/go-redis/v9"
)

type productRepository struct {
	db  *conn.Dao
	rdb *redis.Client

	recommendationExpiration       time.Duration
	categoryProductCacheMutexTable sync.Map
}

func NewProductRepository(conf config.Config, db *conn.Dao, rdb *redis.Client) repository.ProductRepository {
	return &productRepository{
		db:                             db,
		rdb:                            rdb,
		recommendationExpiration:       conf.Product.Recommendation.Expiration,
		categoryProductCacheMutexTable: sync.Map{},
	}
}

func (repo *productRepository) GetUserRanked(ctx context.Context, query repository.GetUserRankedQuery) ([]*entity.Product, int64, error) {
	result, count, err := repo.getCategoryProductList(ctx, query.Category, query.Offset, query.Limit)
	if err != nil {
		return nil, 0, errors.Errorf("get product list, err: %+v", err)
	}

	return result, count, nil
}

func (repo *productRepository) getCategoryProductList(ctx context.Context, category entity.ProductCategory, offset, limit int) ([]*entity.Product, int64, error) {
	start, end := int64(offset), int64(offset+limit-1)
	result, count, err := repo.loadCategoryProductListFromCache(ctx, category, start, end)
	if err == nil {
		return result, count, nil
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return nil, 0, errors.Errorf("load product list from cache, err: %+v", err)
	}

	mu, _ := repo.categoryProductCacheMutexTable.LoadOrStore(category, &sync.Mutex{})
	lock := mu.(*sync.Mutex)
	lock.Lock()
	defer lock.Unlock()

	// reload list from cache again after get lock

	result, count, err = repo.loadCategoryProductListFromCache(ctx, category, start, end)
	if err == nil {
		return result, count, nil
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return nil, 0, errors.Errorf("load product list from cache, err: %+v", err)
	}

	result, err = repo.getCategoryProductListFromMySQL(ctx, category)
	if err != nil {
		return nil, 0, errors.Errorf("get product list from mysql, err: %+v", err)
	}

	err = repo.cacheCategoryProductList(ctx, category, result)
	if err != nil {
		return nil, 0, errors.Errorf("cache product list, err: %+v", err)
	}

	count = int64(len(result))
	if offset < len(result) {
		result = result[offset:]
	}

	if limit < len(result) {
		result = result[:limit]
	}

	return result, count, nil
}

func (*productRepository) categoryProductListCacheKey(category entity.ProductCategory) string {
	return fmt.Sprintf("PRODUCT:CATEGORY:%d:LIST", category)
}

func (repo *productRepository) loadCategoryProductListFromCache(ctx context.Context, category entity.ProductCategory, start, end int64) ([]*entity.Product, int64, error) {
	key := repo.categoryProductListCacheKey(category)
	count, err := repo.rdb.LLen(ctx, key).Result()
	if err != nil {
		if conn.IsNotFoundError(err) {
			return []*entity.Product{}, 0, repository.ErrNotFound
		}

		return nil, 0, errors.Errorf("get product list from cache, err: %+v", err)
	}

	var products []*entity.Product
	if err := repo.rdb.LRange(ctx, key, start, end).ScanSlice(&products); err != nil {
		if conn.IsNotFoundError(err) {
			return []*entity.Product{}, 0, repository.ErrNotFound
		}

		return nil, 0, errors.Errorf("get product list from cache, err: %+v", err)
	}

	if count == 0 {
		return nil, 0, repository.ErrNotFound
	}

	return products, count, nil
}

func (repo *productRepository) cacheCategoryProductList(ctx context.Context, category entity.ProductCategory, products []*entity.Product) error {
	key := repo.categoryProductListCacheKey(category)

	err := repo.rdb.Watch(ctx, func(tx *redis.Tx) error {
		for _, p := range products {
			_, err := tx.RPush(ctx, key, p).Result()
			if err != nil {
				return errors.Errorf("set product list to cache, err: %+v", err)
			}
		}

		_, err := tx.Expire(ctx, key, repo.recommendationExpiration).Result()
		if err != nil {
			return errors.Errorf("set product list to cache, err: %+v", err)
		}

		return nil
	})
	if err != nil {
		return errors.Errorf("set product list to cache, err: %+v", err)
	}

	return nil
}

func (repo *productRepository) getCategoryProductListFromMySQL(ctx context.Context, category entity.ProductCategory) ([]*entity.Product, error) {
	products, err := repo.db.ListProducts(ctx, int64(category))
	if err != nil {
		if conn.IsNotFoundError(err) {
			return []*entity.Product{}, nil
		}

		return nil, errors.Errorf("get product, err: %+v", err)
	}

	result := make([]*entity.Product, 0, len(products))
	for _, p := range products {
		result = append(result, &entity.Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			CategoryID:  entity.ProductCategory(p.CategoryID),
			Price:       p.Price,
			Rank:        p.Rank,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt.ValueOrZero(),
		})
	}

	return result, nil
}
