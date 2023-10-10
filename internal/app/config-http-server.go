package app

import (
	"isling-be/pkg/facade"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echo_swagger "github.com/swaggo/echo-swagger"
	"golang.org/x/time/rate"
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
	limit := rate.Limit(facade.Config().HTTP.RateLimit)

	handler.Use(
		middleware.RateLimiter(
			middleware.NewRateLimiterMemoryStore(limit),
		),
	)

	handler.Use(middleware.Logger())
	handler.Use(middleware.Recover())
	handler.Use(middleware.CORS())
	handler.Use(echoprometheus.NewMiddleware(facade.Config().App.Name))
	handler.Use(middleware.Secure())
	handler.Validator = NewCustomValidator()

	// health check
	handler.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	handler.GET("/healthz", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// stress test empty case
	handler.GET("/stress-test", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	handler.POST("/stress-test", func(c echo.Context) error {
		body := new(struct {
			title       string
			description string
		})

		if err := c.Bind(body); err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		return c.NoContent(http.StatusOK)
	})

	handler.GET("/metrics", echoprometheus.NewHandler())

	handler.GET("/swagger/*", echo_swagger.WrapHandler)
}
