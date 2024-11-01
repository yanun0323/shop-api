package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"main/internal/domain/entity"
	"main/internal/domain/repository"
	"main/internal/domain/usecase"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	emailFormatValidator       = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordUpperCastValidator = regexp.MustCompile(`.*[A-Z].*`)
	passwordLowerCastValidator = regexp.MustCompile(`.*[a-z].*`)
	passwordLengthValidator    = regexp.MustCompile(".*[()[\\]{}<>+\\-*/?,.:;\"'_\\\\|~`!@#$%^&=].*")
	passwordSpecialValidator   = regexp.MustCompile(`^.{6,16}$`)
)

type authUsecase struct {
	tokenUse usecase.TokenUsecase
	optUse   usecase.OTPUsecase
	userRepo repository.UserRepository
}

func NewAuthUsecase(
	tokenUse usecase.TokenUsecase,
	optUse usecase.OTPUsecase,
	userRepo repository.UserRepository,
) usecase.AuthUsecase {
	return &authUsecase{
		tokenUse: tokenUse,
		optUse:   optUse,
		userRepo: userRepo,
	}
}

func (use *authUsecase) Register(ctx context.Context, email, password, code string) error {
	if !use.verifyEmailFormat(email) {
		return errors.Errorf("mismatch email format (%s), err: %+v", email, usecase.ErrInvalidEmailFormat)
	}

	if !use.verifyPasswordFormat(password) {
		return errors.Errorf("mismatch password format (%s), err: %+v", password, usecase.ErrInvalidPasswordFormat)
	}

	pass, err := use.optUse.VerifyEmail(ctx, email, code)
	if err != nil {
		return errors.Errorf("verify email: %+v", err)
	}

	if !pass {
		return errors.Errorf("email (%s) verify code (%s) mismatch", email, code)
	}

	exist, err := use.userRepo.Exist(ctx, email)
	if err != nil {
		return errors.Errorf("check user exist: %+v", err)
	}

	if exist {
		return errors.Errorf("user email (%s) already exists", email)
	}

	user := &entity.User{
		Name:     "",
		Email:    email,
		Password: use.encryptedPassword(password),
	}

	if err := use.userRepo.Create(ctx, user); err != nil {
		return errors.Errorf("create user: %+v", err)
	}

	return nil
}

func (authUsecase) verifyEmailFormat(email string) bool {
	return emailFormatValidator.MatchString(email)
}

func (authUsecase) verifyPasswordFormat(password string) bool {
	return passwordUpperCastValidator.MatchString(password) &&
		passwordLowerCastValidator.MatchString(password) &&
		passwordLengthValidator.MatchString(password) &&
		passwordSpecialValidator.MatchString(password)
}

func (authUsecase) encryptedPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func (use *authUsecase) Login(ctx context.Context, param usecase.LoginParam) (*entity.AuthToken, error) {
	pass, err := use.optUse.VerifyEmail(ctx, param.Email, param.Code)
	if err != nil {
		return nil, errors.Errorf("verify email: %+v", err)
	}

	if !pass {
		return nil, errors.Errorf("email (%s) verify code (%s) mismatch", param.Email, param.Code)
	}

	user, err := use.userRepo.Get(ctx, repository.GetUserOption{
		Email: param.Email,
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, usecase.ErrInvalidEmail
		}

		return nil, errors.Errorf("get user: %+v", err)
	}

	password := use.encryptedPassword(param.Password)
	if !strings.EqualFold(password, user.Password) {
		return nil, usecase.ErrMismatchPassword
	}

	token, err := use.tokenUse.CreateToken(ctx, user.ID, param.DeviceID)
	if err != nil {
		return nil, errors.Errorf("create token: %+v", err)
	}

	return token, nil
}

func (use *authUsecase) SendVerifyCode(ctx context.Context, email string) error {
	if err := use.optUse.SendEmail(ctx, email); err != nil {
		return errors.Errorf("send email: %+v", err)
	}

	return nil
}

func (use *authUsecase) RefreshToken(ctx context.Context, userID int64, deviceID string) (*entity.AuthToken, error) {
	auth, err := use.tokenUse.RefreshToken(ctx, userID, deviceID)
	if err != nil {
		return nil, errors.Errorf("create token: %+v", err)
	}

	return auth, nil
}

func (use *authUsecase) Logout(ctx context.Context, userID int64, deviceID string) error {
	err := use.tokenUse.DeleteToken(ctx, userID, deviceID)
	if err != nil {
		return errors.Errorf("delete token: %+v", err)
	}

	return nil
}
