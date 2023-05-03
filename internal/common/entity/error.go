package entity

import (
	"errors"
)

var (
	ErrNoRows     = errors.New("error no rows")
	ErrDuplicated = errors.New("error duplicated")
)
