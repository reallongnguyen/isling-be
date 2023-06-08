package email_sender

import (
	controller_v1 "isling-be/internal/email-sender/controller/http/v1"
	"isling-be/internal/email-sender/repo"
	"isling-be/internal/email-sender/usecase"
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
	emailClientFaker := repo.NewEmailClientFaker()

	emailSenderUC := usecase.NewEmailSenderUC(l, emailClientFaker)

	groupV1 := handler.Group("/v1")

	groupV1.GET("/swagger/email-senders/*", echo_swagger.WrapHandler)

	{
		controller_v1.NewEmailSenderRouter(groupV1, l, emailSenderUC)
	}
}
