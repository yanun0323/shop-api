package repository

import (
	"context"
	"main/internal/domain/entity"
)

//go:generate domaingen -destination=../../repository/token_repository.go -package=repository -constructor
type TokenRepository interface {
	Exist(ctx context.Context, query TokenQuery) (bool, error)
	Get(ctx context.Context, query TokenQuery) (string, error)
	Create(ctx context.Context, query CreateTokenQuery) error
	Delete(ctx context.Context, userID int64, deviceID string) error
}

type TokenQuery struct {
	UserID    int64
	DeviceID  string
	TokenType entity.TokenType
}

type CreateTokenQuery struct {
	TokenQuery

	Token     string
	ExpiredAt int64
}
