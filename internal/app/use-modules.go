package app

import (
	"isling-be/config"
	account "isling-be/internal/account"
	eventtracking "isling-be/internal/event-tracking"
	playisling "isling-be/internal/play-isling"
	tool "isling-be/internal/tool"
	"isling-be/pkg/logger"
	"isling-be/pkg/postgres"
	"isling-be/pkg/surreal"

	"github.com/dgraph-io/ristretto"
	"github.com/labstack/echo/v4"
)

func useModules(
	l logger.Interface,
	cache *ristretto.Cache,
	cfg *config.Config,
	handler *echo.Echo,
	pg *postgres.Postgres,
	sur *surreal.Surreal,
	msgBus *map[string]chan string,
) func() {
	account.Register(l, cfg, cache, handler, pg, msgBus)
	stopPlay := playisling.Register(l, cfg, cache, handler, pg, msgBus)
	tool.Register(l, cfg, cache, pg, sur)
	stopEventTracking := eventtracking.Register(l, cfg, cache, handler, sur, msgBus)

	return func() {
		stopPlay()
		stopEventTracking()
	}
}
