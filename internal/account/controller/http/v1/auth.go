package v1

import (
	"errors"
	"net/http"

	"isling-be/internal/account/controller/http/v1/dto"
	"isling-be/internal/account/usecase"
	common_mw "isling-be/internal/common/controller/http/middleware"
	common_entity "isling-be/internal/common/entity"

	"github.com/labstack/echo/v4"
)

type AuthRouter struct {
	authUC    usecase.AuthUsecase
	accountUC usecase.AccountUsecase
}

func NewAuthRouter(e *echo.Group, authUC usecase.AuthUsecase, accountUC usecase.AccountUsecase) *AuthRouter {
	router := &AuthRouter{
		authUC:    authUC,
		accountUC: accountUC,
	}

	group := e.Group("/auth")
	group.POST("/signup", router.signUp)
	group.POST("/tokens", router.getToken)
	group.POST("/logout", router.logout, common_mw.VerifyJWT())

	return router
}

func (router *AuthRouter) signUp(c echo.Context) error {
	createAccountDto := dto.CreateAccountDto{}

	if err := c.Bind(&createAccountDto); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "bad request", []error{err})
	}

	if err := c.Validate(createAccountDto); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "validation failed", []error{err})
	}

	token, err := router.authUC.SignUp(c.Request().Context(), createAccountDto.ToCreateAccountRequest())

	if errors.Is(err, common_entity.ErrEmailDuplicated) {
		return common_entity.ResponseError(c, http.StatusBadRequest, "bad request", []error{err})
	}

	if err != nil {
		return common_entity.ResponseError(c, http.StatusInternalServerError, "server error", []error{err})
	}

	return common_entity.ResponseSuccess(c, http.StatusCreated, "sign up successfully", dto.FromGetTokenRequestToDTO(token))
}

func (router *AuthRouter) getToken(c echo.Context) error {
	grantType := c.QueryParam("grant_type")
	if grantType != dto.GrantTypePassword && grantType != dto.GrantTypeRefreshToken {
		err := common_entity.ErrGrantTypeInvalid

		return common_entity.ResponseError(c, http.StatusBadRequest, "validation error", []error{err})
	}

	if grantType == dto.GrantTypePassword {
		return router.getTokenByPassword(c)
	}

	return router.getTokenByRefreshToken(c)
}

func (router *AuthRouter) getTokenByPassword(c echo.Context) error {
	getTokenByPasswordDTO := dto.GetTokenByPasswordRequestDTO{
		Email:    c.QueryParam("email"),
		Password: c.QueryParam("password"),
	}

	if err := c.Validate(getTokenByPasswordDTO); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "validation error", []error{err})
	}

	token, err := router.authUC.GetTokenByPassword(c.Request().Context(), getTokenByPasswordDTO.ToRequest())

	if err != nil {
		switch {
		case errors.Is(err, common_entity.ErrEmailPasswordNotMatch) || errors.Is(err, common_entity.ErrAccountNotFound):
			return common_entity.ResponseError(c, http.StatusUnauthorized, "sign in error", []error{common_entity.ErrEmailPasswordNotMatch})
		default:
			return common_entity.ResponseError(c, http.StatusInternalServerError, "server error", []error{err})
		}
	}

	return common_entity.ResponseSuccess(c, http.StatusOK, "sign in successfully", dto.FromGetTokenRequestToDTO(token))
}

func (router *AuthRouter) getTokenByRefreshToken(c echo.Context) error {
	refreshTokenCredential := dto.GetTokenByRefreshTokenRequestDTO{
		RefreshToken: c.QueryParam("refresh_token"),
	}

	if err := c.Validate(refreshTokenCredential); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "validation error", []error{err})
	}

	token, err := router.authUC.GetTokenByRefreshToken(c.Request().Context(), refreshTokenCredential.ToRequest())

	if err != nil {
		switch {
		case errors.Is(err, common_entity.ErrRefreshTokenNotFound) ||
			errors.Is(err, common_entity.ErrRefreshTokenInvalid) ||
			errors.Is(err, common_entity.ErrAccountNotFound):
			return common_entity.ResponseError(c, http.StatusUnauthorized, "get access token error", []error{err})
		default:
			return common_entity.ResponseError(c, http.StatusInternalServerError, "server error", []error{err})
		}
	}

	return common_entity.ResponseSuccess(c, http.StatusOK, "get access token successfully", dto.FromGetTokenRequestToDTO(token))
}

func (router *AuthRouter) logout(c echo.Context) error {
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "invalid JWT", []error{err})
	}

	refreshToken := c.QueryParam("refresh_token")

	err = router.authUC.Logout(c.Request().Context(), accountID, refreshToken)
	if err != nil {
		if errors.Is(err, common_entity.ErrRefreshTokenNotFound) {
			return common_entity.ResponseError(c, http.StatusNotFound, "refresh token not found", []error{err})
		}

		return common_entity.ResponseError(c, http.StatusInternalServerError, "server error", []error{err})
	}

	return common_entity.ResponseSuccess(c, http.StatusOK, "logout success fully", "")
}
