package eventtracking

import (
	"encoding/json"
	appresponse "isling-be/internal/common/controller/http"
	"isling-be/internal/common/controller/http/middleware"
	cm_entity "isling-be/internal/common/entity"
	"isling-be/internal/event-tracking/entity"
	"isling-be/internal/event-tracking/repo"
	"isling-be/internal/event-tracking/usecase"
	"isling-be/pkg/facade"
	"isling-be/pkg/surreal"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mileusna/useragent"
	"golang.org/x/exp/slices"
)

var recommendEventList = []string{
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

func Register(
	handler *echo.Echo,
	sur *surreal.Surreal,
) func() {
	userActRepo := usecase.UserActBatch{
		MaxBatchSize:        10000,
		BatchTimeout:        2 * time.Second,
		PendingWorkCapacity: 80000,
		UserActRepo:         repo.NewUserActSurRepo(sur),
	}

	if err := userActRepo.Start(); err != nil {
		facade.Log().Error("event-tracking: start userActBatch: %w", err)

		return func() {}
	}

	// TODO: separate the code below to route, usecase file
	handler.POST("/v1/tracking/user-activities", func(c echo.Context) error {
		accountID, _ := middleware.GetAccountIDFromJWT(c)
		uaString := c.Request().Header.Get("User-Agent")
		var userAgent cm_entity.UserAgent
		hasUserAgentInCache := false

		data, found := facade.Cache().Get(uaString)

		if found {
			userAgent, hasUserAgentInCache = data.(cm_entity.UserAgent)
		}

		if !hasUserAgentInCache {
			ua := useragent.Parse(uaString)
			userAgent.From(ua)
			facade.Cache().Set(uaString, userAgent, 100)
			facade.Cache().Wait()
		}

		dto := new(entity.CreateUserActivityDTO[any])

		if err := c.Bind(dto); err != nil {
			return appresponse.ResponseCustomError(c, http.StatusBadRequest, "", []error{err})
		}

		userActivity := entity.UserActivity[any]{
			UserID:    strconv.Itoa(int(accountID)), // 0 meaning empty
			EventName: dto.EventName,
			Data:      dto.Data,
			Device:    userAgent.Device,
			OS:        userAgent.OS,
			App:       dto.App,
			Timestamp: time.Now(),
		}

		if userActivity.UserID != "0" && slices.Contains(recommendEventList, userActivity.EventName) {
			byteOfData, _ := json.Marshal(userActivity.Data)
			actOnItemData := new(entity.ActOnItemData)
			if err := json.Unmarshal(byteOfData, actOnItemData); err != nil {
				return appresponse.ResponseCustomError(c, http.StatusBadRequest, "", []error{err})
			}

			userActivity.Data = actOnItemData

			go func() {
				data, err := json.Marshal(userActivity)
				if err != nil {
					return
				}

				if err := facade.Pubsub().Publish("recommend.feedback", data, nil); err != nil {
					facade.Log().Info("publish 'recommend.feedback' %s error %w", data, err)
				}

				if userActivity.EventName == "read" {
					if err := facade.Pubsub().Publish("room.watched", data, nil); err != nil {
						facade.Log().Info("publish 'room.watched' %s error %w", data, err)
					}
				}
			}()
		}

		err := userActRepo.Add(userActivity)
		if err != nil {
			return appresponse.ResponseError(c, err)
		}

		return appresponse.ResponseSuccess(c, true)
	}, middleware.ParseJWT())

	return func() {
		if err := userActRepo.Stop(); err != nil {
			facade.Log().Error("stop userActMuster: %w", err)
		}
	}
}
