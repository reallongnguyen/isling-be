package app

import (
	account "isling-be/internal/account"
	emailSender "isling-be/internal/email-sender"
	"isling-be/pkg/logger"
	"isling-be/pkg/postgres"

	"github.com/labstack/echo/v4"
)

func useModules(pg *postgres.Postgres, l logger.Interface, handler *echo.Echo) {
	account.Register(pg, l, handler)
	emailSender.Register(pg, l, handler)
}
