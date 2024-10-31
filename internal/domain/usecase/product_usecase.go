package usecase

import (
	"context"
	"main/internal/domain/entity"
	"main/internal/helper/pager"
)

//go:generate domaingen -destination=../../usecase/product_usecase.go -package=usecase -constructor
type ProductUsecase interface {
	ListUserRecommended(ctx context.Context, userID int64, pagination pager.Request) ([]*entity.Product, *pager.Response, error)
}

type GetProductRequest struct {
	ID   int64
	Name string
}

type ListProductRequest struct {
	ID   int64
	Name string
	Rank int64
}
