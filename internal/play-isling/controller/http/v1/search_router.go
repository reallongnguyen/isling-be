package v1

import (
	common_http "isling-be/internal/common/controller/http"
	common_mw "isling-be/internal/common/controller/http/middleware"
	"isling-be/internal/play-isling/controller/http/v1/dto"
	"isling-be/internal/play-isling/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SearchRouter struct {
	searchUC usecase.SearchUsecase
}

func NewSearchRouter(searchUC usecase.SearchUsecase) *SearchRouter {
	return &SearchRouter{
		searchUC: searchUC,
	}
}

func (r *SearchRouter) Search(c echo.Context) error {
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	reqDTO := dto.SearchReqDTO{
		Query:  c.QueryParam("query"),
		Limit:  c.QueryParam("limit"),
		Offset: c.QueryParam("offset"),
	}

	if err = c.Validate(reqDTO); err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, "validation failed", []error{err})
	}

	res, err := r.searchUC.Search(c.Request().Context(), accountID, reqDTO.ToReq())
	if err != nil {
		return common_http.ResponseError(c, err)
	}

	return common_http.ResponseSuccess(c, res)
}
