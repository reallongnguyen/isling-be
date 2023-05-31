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

type AuthRouter struct {
	log    logger.Interface
	authUC usecase.AuthUsecase
}

func NewAuthRouter(e *echo.Group, log logger.Interface, authUC usecase.AuthUsecase) *AuthRouter {
	router := &AuthRouter{
		log:    log,
		authUC: authUC,
	}

	group := e.Group("/auth")
	group.POST("/tokens", router.GetToken)
	group.POST("/logout", router.Logout, common_mw.VerifyJWT())

	return router
}

func (ar *AuthRouter) GetToken(c echo.Context) error {
	getTokenByPasswordDTO := dto.GetTokenByPasswordRequestDTO{
		Email:    c.QueryParam("email"),
		Password: c.QueryParam("password"),
	}

	if err := c.Validate(getTokenByPasswordDTO); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "validation error", []error{err})
	}

	token, err := ar.authUC.GetTokenByPassword(c.Request().Context(), getTokenByPasswordDTO.ToRequest())

	if err != nil {
		switch {
		case errors.Is(err, common_entity.ErrEmailPasswordNotMatch):
			return common_entity.ResponseError(c, http.StatusUnauthorized, "sign in error", []error{err})
		case errors.Is(err, common_entity.ErrAccountNotFound):
			return common_entity.ResponseError(c, http.StatusUnauthorized, "sign in error", []error{err})
		default:
			return common_entity.ResponseError(c, http.StatusInternalServerError, "server error", []error{err})
		}
	}

	return common_entity.ResponseSuccess(c, http.StatusOK, "sign in successfully", dto.FromGetTokenRequestToDTO(token))
}

func (ar *AuthRouter) Logout(c echo.Context) error {
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "invalid JWT", []error{err})
	}

	refreshToken := c.QueryParam("refresh_token")

	err = ar.authUC.Logout(c.Request().Context(), accountID, refreshToken)
	if err != nil {
		if errors.Is(err, common_entity.ErrRefreshTokenNotFound) {
			return common_entity.ResponseError(c, http.StatusNotFound, "refresh token not found", []error{err})
		}

		return common_entity.ResponseError(c, http.StatusInternalServerError, "server error", []error{err})
	}

	return common_entity.ResponseSuccess(c, http.StatusOK, "logout success fully", "")
}
