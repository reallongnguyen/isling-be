package playisling

import (
	"context"
	"encoding/json"
	acc_entity "isling-be/internal/account/entity"
	"isling-be/internal/common/controller/http/middleware"
	v1 "isling-be/internal/play-isling/controller/http/v1"
	"isling-be/internal/play-isling/controller/worker"
	"isling-be/internal/play-isling/entity"
	"isling-be/internal/play-isling/repo"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/facade"
	"isling-be/pkg/postgres"
	"isling-be/pkg/surreal"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Register(handler *echo.Echo, pg *postgres.Postgres, sur *surreal.Surreal) func() {
	protectedRoutes := handler.Group("", middleware.VerifyJWT())

	roomRepo := repo.NewRoomRepo(pg)
	playUserRepo := repo.NewPlayUserRepo(pg)

	roomUC := usecase.NewRoomUsecase(roomRepo)
	homeUC := usecase.NewHomeUsecase(playUserRepo, roomRepo)
	recommendationUC := usecase.NewRecommendationUC()
	playUserUC := usecase.NewPlayUserUC(playUserRepo)

	roomRouter := v1.NewRoomRouter(roomUC)
	homeRouter := v1.NewHomeRouter(homeUC)

	{
		protectedRoutes.POST("/play-isling/v1/rooms", roomRouter.Create)
		protectedRoutes.GET("/play-isling/v1/rooms", roomRouter.List)
		handler.GET("/play-isling/v1/rooms/:slugName", roomRouter.Show, middleware.ParseJWT())
		protectedRoutes.PATCH("/play-isling/v1/rooms/:id", roomRouter.Update)
		protectedRoutes.DELETE("/play-isling/v1/rooms/:id", roomRouter.Delete)

		protectedRoutes.GET("/play-isling/v1/home", homeRouter.Show)
		handler.GET("/play-isling/v1/guest/home", homeRouter.ShowGuest)
	}

	gorseETLWorker := worker.NewGorseETL(recommendationUC)
	gorseETLWorker.Run()

	roomAudCounter := worker.NewRoomAudCounter(sur)
	roomAudCounter.Run()

	watchAgainUpdater := worker.NewWatchAgainUpdater(playUserUC)
	watchAgainUpdater.Run()

	// TODO: move to /controller/worker
	go func() {
		handler := func(uuid string, payload []byte, metadata map[string]string) error {
			account := new(acc_entity.Account)

			err := json.Unmarshal(payload, account)
			if err != nil {
				return nil
			}

			ctx := context.Background()

			_, err = playUserRepo.Create(ctx, account.ID)
			if err != nil {
				facade.Log().Error("create play user: %w", err)

				return err
			}

			return nil
		}

		err := facade.Pubsub().Subscribe("account.created", handler)
		if err != nil {
			facade.Log().Error("subscribe topic 'account.created' error %w", err)
		}
	}()

	// TODO: move to /controller/worker
	go func() {
		handler := func(uuid string, payload []byte, metadata map[string]string) error {
			room := new(entity.Room)

			err := json.Unmarshal(payload, room)
			if err != nil {
				return nil
			}

			ctx := context.Background()

			err = recommendationUC.HideItem(ctx, strconv.FormatInt(room.ID, 10))
			if err != nil {
				facade.Log().Error("hide room in recommendation: %w", err)
			}

			return nil
		}

		err := facade.Pubsub().Subscribe("room.deleted", handler)
		if err != nil {
			facade.Log().Error("subscribe topic 'room.deleted' error %w", err)
		}
	}()

	return func() {
		if err := recommendationUC.InsertFBBatch.Stop(); err != nil {
			facade.Log().Error("play-isling module: insertFBBatch stop: %w", err)
		}
	}
}
