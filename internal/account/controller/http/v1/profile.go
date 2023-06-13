package v1

import (
	"errors"
	"net/http"

	"isling-be/internal/account/usecase"
	common_mw "isling-be/internal/common/controller/http/middleware"
	common_entity "isling-be/internal/common/entity"
	"isling-be/pkg/logger"

	"github.com/labstack/echo/v4"
)

type ProfilesRouter struct {
	profileUC usecase.ProfileUsecase
	log       logger.Interface
}

func NewProfilesRouter(e *echo.Group, profileUC usecase.ProfileUsecase, log logger.Interface) *ProfilesRouter {
	router := ProfilesRouter{profileUC: profileUC, log: log}
	group := e.Group("/profiles", common_mw.VerifyJWT())
	group.POST("/me", router.create)
	group.GET("/me", router.getProfile)

	return &router
}

func (router *ProfilesRouter) create(c echo.Context) error {
	return common_entity.ResponseSuccess(c, http.StatusCreated, "create one user successfully", "")
}

func (router *ProfilesRouter) getProfile(c echo.Context) error {
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	profile, err := router.profileUC.GetProfile(c.Request().Context(), common_entity.AccountID(accountID))

	if errors.Is(err, common_entity.ErrAccountNotFound) {
		router.log.Info("get profile: account id %d not found", accountID)

		return common_entity.ResponseError(c, http.StatusNotFound, "not found", []error{err})
	}

	if err != nil {
		return common_entity.ResponseError(c, http.StatusInternalServerError, "server error", []error{err})
	}

	return common_entity.ResponseSuccess(c, http.StatusOK, "success", profile)
}
