package entity

import (
	"github.com/labstack/echo/v4"
)

type AppResponse[T any] struct {
	Success bool   `json:"success" example:"true"`
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"SUCCESS"`
	Data    T      `json:"data,omitempty" example:"{ 'userID': 1, 'name': 'Luffy' }"`
}

func ResponseError(c echo.Context, code int, msg string) error {
	err := c.JSON(code, AppResponse[any]{
		Success: false,
		Code:    code,
		Message: msg,
	})
	if err != nil {
		return err
	}

	return echo.NewHTTPError(code, msg)
}

func ResponseSuccess[T any](c echo.Context, code int, msg string, data T) error {
	return c.JSON(code, AppResponse[T]{
		Success: true,
		Code:    code,
		Message: msg,
		Data:    data,
	})
}
