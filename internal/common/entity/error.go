package entity

import (
	"errors"
)

var (
	ErrAccountNotFound       = errors.New("account not found")
	ErrDuplicated            = errors.New("error duplicated")
	ErrEmailPasswordNotMatch = errors.New("email password not match")
	ErrGrantTypeInvalid      = errors.New("grant type must be 'password' or 'refresh_token'")
	ErrInvalidJWT            = errors.New("invalid JWT")
	ErrNoRows                = errors.New("error no rows")
	ErrRefreshTokenNotFound  = errors.New("refresh token not found")
	ErrRefreshTokenInvalid   = errors.New("refresh token invalid")
)
