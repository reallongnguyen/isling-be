package playisling

import (
	v1 "isling-be/internal/play-isling/controller/http/v1"
	"isling-be/internal/play-isling/repo"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/logger"
	"isling-be/pkg/postgres"

	"github.com/labstack/echo/v4"
)

func Register(l logger.Interface, handler *echo.Echo, pg *postgres.Postgres) {
	groupV1 := handler.Group("/play-isling/v1")

	roomRepo := repo.NewRoomRepo(pg)

	roomUC := usecase.NewRoomUsecase(l, roomRepo)

	{
		v1.NewRoomRouter(groupV1, l, roomUC)
	}
}
