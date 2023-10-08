package appresponse

import (
	errormessage "isling-be/internal/common/controller/http/error-message"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AppSuccessRes[T any] struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message,omitempty" example:"SUCCESS"`
	Data    T      `json:"data,omitempty"`
}

type AppErrorRes struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message,omitempty"`
	Errors  []*errormessage.HTTPError `json:"errors"`
}

func ResponseSuccess[T any](c echo.Context, data T) error {
	return c.JSON(http.StatusOK, AppSuccessRes[T]{
		Code: http.StatusOK,
		Data: data,
	})
}

func ResponseError(c echo.Context, err error) error {
	httpError, ok := errormessage.ErrorMap[err]
	if !ok {
		httpError = &errormessage.HTTPError{
			HTTPCode: http.StatusInternalServerError,
			Code:     1000500,
			Message:  err.Error(),
		}
	}

	return c.JSON(httpError.HTTPCode, AppErrorRes{
		Code:   httpError.HTTPCode,
		Errors: []*errormessage.HTTPError{httpError},
	})
}

func ResponseCustomError(c echo.Context, code int, msg string, errs []error) error {
	errMsgs := make([]*errormessage.HTTPError, 0, len(errs))

	for _, err := range errs {
		errMsgs = append(errMsgs, &errormessage.HTTPError{
			Code:    code,
			Message: err.Error(),
		})
	}

	return c.JSON(code, AppErrorRes{
		Code:    code,
		Message: msg,
		Errors:  errMsgs,
	})
}

func ResponseCustomSuccess[T any](c echo.Context, code int, msg string, data T) error {
	return c.JSON(code, AppSuccessRes[T]{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}
