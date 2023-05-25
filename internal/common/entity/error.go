package entity

import (
	"errors"
)

var (
	ErrAccountNotFound       = errors.New("account not found")
	ErrDuplicated            = errors.New("error duplicated")
	ErrEmailPasswordNotMatch = errors.New("email password not match")
	ErrNoRows                = errors.New("error no rows")
)
