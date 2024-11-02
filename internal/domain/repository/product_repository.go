package repository

import (
	"context"
	"main/internal/domain/entity"
)

//go:generate domaingen -destination=../../repository/product_repository.go -package=repository -constructor
type ProductRepository interface {
	ListRankedProductsByCategory(ctx context.Context, query ListRankedProductsByCategoryQuery) ([]*entity.Product, int64, error)
}

type ListRankedProductsByCategoryQuery struct {
	Category entity.ProductCategory
	Offset   int
	Limit    int
}
