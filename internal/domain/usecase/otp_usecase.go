package usecase

import "context"

//go:generate domaingen -destination=../../usecase/otp_usecase.go -package=usecase -constructor
type OTPUsecase interface {
	SendEmail(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, email, code string) (bool, error)
}
