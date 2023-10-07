package playisling

import (
	"context"
	"encoding/json"
	"isling-be/config"
	acc_entity "isling-be/internal/account/entity"
	"isling-be/internal/common/controller/http/middleware"
	v1 "isling-be/internal/play-isling/controller/http/v1"
	"isling-be/internal/play-isling/controller/worker"
	"isling-be/internal/play-isling/entity"
	"isling-be/internal/play-isling/repo"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/logger"
	"isling-be/pkg/postgres"
	"strconv"

	"github.com/dgraph-io/ristretto"
	"github.com/labstack/echo/v4"
)

func Register(l logger.Interface, cfg *config.Config, cache *ristretto.Cache, handler *echo.Echo, pg *postgres.Postgres, msgBus *map[string]chan string) func() {
	protectedRoutes := handler.Group("", middleware.VerifyJWT())

	roomRepo := repo.NewRoomRepo(pg)
	playUserRepo := repo.NewPlayUserRepo(pg)

	roomUC := usecase.NewRoomUsecase(roomRepo, msgBus)
	homeUC := usecase.NewHomeUsecase(playUserRepo, roomRepo)
	recommendationUC := usecase.NewRecommendationUC()
	playUserUC := usecase.NewPlayUserUC(playUserRepo)

	roomRouter := v1.NewRoomRouter(roomUC)
	homeRouter := v1.NewHomeRouter(homeUC)
	trackingRouter := v1.NewTrackingRouter(recommendationUC, playUserUC)

	{
		protectedRoutes.POST("/play-isling/v1/rooms", roomRouter.Create)
		protectedRoutes.GET("/play-isling/v1/rooms", roomRouter.List)
		handler.GET("/play-isling/v1/rooms/:slugName", roomRouter.Show, middleware.ParseJWT())
		protectedRoutes.PATCH("/play-isling/v1/rooms/:id", roomRouter.Update)
		protectedRoutes.DELETE("/play-isling/v1/rooms/:id", roomRouter.Delete)

		protectedRoutes.GET("/play-isling/v1/home", homeRouter.Show)
		handler.GET("/play-isling/v1/guest/home", homeRouter.ShowGuest)

		protectedRoutes.POST("/play-isling/v1/actions", trackingRouter.Create)
	}

	userActChan := (*msgBus)["userActivityOnItem"]
	gorseETLWorker := worker.NewGorseETL(userActChan, recommendationUC)

	gorseETLWorker.Run()

	go func() {
		if msgBus == nil {
			return
		}

		accountCreatedChan, ok := (*msgBus)["accountCreated"]
		if !ok {
			return
		}

		for acc := range accountCreatedChan {
			account := new(acc_entity.Account)

			err := json.Unmarshal([]byte(acc), account)
			if err != nil {
				return
			}

			ctx := context.Background()

			err = recommendationUC.InsertUser(ctx, account)
			if err != nil {
				l.Error("insert user in recommendation: %w", err)
			}

			_, err = playUserRepo.Create(ctx, account.ID)
			if err != nil {
				l.Error("create play user: %w", err)
			}
		}
	}()

	go func() {
		if msgBus == nil {
			return
		}

		roomCreatedChan, ok := (*msgBus)["roomCreated"]
		if !ok {
			return
		}

		for acc := range roomCreatedChan {
			room := new(entity.Room)

			err := json.Unmarshal([]byte(acc), room)
			if err != nil {
				return
			}

			ctx := context.Background()

			err = recommendationUC.InsertRoom(ctx, room)
			if err != nil {
				l.Error("insert room in recommendation: %w", err)
			}
		}
	}()

	go func() {
		if msgBus == nil {
			return
		}

		roomDeletedChan, ok := (*msgBus)["roomDeleted"]
		if !ok {
			return
		}

		for acc := range roomDeletedChan {
			room := new(entity.Room)

			err := json.Unmarshal([]byte(acc), room)
			if err != nil {
				return
			}

			ctx := context.Background()

			err = recommendationUC.HideItem(ctx, strconv.FormatInt(room.ID, 10))
			if err != nil {
				l.Error("hide room in recommendation: %w", err)
			}
		}
	}()

	return func() {
		if err := recommendationUC.InsertFBBatch.Stop(); err != nil {
			l.Error("play-isling module: insertFBBatch stop: %w", err)
		}
	}
}
