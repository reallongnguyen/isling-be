package entity

import (
	"github.com/labstack/echo/v4"
)

const (
	errMsgsInitialCap = 1
)

type AppResponse[T any] struct {
	Success bool     `json:"success" example:"true"`
	Code    int      `json:"code" example:"200"`
	Message string   `json:"message" example:"SUCCESS"`
	Data    T        `json:"data,omitempty" example:"{ 'accountID': 1, 'name': 'Luffy' }"`
	Errors  []string `json:"errors,omitempty"`
}

func ResponseError(c echo.Context, code int, msg string, errs []error) error {
	errMsgs := make([]string, 0, errMsgsInitialCap)

	for _, err := range errs {
		errMsgs = append(errMsgs, err.Error())
	}

	return c.JSON(code, AppResponse[any]{
		Success: false,
		Code:    code,
		Message: msg,
		Errors:  errMsgs,
	})
}

func ResponseSuccess[T any](c echo.Context, code int, msg string, data T) error {
	return c.JSON(code, AppResponse[T]{
		Success: true,
		Code:    code,
		Message: msg,
		Data:    data,
	})
}
