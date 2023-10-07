package v1

import (
	"errors"
	"net/http"
	"strconv"

	"isling-be/internal/account/controller/http/v1/dto"
	"isling-be/internal/account/usecase"
	common_mw "isling-be/internal/common/controller/http/middleware"
	common_entity "isling-be/internal/common/entity"
	"isling-be/pkg/facade"

	"github.com/labstack/echo/v4"
)

type AccountsRouter struct {
	accountUC usecase.AccountUsecase
}

func NewAccountsRouter(e *echo.Group, accountUC usecase.AccountUsecase) *AccountsRouter {
	router := AccountsRouter{accountUC: accountUC}
	group := e.Group("/accounts", common_mw.VerifyJWT())
	group.POST("", router.create)
	group.GET("/:accountID", router.getOne)
	group.PATCH("/me/password", router.changePassword)

	return &router
}

func (router *AccountsRouter) create(c echo.Context) error {
	createAccountDto := dto.CreateAccountDto{}

	if err := c.Bind(&createAccountDto); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	if err := c.Validate(createAccountDto); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	account, err := router.accountUC.CreateAccount(c.Request().Context(), createAccountDto.ToCreateAccountRequest())
	if err != nil {
		code := http.StatusInternalServerError

		if errors.Is(err, common_entity.ErrEmailDuplicated) {
			code = http.StatusConflict
		}

		return common_entity.ResponseError(c, code, err.Error(), []error{err})
	}

	return common_entity.ResponseSuccess(c, http.StatusCreated, "create one user successfully", account)
}

func (router *AccountsRouter) getOne(c echo.Context) error {
	accountIDParam := c.Param("accountID")

	accountID, err := strconv.Atoi(accountIDParam)
	if err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	account, err := router.accountUC.GetAccountByID(c.Request().Context(), common_entity.AccountID(accountID))

	if errors.Is(err, common_entity.ErrAccountNotFound) {
		facade.Log().Info("get one: account id %d not found", accountID)

		return common_entity.ResponseError(c, http.StatusNotFound, "not found", []error{err})
	}

	if err != nil {
		return common_entity.ResponseError(c, http.StatusInternalServerError, "server error", []error{err})
	}

	return common_entity.ResponseSuccess(c, http.StatusOK, "success", account)
}

func (router *AccountsRouter) changePassword(c echo.Context) error {
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	changePasswordDTO := dto.ChangePasswordDto{}

	if err := c.Bind(&changePasswordDTO); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "parse request body failed", []error{err})
	}

	if err := c.Validate(changePasswordDTO); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "validation failed", []error{err})
	}

	if err := router.accountUC.ChangePassword(c.Request().Context(), accountID, changePasswordDTO.ToChangePasswordRequest()); err != nil {
		switch {
		case errors.Is(err, common_entity.ErrAccountNotFound):
			return common_entity.ResponseError(c, http.StatusBadRequest, "account not found", []error{err})
		case errors.Is(err, common_entity.ErrPasswordNotCorrect):
			return common_entity.ResponseError(c, http.StatusBadRequest, "old password not correct", []error{err})
		default:
			return common_entity.ResponseError(c, http.StatusInternalServerError, "server error", []error{err})
		}
	}

	return common_entity.ResponseSuccess(c, http.StatusOK, "change password successfully", "")
}
