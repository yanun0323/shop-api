package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailFormatValidator(t *testing.T) {
	testCases := []struct {
		desc  string
		email string
		want  bool
	}{
		{
			desc:  "pass domain",
			email: "test@test-123.com",
			want:  true,
		},
		{
			desc:  "pass username",
			email: "test.user-name+tag@test.com",
			want:  true,
		},
		{
			desc:  "empty username",
			email: "@test.com",
			want:  false,
		},
		{
			desc:  "empty domain",
			email: "test.user-name+tag",
			want:  false,
		},
		{
			desc:  "short domain",
			email: "test@test",
			want:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := emailFormatValidator.MatchString(tc.email)
			assert.Equal(t, tc.want, result)
		})
	}
}

func TestVerifyPasswordFormat(t *testing.T) {
	testCases := []struct {
		desc string
		pass string
		want bool
	}{
		{
			desc: "pass",
			pass: "Password@123",
			want: true,
		},
		{
			desc: "pass 6 length",
			pass: "Pass@6",
			want: true,
		},
		{
			desc: "pass 16 length",
			pass: "Password@1234567",
			want: true,
		},
		{
			desc: "no upper case",
			pass: "password@123",
			want: false,
		},
		{
			desc: "no lower case",
			pass: "PASSWORD@123",
			want: false,
		},
		{
			desc: "too short, 5 length",
			pass: "5h@rt",
			want: false,
		},
		{
			desc: "no special character",
			pass: "noSpecial123",
			want: false,
		},
		{
			desc: "too long, 17 length",
			pass: "Password@12345678",
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := new(authUsecase).verifyPasswordFormat(tc.pass)
			assert.Equal(t, tc.want, result)
		})
	}
}

func TestEncryptedPassword(t *testing.T) {
	testCases := []struct {
		desc string
		pass string
		want string
	}{
		{
			"pass",
			"Password@123",
			"ff7bd97b1a7789ddd2775122fd6817f3173672da9f802ceec57f284325bf589f",
		},
		{
			"pass",
			"pASSword@123",
			"7a6e4766b5532d723b44cd712a522573d68ca8cd8e508c3832fac58041213daa",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := new(authUsecase).encryptedPassword(tc.pass)
			assert.Equal(t, tc.want, result)
		})
	}
}
