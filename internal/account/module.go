package account

import (
	"isling-be/config"
	controller_v1 "isling-be/internal/account/controller/http/v1"
	"isling-be/internal/account/repo"
	"isling-be/internal/account/usecase"
	"isling-be/pkg/logger"
	"isling-be/pkg/postgres"

	"github.com/dgraph-io/ristretto"
	"github.com/labstack/echo/v4"
	echo_swagger "github.com/swaggo/echo-swagger"
)

// Swagger spec:
// @title Isling Open API
// @version 1.0
// @description This is a Isling Open API.

// @contact.name Isling Open API Support
// @contact.email api@isling.me

// @host https://api.isling.me
// @BasePath /v1.
func Register(l logger.Interface, _ *config.Config, _ *ristretto.Cache, handler *echo.Echo, pg *postgres.Postgres, msgBus *map[string]chan string) {
	accountRepo := repo.NewAccountRepo(pg)
	refreshTokenRepo := repo.NewRefreshTokenRepo(pg)
	profileRepo := repo.NewProfileRepo(pg)

	accountUC := usecase.NewAccountUC(accountRepo, msgBus)
	authUC := usecase.NewAuthUsecase(accountUC, accountRepo, refreshTokenRepo)
	profileUC := usecase.NewProfileUC(profileRepo)

	groupV1 := handler.Group("/v1")

	groupV1.GET("/swagger/accounts/*", echo_swagger.WrapHandler)

	{
		controller_v1.NewAccountsRouter(groupV1, accountUC)
		controller_v1.NewProfilesRouter(groupV1, profileUC)
		controller_v1.NewAuthRouter(groupV1, authUC, accountUC)
	}
}
