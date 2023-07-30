package middleware

import (
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

			return entity.ResponseError(c, http.StatusUnauthorized, "unauthorized", []error{customErr})
		},
		ContextKey: "account",
	})
}

func ParseJWT() echo.MiddlewareFunc {
	return echo_jwt.WithConfig(echo_jwt.Config{
		SigningKey: []byte(cfg.JWT.JWTSecretKey),
		ErrorHandler: func(c echo.Context, err error) error {
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
