package entity

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrTokenClaimsUnknownUserID   = errors.Errorf("unknown user id")
	ErrTokenClaimsUnknownDeviceID = errors.Errorf("unknown device id")
	ErrTokenClaimsUnknownType     = errors.Errorf("unknown token type")
	ErrTokenClaimsExpired         = errors.Errorf("token expired")
)

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` /* token expiration time in seconds */
}

type TokenType string

const (
	TokenTypeAccessToken  TokenType = "accessToken"
	TokenTypeRefreshToken TokenType = "refreshToken"
)

type TokenClaims struct {
	UserID    int64     `json:"user_id"`
	DeviceID  string    `json:"device_id"`
	TokenType TokenType `json:"token_type"`
	ExpiresAt int64     `json:"expires_at"`
}

func (c *TokenClaims) Valid() error {
	if c.UserID == 0 {
		return ErrTokenClaimsUnknownUserID
	}

	if len(c.DeviceID) == 0 {
		return ErrTokenClaimsUnknownDeviceID
	}

	switch c.TokenType {
	case TokenTypeAccessToken, TokenTypeRefreshToken:
	default:
		return ErrTokenClaimsUnknownType
	}

	if c.ExpiresAt <= time.Now().Unix() {
		return ErrTokenClaimsExpired
	}

	return nil
}
