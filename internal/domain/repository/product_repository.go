package repository

import (
	"context"
	"main/internal/domain/entity"
)

//go:generate domaingen -destination=../../repository/product_repository.go -package=repository -constructor
type ProductRepository interface {
	GetUserRanked(ctx context.Context, query GetUserRankedQuery) ([]*entity.Product, int64, error)
}

type GetUserRankedQuery struct {
	UserID   int64
	Category entity.ProductCategory
	Offset   int
	Limit    int
}
