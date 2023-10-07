package v1

import (
	common_http "isling-be/internal/common/controller/http"
	common_mw "isling-be/internal/common/controller/http/middleware"
	"isling-be/internal/play-isling/controller/http/v1/dto"
	"isling-be/internal/play-isling/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RoomRouter struct {
	roomUC usecase.RoomUsecase
}

func NewRoomRouter(roomUC usecase.RoomUsecase) *RoomRouter {
	return &RoomRouter{
		roomUC: roomUC,
	}
}

func (router *RoomRouter) Create(c echo.Context) error {
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	createRoomDTO := dto.CreateRoomDTO{}

	if err = c.Bind(&createRoomDTO); err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, "bad request", []error{err})
	}

	if err = c.Validate(createRoomDTO); err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, "validation error", []error{err})
	}

	room, err := router.roomUC.CreateRoom(c.Request().Context(), accountID, createRoomDTO.ToCreateRoomReq())
	if err != nil {
		return common_http.ResponseCustomError(c, http.StatusInternalServerError, "server error", []error{err})
	}

	return common_http.ResponseSuccess(c, room)
}

func (router *RoomRouter) List(c echo.Context) error {
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	roomCollection, err := router.roomUC.GetManyRoomOfUser(c.Request().Context(), accountID)
	if err != nil {
		return common_http.ResponseCustomError(c, http.StatusInternalServerError, "server error", []error{err})
	}

	return common_http.ResponseSuccess(c, roomCollection)
}

func (router *RoomRouter) Show(c echo.Context) error {
	accountID, _ := common_mw.GetAccountIDFromJWT(c)

	slugName := c.Param("slugName")

	room, err := router.roomUC.GetRoom(c.Request().Context(), accountID, slugName)
	if err != nil {
		return common_http.ResponseError(c, err)
	}

	return common_http.ResponseSuccess(c, room)
}

func (router *RoomRouter) Update(c echo.Context) error {
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	updateRoomDTO := dto.UpdateRoomDTO{}

	err = c.Bind(&updateRoomDTO)
	if err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, "bad request", []error{err})
	}

	if err = c.Validate(updateRoomDTO); err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, "validation error", []error{err})
	}

	room, err := router.roomUC.UpdateRoom(c.Request().Context(), accountID, int64(id), updateRoomDTO.ToUpdateRoomReq())
	if err != nil {
		return common_http.ResponseError(c, err)
	}

	return common_http.ResponseSuccess(c, room)
}

func (router *RoomRouter) Delete(c echo.Context) error {
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return common_http.ResponseCustomError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	err = router.roomUC.DeleteRoom(c.Request().Context(), accountID, int64(id))

	if err != nil {
		return common_http.ResponseError(c, err)
	}

	return common_http.ResponseSuccess(c, "")
}
