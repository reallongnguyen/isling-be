package account

import (
	controller_v1 "isling-be/internal/account/controller/http/v1"
	"isling-be/internal/account/repo"
	"isling-be/internal/account/usecase"
	"isling-be/pkg/logger"
	"isling-be/pkg/postgres"

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
func Register(pg *postgres.Postgres, l logger.Interface, handler *echo.Echo) {
	accountRepo := repo.NewAccountRepo(pg)
	refreshTokenRepo := repo.NewRefreshTokenRepo(pg)

	accountUC := usecase.NewAccountUC(accountRepo, l)
	authUC := usecase.NewAuthUsecase(l, accountRepo, refreshTokenRepo)

	groupV1 := handler.Group("/v1")

	groupV1.GET("/swagger/accounts/*", echo_swagger.WrapHandler)

	{
		controller_v1.NewAccountsRouter(groupV1, accountUC, l)
		controller_v1.NewAuthRouter(groupV1, l, authUC, accountUC)
	}
}
