package app

import (
	account "isling-be/internal/account"
	eventtracking "isling-be/internal/event-tracking"
	playisling "isling-be/internal/play-isling"
	tool "isling-be/internal/tool"
	"isling-be/pkg/postgres"
	"isling-be/pkg/surreal"

	"github.com/labstack/echo/v4"
)

func useModules(
	handler *echo.Echo,
	pg *postgres.Postgres,
	sur *surreal.Surreal,
) func() {
	account.Register(handler, pg)
	stopPlay := playisling.Register(handler, pg, sur)
	tool.Register(pg, sur)
	stopEventTracking := eventtracking.Register(handler, sur)

	return func() {
		stopPlay()
		stopEventTracking()
	}
}
