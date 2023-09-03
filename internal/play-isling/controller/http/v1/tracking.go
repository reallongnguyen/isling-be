package v1

import (
	"context"
	appresponse "isling-be/internal/common/controller/http"
	"isling-be/internal/common/controller/http/middleware"
	"isling-be/internal/play-isling/controller/http/v1/dto"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/logger"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type TrackingRouter struct {
	log              logger.Interface
	recommendationUC usecase.RecommendationUsecase
	playUserUC       usecase.PlayUserUsecase
}

func NewTrackingRouter(log logger.Interface, recommendationUC usecase.RecommendationUsecase, playUserUC usecase.PlayUserUsecase) *TrackingRouter {
	return &TrackingRouter{
		log:              log,
		recommendationUC: recommendationUC,
		playUserUC:       playUserUC,
	}
}

func (r *TrackingRouter) Create(c echo.Context) error {
	accountID, err := middleware.GetAccountIDFromJWT(c)
	if err != nil {
		return appresponse.ResponseCustomError(c, 400, "can not detect account id", []error{err})
	}

	createActionDTO := new(dto.CreateAction)
	if err = c.Bind(createActionDTO); err != nil {
		return appresponse.ResponseCustomError(c, 400, "can not parse request body", []error{err})
	}

	if err = c.Validate(createActionDTO); err != nil {
		return appresponse.ResponseCustomError(c, 400, "validation error", []error{err})
	}

	err = r.recommendationUC.InsertFeedback(c.Request().Context(), []usecase.CreateActionRequest{*createActionDTO.ToCreateActionRequest(accountID)})
	if err != nil {
		return appresponse.ResponseError(c, err)
	}

	if strings.HasPrefix(createActionDTO.Type, "watch-") && createActionDTO.ObjectID != nil {
		go func() {
			roomID, err := strconv.Atoi(*createActionDTO.ObjectID)
			if err != nil {
				r.log.Error("tracking router: parse ObjectID: %w", err)

				return
			}

			ctx := context.Background()

			err = r.playUserUC.InsertRecentlyJoinedRoom(ctx, accountID, int64(roomID))
			if err != nil {
				r.log.Error("tracking router: insert recently joined room: %w", err)
			}
		}()
	}

	return appresponse.ResponseSuccess(c, true)
}
