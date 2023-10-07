package v1

import (
	"context"
	appresponse "isling-be/internal/common/controller/http"
	"isling-be/internal/common/controller/http/middleware"
	"isling-be/internal/play-isling/controller/http/v1/dto"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/facade"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zhenghaoz/gorse/client"
	"golang.org/x/exp/slices"
)

var feedbackTypes = []string{
	"read",
	"comment",
	"like",
	"share",
	"save",
	"reaction",
	"add-item",
	"watch-15min",
	"watch-1h",
}

type TrackingRouter struct {
	recommendationUC usecase.RecommendationUsecase
	playUserUC       usecase.PlayUserUsecase
}

func NewTrackingRouter(recommendationUC usecase.RecommendationUsecase, playUserUC usecase.PlayUserUsecase) *TrackingRouter {
	return &TrackingRouter{
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

	action := *createActionDTO.ToCreateActionRequest(accountID)

	if action.ObjectID == nil || !slices.Contains(feedbackTypes, action.Type) {
		return appresponse.ResponseSuccess(c, true)
	}

	feedback := client.Feedback{
		FeedbackType: action.Type,
		UserId:       strconv.FormatInt(int64(action.AccountID), 10),
		ItemId:       *action.ObjectID,
		Timestamp:    action.Timestamp.Format(time.RFC3339),
	}

	err = r.recommendationUC.InsertFeedback(c.Request().Context(), feedback)
	if err != nil {
		return appresponse.ResponseError(c, err)
	}

	if strings.HasPrefix(createActionDTO.Type, "watch-") && createActionDTO.ObjectID != nil {
		go func() {
			roomID, err := strconv.Atoi(*createActionDTO.ObjectID)
			if err != nil {
				facade.Log().Error("tracking router: parse ObjectID: %w", err)

				return
			}

			ctx := context.Background()

			err = r.playUserUC.InsertRecentlyJoinedRoom(ctx, accountID, int64(roomID))
			if err != nil {
				facade.Log().Error("tracking router: insert recently joined room: %w", err)
			}
		}()
	}

	return appresponse.ResponseSuccess(c, true)
}
