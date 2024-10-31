package usecase

import (
	"context"
	"main/internal/domain/entity"
	"main/internal/domain/repository"
	"main/internal/domain/usecase"
	"main/internal/helper/pager"

	"github.com/pkg/errors"
)

type productUsecase struct {
	productRepo            repository.ProductRepository
	userCategoryPreferRepo repository.UserCategoryPreferenceRepository
}

func NewProductUsecase(
	productRepo repository.ProductRepository,
	userCategoryPreferRepo repository.UserCategoryPreferenceRepository,
) usecase.ProductUsecase {
	return &productUsecase{
		productRepo:            productRepo,
		userCategoryPreferRepo: userCategoryPreferRepo,
	}
}

func (use *productUsecase) ListUserRecommended(ctx context.Context, userID int64, pagination pager.Request) ([]*entity.Product, *pager.Response, error) {
	categoryID, err := use.userCategoryPreferRepo.Get(ctx, userID)
	if err != nil {
		return nil, nil, errors.Errorf("get user category preference, err: %+v", err)
	}

	products, count, err := use.productRepo.GetUserRanked(ctx, repository.GetUserRankedQuery{
		UserID:   userID,
		Category: categoryID,
		Offset:   pagination.Offset(),
		Limit:    pagination.Limit(),
	})
	if err != nil {
		return nil, nil, errors.Errorf("get user ranked, err: %+v", err)
	}

	return products, pagination.Response(count), nil
}
