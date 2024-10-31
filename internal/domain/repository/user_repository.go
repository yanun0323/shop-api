package repository

import (
	"context"
	"main/internal/domain/entity"
)

//go:generate domaingen -destination=../../repository/user_repository.go -package=repository -constructor
type UserRepository interface {
	Exist(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user *entity.User) error
	Get(ctx context.Context, opt GetUserOption) (*entity.User, error)
	Delete(ctx context.Context, opt GetUserOption) error
}

type GetUserOption struct {
	ID    int64
	Email string
}
