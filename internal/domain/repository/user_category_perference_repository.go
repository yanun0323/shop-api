package repository

import (
	"context"
	"main/internal/domain/entity"
)

//go:generate domaingen -destination=../../repository/user_category_preference_repository.go -package=repository -constructor
type UserCategoryPreferenceRepository interface {
	Get(ctx context.Context, userID int64) (entity.ProductCategory, error)
}
