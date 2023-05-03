package app

import (
	"github.com/labstack/echo/v4"

	account "isling-be/internal/account"
	"isling-be/pkg/logger"
	"isling-be/pkg/postgres"
)

func useModules(pg *postgres.Postgres, l logger.Interface, handler *echo.Echo) {
	account.Register(pg, l, handler)
}
