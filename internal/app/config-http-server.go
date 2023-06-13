package app

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echo_swagger "github.com/swaggo/echo-swagger"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewCustomValidator() *CustomValidator {
	customValidator := CustomValidator{validator: validator.New()}
	customValidator.validator.RegisterValidation("alphaunicodespace", ValidateAlphaUnicodeWithSpace)
	customValidator.validator.RegisterValidation("beforenow", ValidateBeforeNow)

	return &customValidator
}

// Swagger spec:
// @title Isling Open API
// @version 1.0
// @description This is a Isling Open API.

// @contact.name Isling Open API Support
// @contact.email api@isling.me

// @host https://api.isling.me
// @BasePath /v1.
func configHTTPServer(handler *echo.Echo) {
	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())
	handler.Use(middleware.CORS())
	handler.Validator = NewCustomValidator()

	handler.GET("/healthz", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	handler.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	handler.GET("/swagger/*", echo_swagger.WrapHandler)
}
