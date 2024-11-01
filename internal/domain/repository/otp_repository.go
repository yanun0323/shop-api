package repository

import "context"

//go:generate domaingen -destination=../../repository/otp_repository.go -package=repository -constructor
type OTPRepository interface {
	Store(ctx context.Context, email string, code string) error
	Get(ctx context.Context, email string) (string, error)
	Delete(ctx context.Context, email string) error
}
