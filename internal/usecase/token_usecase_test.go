package usecase

import (
	"main/internal/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseToken(t *testing.T) {
	testCases := []struct {
		desc   string
		claims entity.TokenClaims
		secret []byte
	}{
		{
			desc: "generate and parse token",
			claims: entity.TokenClaims{
				UserID:    1,
				DeviceID:  "device-id",
				TokenType: entity.TokenTypeAccessToken,
				ExpiresAt: 9999999999,
			},
			secret: []byte("secret"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			use := new(tokenUsecase)
			token, err := use.generateToken(tc.secret, &tc.claims)
			require.NoError(t, err)

			claims, err := use.parseToken(tc.secret, token)
			require.NoError(t, err)

			assert.EqualExportedValues(t, tc.claims, *claims)
		})
	}
}
