package app

import (
	account "isling-be/internal/account"
	emailSender "isling-be/internal/email-sender"
	playisling "isling-be/internal/play-isling"
	tool "isling-be/internal/tool"
	"isling-be/pkg/logger"
	"isling-be/pkg/postgres"
	"isling-be/pkg/surreal"

	"github.com/labstack/echo/v4"
)

func useModules(
	pg *postgres.Postgres,
	sur *surreal.Surreal,
	l logger.Interface,
	handler *echo.Echo,
	msgBus *map[string]chan string,
) {
	account.Register(pg, l, handler, msgBus)
	playisling.Register(l, handler, pg, msgBus)
	emailSender.Register(pg, l, handler)
	tool.Register(l, pg, sur)
}
