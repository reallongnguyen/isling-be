package eventtracking

import (
	"encoding/json"
	appresponse "isling-be/internal/common/controller/http"
	mymiddleware "isling-be/internal/common/controller/http/middleware"
	cm_entity "isling-be/internal/common/entity"
	"isling-be/internal/event-tracking/delivery/http/v1/dto"
	"isling-be/internal/event-tracking/entity"
	"isling-be/internal/event-tracking/repo"
	"isling-be/internal/event-tracking/usecase"
	"isling-be/pkg/facade"
	"isling-be/pkg/surreal"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mileusna/useragent"
	"golang.org/x/exp/slices"
	"golang.org/x/time/rate"
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
	userActBatchUC := usecase.UserActBatch{
		MaxBatchSize:        1000,
		BatchTimeout:        1 * time.Minute,
		PendingWorkCapacity: 4000,
		UserActRepo:         repo.NewUserActSurRepo(sur),
	}

	if err := userActBatchUC.Start(); err != nil {
		facade.Log().Error("event-tracking: start userActBatch: %w", err)

		return func() {}
	}

	ccuLogRepo := repo.NewCCULogRepo(facade.Redis())

	ccuLogBatchUC := usecase.CCULogBatch{
		MaxBatchSize:        1000,
		BatchTimeout:        30 * time.Second,
		PendingWorkCapacity: 4000,
		CCULogRepo:          ccuLogRepo,
	}

	if err := ccuLogBatchUC.Start(); err != nil {
		facade.Log().Error("event-tracking: start ccuLogBatch: %w", err)

		return func() {}
	}

	rateLimit := rate.Limit(facade.Config().HTTP.RateLimitUserActivitiesPost)

	middlewares := []echo.MiddlewareFunc{
		middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rateLimit)),
		mymiddleware.ParseJWT(),
	}

	// TODO: separate the code below to route, usecase file
	handler.POST("/tracking/v1/user-activities", func(c echo.Context) error {
		accountID, _ := mymiddleware.GetAccountIDFromJWT(c)
		guestID := c.Request().Header.Get("X-Guest-ID")
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
			IP:        c.RealIP(),
			GuestID:   guestID,
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

		err := userActBatchUC.Add(&userActivity)
		if err != nil {
			return appresponse.ResponseError(c, err)
		}

		return appresponse.ResponseSuccess(c, true)
	}, middlewares...)

	handler.POST("/tracking/v1/ccu-logs", func(c echo.Context) error {
		accountID, _ := mymiddleware.GetAccountIDFromJWT(c)
		guestID := c.Request().Header.Get("X-Guest-ID")
		uaString := c.Request().Header.Get("User-Agent")
		var userAgent cm_entity.UserAgent
		hasUserAgentInCache := false

		data, found := facade.Cache().Get(uaString)

		if found {
			userAgent, hasUserAgentInCache = data.(cm_entity.UserAgent)
		}

		// facade.Log().Trace("%w", useragent.Parse(uaString))
		// client: ios, android, web_pc, web_mobile, other

		if !hasUserAgentInCache {
			ua := useragent.Parse(uaString)

			userAgent.From(ua)
			facade.Cache().Set(uaString, userAgent, 100)
			facade.Cache().Wait()
		}

		ccuLog := entity.CCULog{
			Device:    userAgent.Device,
			OS:        userAgent.OS,
			IP:        c.RealIP(),
			Timestamp: time.Now(),
		}

		if accountID != 0 {
			ccuLog.UserID = strconv.FormatInt(int64(accountID), 10)
		} else {
			ccuLog.UserID = guestID
		}

		err := ccuLogBatchUC.Add(&ccuLog)
		if err != nil {
			return appresponse.ResponseError(c, err)
		}

		return appresponse.ResponseSuccess(c, true)
	}, middlewares...)

	handler.GET("/tracking/v1/ccu-logs/count", func(c echo.Context) error {
		ccuQueryDTO := dto.CountCCUReq{
			Timestamp:  c.QueryParam("timestamp"),
			WindowSize: c.QueryParam("window_size"),
		}

		if ccuQueryDTO.WindowSize == "" {
			ccuQueryDTO.WindowSize = "1"
		}

		req, err := ccuQueryDTO.ToReq()
		if err != nil {
			return appresponse.ResponseCustomError(c, http.StatusBadRequest, "bad request", []error{err})
		}

		count, err := ccuLogRepo.CountCCU(req.Timestamp, req.WindowSize)
		if err != nil {
			return appresponse.ResponseError(c, err)
		}

		res := struct {
			Count int64 `json:"count"`
		}{Count: count}

		return appresponse.ResponseSuccess(c, res)
	}, mymiddleware.VerifyJWT())

	handler.GET("/tracking/v1/ccu-logs/metrics", func(c echo.Context) error {
		latestLogTime := time.Now().Add(-time.Minute)

		countLimit, err := ccuLogRepo.CountCCU(latestLogTime, 30)
		if err != nil {
			return appresponse.ResponseError(c, err)
		}

		countSession, err := ccuLogRepo.CountCCU(latestLogTime, 5)
		if err != nil {
			return appresponse.ResponseError(c, err)
		}

		countPeak, err := ccuLogRepo.CountCCU(latestLogTime, 1)
		if err != nil {
			return appresponse.ResponseError(c, err)
		}

		data := `
isling_be_ccu_peak ` + strconv.Itoa(int(countPeak)) + `
isling_be_ccu_session ` + strconv.Itoa(int(countSession)) + `
isling_be_ccu_limit ` + strconv.Itoa(int(countLimit)) + `
		`

		return c.String(http.StatusOK, data)
	})

	return func() {
		if err := userActBatchUC.Stop(); err != nil {
			facade.Log().Error("stop userActBatchUC: %w", err)
		}

		if err := ccuLogBatchUC.Stop(); err != nil {
			facade.Log().Error("stop ccuLogBatchUC: %w", err)
		}
	}
}
