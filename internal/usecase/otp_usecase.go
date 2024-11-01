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
	"strings"

	"github.com/pkg/errors"
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
	if !emailFormatValidator.MatchString(email) {
		return errors.Errorf("invalid email format, err: %+v", usecase.ErrInvalidEmailFormat)
	}

	code := use.generateOTPCode(_verifyCodeCount)
	if err := use.otpRepo.Store(ctx, email, code); err != nil {
		return errors.Errorf("store otp: %+v", err)
	}

	if err := use.notifyRepo.SendEmail(email, entity.Notification{
		Type:    entity.NotifyTypeRegisterOtp,
		Subject: "OTP",
		Body:    fmt.Sprintf("Your OTP code is %s", use.fillUpOTPCode(code, _verifyCodeCount)),
	}); err != nil {
		return errors.Errorf("send email: %+v", err)
	}

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

	if otp != code {
		return false, nil
	}

	if err := use.otpRepo.Delete(ctx, email); err != nil {
		return false, errors.Errorf("delete otp: %+v", err)
	}

	return true, nil
}

func (otpUsecase) generateOTPCode(length int) string {
	limit := int(math.Pow10(length)) - 1
	return strconv.Itoa(rand.Intn(limit))
}

func (otpUsecase) fillUpOTPCode(code string, length int) string {
	if len(code) >= length {
		return code
	}

	return strings.Repeat("0", length-len(code)) + code
}
