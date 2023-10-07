package v1

import (
	common_http "isling-be/internal/common/controller/http"
	common_mw "isling-be/internal/common/controller/http/middleware"
	"isling-be/internal/play-isling/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HomeRouter struct {
	homeUC usecase.HomeUsecase
}

func NewHomeRouter(homeUC usecase.HomeUsecase) *HomeRouter {
	return &HomeRouter{
		homeUC: homeUC,
	}
}

func (router *HomeRouter) Show(c echo.Context) error {
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	homeRes, err := router.homeUC.Show(c.Request().Context(), accountID)
	if err != nil {
		return common_http.ResponseError(c, err)
	}

	return common_http.ResponseSuccess(c, homeRes)
}

func (router *HomeRouter) ShowGuest(c echo.Context) error {
	homeRes, err := router.homeUC.ShowGuest(c.Request().Context())
	if err != nil {
		return common_http.ResponseError(c, err)
	}

	return common_http.ResponseSuccess(c, homeRes)
}
