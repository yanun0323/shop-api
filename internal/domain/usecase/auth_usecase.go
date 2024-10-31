package usecase

import (
	"context"
	"main/internal/domain/entity"
)

//go:generate domaingen -destination=../../usecase/auth_usecase.go -package=usecase -constructor
type AuthUsecase interface {
	Register(ctx context.Context, email, password, code string) error
	Login(ctx context.Context, param LoginParam) (*entity.AuthToken, error)
	SendVerifyCode(ctx context.Context, email string) error
	RefreshToken(ctx context.Context, userID int64, deviceID string) (*entity.AuthToken, error)
	Logout(ctx context.Context, userID int64, deviceID string) error
}

type LoginParam struct {
	Email    string
	Password string
	Code     string
	DeviceID string
}
