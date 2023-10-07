package eventtracking

import (
	"encoding/json"
	"isling-be/config"
	appresponse "isling-be/internal/common/controller/http"
	"isling-be/internal/common/controller/http/middleware"
	cm_entity "isling-be/internal/common/entity"
	"isling-be/internal/event-tracking/entity"
	"isling-be/internal/event-tracking/repo"
	"isling-be/pkg/logger"
	"isling-be/pkg/surreal"
	"net/http"
	"strconv"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/labstack/echo/v4"
	"github.com/mileusna/useragent"
	"golang.org/x/exp/slices"
)

var eventOnItems = []string{
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
	l logger.Interface,
	cfg *config.Config,
	cache *ristretto.Cache,
	handler *echo.Echo,
	sur *surreal.Surreal,
	msgBus *map[string]chan string,
) func() {
	userActRepo := repo.UserActBatch{
		MaxBatchSize:        10000,
		BatchTimeout:        2 * time.Second,
		PendingWorkCapacity: 80000,
		UserActRepo:         repo.NewUserActSurRepo(sur),
	}

	if err := userActRepo.Start(); err != nil {
		l.Error("event-tracking: start userActBatch: %w", err)

		return func() {}
	}

	userActChan := (*msgBus)["userActivityOnItem"]

	// TODO: separate the code below to route, usecase file
	handler.POST("/v1/tracking/user-activities", func(c echo.Context) error {
		accountID, _ := middleware.GetAccountIDFromJWT(c)
		uaString := c.Request().Header.Get("User-Agent")
		var userAgent cm_entity.UserAgent

		data, found := cache.Get(uaString)

		if found {
			userAgent = data.(cm_entity.UserAgent)
		} else {
			ua := useragent.Parse(uaString)
			userAgent.From(ua)
			cache.Set(uaString, userAgent, 100)
			cache.Wait()
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

		if slices.Contains(eventOnItems, userActivity.EventName) {
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

				userActChan <- string(data)
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
			l.Error("stop userActMuster: %w", err)
		}
	}
}
