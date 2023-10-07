package middleware

import (
	"errors"
	"fmt"
	"isling-be/config"
	"isling-be/internal/common/entity"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echo_jwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var cfg, _ = config.NewConfig()

func VerifyJWT() echo.MiddlewareFunc {
	return echo_jwt.WithConfig(echo_jwt.Config{
		SigningKey: []byte(cfg.JWT.JWTSecretKey),
		ErrorHandler: func(c echo.Context, err error) error {
			customErr := fmt.Errorf("access token: %w", err)

			entity.ResponseError(c, http.StatusUnauthorized, "unauthorized", []error{customErr})

			return err
		},
		ContextKey: "account",
	})
}

func ParseJWT() echo.MiddlewareFunc {
	return echo_jwt.WithConfig(echo_jwt.Config{
		SigningKey: []byte(cfg.JWT.JWTSecretKey),
		ErrorHandler: func(c echo.Context, err error) error {
			if errors.Is(err, echo_jwt.ErrJWTInvalid) {
				entity.ResponseError(c, http.StatusUnauthorized, "unauthorized", []error{err})

				return err
			}

			return nil
		},
		ContinueOnIgnoredError: true,
		ContextKey:             "account",
	})
}

func GetAccountIDFromJWT(c echo.Context) (entity.AccountID, error) {
	account, ok := c.Get("account").(*jwt.Token)
	if !ok {
		return 0, entity.ErrInvalidJWT
	}

	claim, ok := account.Claims.(jwt.MapClaims)
	if !ok {
		return 0, entity.ErrInvalidJWT
	}

	accountID, ok := claim["account_id"].(float64)
	if !ok {
		return 0, entity.ErrInvalidJWT
	}

	return entity.AccountID(accountID), nil
}
