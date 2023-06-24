package app

import (
	account "isling-be/internal/account"
	emailSender "isling-be/internal/email-sender"
	playisling "isling-be/internal/play-isling"
	"isling-be/pkg/logger"
	"isling-be/pkg/postgres"

	"github.com/labstack/echo/v4"
)

func useModules(pg *postgres.Postgres, l logger.Interface, handler *echo.Echo) {
	account.Register(pg, l, handler)
	playisling.Register(l, handler, pg)
	emailSender.Register(pg, l, handler)
}
