package usecase

import (
	"context"
	"main/internal/domain/entity"
)

//go:generate domaingen -destination=../../usecase/token_usecase.go -package=usecase -constructor
type TokenUsecase interface {
	VerifyToken(ctx context.Context, token string) (*entity.TokenClaims, error)
	RefreshToken(ctx context.Context, userID int64, deviceID string) (*entity.AuthToken, error)
	CreateToken(ctx context.Context, userID int64, deviceID string) (*entity.AuthToken, error)
	DeleteToken(ctx context.Context, userID int64, deviceID string) error
}
