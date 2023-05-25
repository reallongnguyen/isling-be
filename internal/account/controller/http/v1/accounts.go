package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"isling-be/internal/account/controller/http/v1/dto"
	"isling-be/internal/account/usecase"
	common_entity "isling-be/internal/common/entity"
	"isling-be/pkg/logger"

	"github.com/labstack/echo/v4"
)

type AccountsRouter struct {
	accountUC usecase.AccountUsecase
	log       logger.Interface
}

func NewAccountsRouter(e *echo.Group, accountUC usecase.AccountUsecase, log logger.Interface) *AccountsRouter {
	router := AccountsRouter{accountUC: accountUC, log: log}
	group := e.Group("/accounts")
	group.GET("/:accountId", router.getOne)
	group.POST("", router.create)

	return &router
}

func (router *AccountsRouter) getOne(c echo.Context) error {
	accountIDParam := c.Param("accountId")

	accountID, err := strconv.Atoi(accountIDParam)
	if err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, err.Error())
	}

	account, err := router.accountUC.GetAccountByID(c.Request().Context(), common_entity.AccountID(accountID))

	if errors.Is(err, common_entity.ErrNoRows) {
		router.log.Info("account has id %d not found", accountID)

		return common_entity.ResponseError(c, http.StatusNotFound, fmt.Sprintf("user has id %d not found", accountID))
	}

	if err != nil {
		return common_entity.ResponseError(c, http.StatusInternalServerError, "server error")
	}

	return common_entity.ResponseSuccess(c, http.StatusOK, "success", account)
}

func (router *AccountsRouter) create(c echo.Context) error {
	createAccountDto := dto.CreateAccountDto{}

	if err := c.Bind(&createAccountDto); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(createAccountDto); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, err.Error())
	}

	account, err := router.accountUC.CreateAccount(c.Request().Context(), createAccountDto.ToCreateAccountRequest())
	if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, common_entity.ErrDuplicated) {
			code = http.StatusConflict
		}

		return common_entity.ResponseError(c, code, err.Error())
	}

	return common_entity.ResponseSuccess(c, http.StatusCreated, "create one user successfully", account)
}
