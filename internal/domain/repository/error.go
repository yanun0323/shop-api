package repository

import "github.com/pkg/errors"

var (
	ErrNotFound     = errors.New("data not found")
	ErrDuplicateKey = errors.New("duplicate key")
)
