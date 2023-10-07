package entity

import (
	"errors"
)

var (
	ErrAccountIDDuplicated   = errors.New("account id duplicated")
	ErrAccountNotFound       = errors.New("account not found")
	ErrEmailDuplicated       = errors.New("email address duplicated")
	ErrEmailPasswordNotMatch = errors.New("email password not match")
	ErrGrantTypeInvalid      = errors.New("grant type must be 'password' or 'refresh_token'")
	ErrInvalidJWT            = errors.New("invalid JWT")
	ErrNoRows                = errors.New("error no rows")
	ErrRefreshTokenNotFound  = errors.New("refresh token not found")
	ErrRefreshTokenInvalid   = errors.New("refresh token invalid")
	ErrPasswordNotCorrect    = errors.New("password not correct")
	ErrCanNotParseReqBody    = errors.New("can not parse request's body")
)
