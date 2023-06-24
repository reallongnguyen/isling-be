package v1

import (
	"errors"
	"net/http"

	"isling-be/internal/account/controller/http/v1/dto"
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
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	createProfileDTO := dto.CreateProfileReqDTO{}
	if err = c.Bind(&createProfileDTO); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "bad request", []error{err})
	}

	if err = c.Validate(createProfileDTO); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "validation failed", []error{err})
	}

	profile, err := router.profileUC.UpsertProfile(c.Request().Context(), accountID, createProfileDTO.ToRequest())
	if err != nil {
		if errors.Is(err, common_entity.ErrAccountIDDuplicated) {
			return common_entity.ResponseError(c, http.StatusBadRequest, "bad request", []error{err})
		}

		return common_entity.ResponseError(c, http.StatusInternalServerError, "server error", []error{err})
	}

	return common_entity.ResponseSuccess(c, http.StatusCreated, "upsert profile successfully", profile)
}

func (router *ProfilesRouter) getProfile(c echo.Context) error {
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	profile, err := router.profileUC.GetProfile(c.Request().Context(), accountID)

	if errors.Is(err, common_entity.ErrAccountNotFound) {
		router.log.Info("get profile: account id %d not found", accountID)

		return common_entity.ResponseError(c, http.StatusUnauthorized, "invalid account", []error{err})
	}

	if err != nil {
		return common_entity.ResponseError(c, http.StatusInternalServerError, "server error", []error{err})
	}

	return common_entity.ResponseSuccess(c, http.StatusOK, "success", profile)
}
