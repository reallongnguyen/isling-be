package app

import (
	"github.com/labstack/echo/v4"

	account "github.com/btcs-longnp/isling-be/internal/account"
	"github.com/btcs-longnp/isling-be/pkg/logger"
	"github.com/btcs-longnp/isling-be/pkg/postgres"
)

func useModules(pg *postgres.Postgres, l logger.Interface, handler *echo.Echo) {
	account.Register(pg, l, handler)
}
