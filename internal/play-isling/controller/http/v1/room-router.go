package v1

import (
	common_http "isling-be/internal/common/controller/http"
	common_mw "isling-be/internal/common/controller/http/middleware"
	"isling-be/internal/play-isling/controller/http/v1/dto"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/logger"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RoomRouter struct {
	log    logger.Interface
	roomUC usecase.RoomUsecase
}

func NewRoomRouter(e *echo.Group, log logger.Interface, roomUC usecase.RoomUsecase) *RoomRouter {
	router := &RoomRouter{
		log:    log,
		roomUC: roomUC,
	}

	g := e.Group("/rooms")
	g.POST("", router.createRoom, common_mw.VerifyJWT())
	g.GET("", router.listRoom, common_mw.VerifyJWT())
	g.GET("/:slugName", router.getRoom, common_mw.ParseJWT())
	g.PATCH("/:id", router.updateRoom, common_mw.VerifyJWT())
	g.DELETE("/:id", router.deleteRoom, common_mw.VerifyJWT())

	return router
}

func (router *RoomRouter) createRoom(c echo.Context) error {
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

func (router *RoomRouter) listRoom(c echo.Context) error {
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

func (router *RoomRouter) getRoom(c echo.Context) error {
	accountID, _ := common_mw.GetAccountIDFromJWT(c)

	slugName := c.Param("slugName")

	room, err := router.roomUC.GetRoom(c.Request().Context(), accountID, slugName)
	if err != nil {
		return common_http.ResponseError(c, err)
	}

	return common_http.ResponseSuccess(c, room)
}

func (router *RoomRouter) updateRoom(c echo.Context) error {
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

func (router *RoomRouter) deleteRoom(c echo.Context) error {
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
