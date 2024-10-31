package usecase

import (
	"context"
	"fmt"
	"main/internal/domain/entity"
	"main/internal/domain/repository"
	"main/internal/domain/usecase"
	"math"
	"math/rand"
	"strconv"

	"github.com/pkg/errors"
	"github.com/yanun0323/pkg/logs"
)

const (
	_verifyCodeCount = 6
)

type otpUsecase struct {
	otpRepo    repository.OTPRepository
	notifyRepo repository.NotificationRepository
}

func NewOTPUsecase(otpRepo repository.OTPRepository, notifyRepo repository.NotificationRepository) usecase.OTPUsecase {
	return &otpUsecase{
		otpRepo:    otpRepo,
		notifyRepo: notifyRepo,
	}
}

func (use *otpUsecase) SendEmail(ctx context.Context, email string) error {
	code := use.generateOTPCode()
	if err := use.otpRepo.Store(ctx, email, code); err != nil {
		return errors.Errorf("store otp: %+v", err)
	}

	if err := use.notifyRepo.SendEmail(email, entity.Notification{
		Type:    entity.NotifyTypeRegisterOtp,
		Subject: "OTP",
		Body:    fmt.Sprintf("Your OTP code is %s", code),
	}); err != nil {
		return errors.Errorf("send email: %+v", err)
	}

	logs.Debugf("Send OTP code: %s", code)

	return nil
}

func (use *otpUsecase) VerifyEmail(ctx context.Context, email, code string) (bool, error) {
	if len(code) > _verifyCodeCount {
		return false, errors.Errorf("empty code")
	}

	otp, err := use.otpRepo.Get(ctx, email)
	if err != nil {
		return false, errors.Errorf("get otp: %+v", err)
	}

	logs.Debugf("Verify OTP code: %s", otp)

	return otp == code, nil
}

func (otpUsecase) generateOTPCode() string {
	limit := int(math.Pow10(_verifyCodeCount)) - 1
	return strconv.Itoa(rand.Intn(limit))
}
