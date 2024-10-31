package usecase

import (
	"context"
	"main/config"
	"main/internal/domain/entity"
	"main/internal/domain/repository"
	"main/internal/domain/usecase"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type tokenUsecase struct {
	tokenSecret        []byte
	tokenAccessExpire  time.Duration
	tokenRefreshExpire time.Duration

	tokenRepo repository.TokenRepository
}

func NewTokenUsecase(conf config.Config, tokenRepo repository.TokenRepository) usecase.TokenUsecase {
	return &tokenUsecase{
		tokenSecret:        []byte(conf.Token.Secret),
		tokenAccessExpire:  conf.Token.Expiration.Access,
		tokenRefreshExpire: conf.Token.Expiration.Refresh,
		tokenRepo:          tokenRepo,
	}
}

func (use *tokenUsecase) VerifyToken(ctx context.Context, token string) (*entity.TokenClaims, error) {
	claim, err := use.parseToken(use.tokenSecret, token)
	if err != nil {
		return nil, errors.Errorf("verify token, err: %+v", err)
	}

	exist, err := use.tokenRepo.Exist(ctx, repository.TokenQuery{
		UserID:    claim.UserID,
		DeviceID:  claim.DeviceID,
		TokenType: claim.TokenType,
	})
	if err != nil {
		return nil, errors.Errorf("verify token, err: %+v", err)
	}

	if !exist {
		return nil, errors.Errorf("invalid token, err: %+v", usecase.ErrInvalidToken)
	}

	if err := claim.Valid(); err != nil {
		return nil, errors.Errorf("invalid token, err: %+v", err)
	}

	return claim, nil
}

func (use *tokenUsecase) RefreshToken(ctx context.Context, userID int64, deviceID string) (*entity.AuthToken, error) {
	refreshToken, err := use.tokenRepo.Get(ctx, repository.TokenQuery{
		UserID:    userID,
		DeviceID:  deviceID,
		TokenType: entity.TokenTypeRefreshToken,
	})
	if err != nil {
		return nil, errors.Errorf("get refresh token, err: %+v", err)
	}

	if len(refreshToken) == 0 {
		return nil, errors.Errorf("empty refresh token, err: %+v", usecase.ErrInvalidToken)
	}

	now := time.Now()
	accessToken, accessTokenExpiresAt, err := use.createAccessToken(ctx, now, userID, deviceID)
	if err != nil {
		return nil, errors.Errorf("create access token, err: %+v", err)
	}

	return &entity.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessTokenExpiresAt - time.Now().Unix(),
	}, nil
}

func (use *tokenUsecase) CreateToken(ctx context.Context, userID int64, deviceID string) (*entity.AuthToken, error) {
	now := time.Now()
	accessToken, accessTokenExpiresAt, err := use.createAccessToken(ctx, now, userID, deviceID)
	if err != nil {
		return nil, errors.Errorf("create access token, err: %+v", err)
	}

	refreshToken, err := use.createRefreshToken(ctx, now, userID, deviceID)
	if err != nil {
		return nil, errors.Errorf("create refresh token, err: %+v", err)
	}

	return &entity.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessTokenExpiresAt - time.Now().Unix(),
	}, nil
}

func (use *tokenUsecase) createAccessToken(ctx context.Context, now time.Time, userID int64, deviceID string) (accessToken string, expiredAt int64, err error) {
	accessTokenClaims := &entity.TokenClaims{
		UserID:    userID,
		DeviceID:  deviceID,
		TokenType: entity.TokenTypeAccessToken,
		ExpiresAt: now.Add(use.tokenAccessExpire).Unix(),
	}

	token, err := use.generateToken(use.tokenSecret, accessTokenClaims)
	if err != nil {
		return "", 0, errors.Errorf("create token, err: %+v", err)
	}

	if err := use.tokenRepo.Create(ctx, use.claimsToCreateTokenQuery(accessTokenClaims, accessToken)); err != nil {
		return "", 0, errors.Errorf("create access token, err: %+v", err)
	}

	return token, accessTokenClaims.ExpiresAt, nil
}

func (use *tokenUsecase) createRefreshToken(ctx context.Context, now time.Time, userID int64, deviceID string) (refreshToken string, err error) {
	refreshTokenClaims := &entity.TokenClaims{
		UserID:    userID,
		DeviceID:  deviceID,
		TokenType: entity.TokenTypeRefreshToken,
		ExpiresAt: now.Add(use.tokenRefreshExpire).Unix(),
	}

	token, err := use.generateToken(use.tokenSecret, refreshTokenClaims)
	if err != nil {
		return "", errors.Errorf("create token, err: %+v", err)
	}

	if err := use.tokenRepo.Create(ctx, use.claimsToCreateTokenQuery(refreshTokenClaims, token)); err != nil {
		return "", errors.Errorf("create refresh token, err: %+v", err)
	}

	return token, nil
}

func (tokenUsecase) claimsToCreateTokenQuery(claims *entity.TokenClaims, token string) repository.CreateTokenQuery {
	return repository.CreateTokenQuery{
		TokenQuery: repository.TokenQuery{
			UserID:    claims.UserID,
			DeviceID:  claims.DeviceID,
			TokenType: claims.TokenType,
		},
		Token:     token,
		ExpiredAt: claims.ExpiresAt,
	}
}

func (use *tokenUsecase) DeleteToken(ctx context.Context, userID int64, deviceID string) error {
	if err := use.tokenRepo.Delete(ctx, userID, deviceID); err != nil {
		return errors.Errorf("delete token, err: %+v", err)
	}

	return nil
}

func (tokenUsecase) parseToken(secret []byte, token string) (*entity.TokenClaims, error) {
	t, err := jwt.ParseWithClaims(token, &entity.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("invalid token type (%s), err: %+v", token.Method.Alg(), jwt.ErrInvalidKeyType)
		}

		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, errors.Errorf("invalid token type (%s), err: %+v", token.Method.Alg(), jwt.ErrInvalidKeyType)
		}

		return secret, nil
	})
	if err != nil {
		return nil, errors.Errorf("parse token, err: %+v", err)
	}

	if !t.Valid {
		return nil, errors.Errorf("invalid token, err: %+v", jwt.ErrSignatureInvalid)
	}

	claims, ok := t.Claims.(*entity.TokenClaims)
	if !ok {
		return nil, errors.Errorf("invalid token claims, err: %+v", jwt.ErrSignatureInvalid)
	}

	return claims, nil
}

func (tokenUsecase) generateToken(secret []byte, claims *entity.TokenClaims) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims = claims
	token, err := t.SignedString(secret)
	if err != nil {
		return "", errors.Errorf("generate token, err: %+v", err)
	}

	return token, nil
}
