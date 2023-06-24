package v1

import (
	common_mw "isling-be/internal/common/controller/http/middleware"
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/controller/http/v1/dto"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/logger"
	"net/http"

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

	g := e.Group("/rooms", common_mw.VerifyJWT())
	g.POST("", router.createRoom)

	return router
}

func (router *RoomRouter) createRoom(c echo.Context) error {
	accountID, err := common_mw.GetAccountIDFromJWT(c)
	if err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, err.Error(), []error{err})
	}

	createRoomDTO := dto.CreateRoomDTO{}

	if err = c.Bind(&createRoomDTO); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "bad request", []error{err})
	}

	if err = c.Validate(createRoomDTO); err != nil {
		return common_entity.ResponseError(c, http.StatusBadRequest, "validation error", []error{err})
	}

	room, err := router.roomUC.CreateRoom(c.Request().Context(), accountID, createRoomDTO.ToCreateRoomReq())
	if err != nil {
		return common_entity.ResponseError(c, http.StatusInternalServerError, "server error", []error{err})
	}

	return common_entity.ResponseSuccess(c, http.StatusCreated, "success", room)
}
