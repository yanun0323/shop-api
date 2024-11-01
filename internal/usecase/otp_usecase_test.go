package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFillUpOTPCode(t *testing.T) {
	tests := []struct {
		name   string
		code   string
		length int
		want   string
	}{
		{
			name:   "same length",
			code:   "123456",
			length: 6,
			want:   "123456",
		},
		{
			name:   "larger length",
			code:   "123456",
			length: 10,
			want:   "0000123456",
		},
		{
			name:   "smaller length",
			code:   "12345678",
			length: 2,
			want:   "12345678",
		},
		{
			name:   "zero length",
			code:   "12345678",
			length: 0,
			want:   "12345678",
		},
		{
			name:   "negative length",
			code:   "12345678",
			length: -1,
			want:   "12345678",
		},
		{
			name:   "longer code",
			code:   "1234567890",
			length: 6,
			want:   "1234567890",
		},
		{
			name:   "shorter code",
			code:   "1",
			length: 6,
			want:   "000001",
		},
		{
			name:   "empty code",
			code:   "",
			length: 6,
			want:   "000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := new(otpUsecase).fillUpOTPCode(tt.code, tt.length)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGenerateOTPCode(t *testing.T) {
	testCases := []struct {
		desc   string
		length int
	}{
		{
			desc:   "generate otp code",
			length: 6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			use := new(otpUsecase)
			for i := 0; i < 999_999; i++ {
				code := use.generateOTPCode(tc.length)
				assert.LessOrEqual(t, len(code), tc.length)
			}
		})
	}
}
