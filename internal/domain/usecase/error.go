package usecase

import (
	"fmt"

	"github.com/pkg/errors"
)

var _ error = (*UsecaseError)(nil) /* Check Interface Implementation */

type UsecaseError struct {
	Code int
	err  error
}

func (e *UsecaseError) Error() string {
	return fmt.Sprintf("code: %d, err: %+v", e.Code, e.err)
}

func ue(code int, text string) *UsecaseError {
	return &UsecaseError{
		Code: code,
		err:  errors.New(text),
	}
}

var (
	ErrNil                   = UsecaseError{200, nil}
	ErrInternal              = ue(999999, "internal error")
	ErrNotFound              = ue(000001, "not found")
	ErrInvalidEmailFormat    = ue(000002, "invalid email format")
	ErrInvalidPasswordFormat = ue(000003, "invalid password format")
	ErrInvalidEmail          = ue(000004, "invalid email")
	ErrMismatchPassword      = ue(000005, "mismatch password")
	ErrTokenExpired          = ue(000006, "token expired")
	ErrInvalidToken          = ue(000007, "invalid token")
)
