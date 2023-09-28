package tool

import (
	"isling-be/internal/tool/pg_surreal_sync"
	"isling-be/internal/tool/room_audience"
	"isling-be/pkg/logger"
	"isling-be/pkg/postgres"
	"isling-be/pkg/surreal"
)

func Register(l logger.Interface, pg *postgres.Postgres, sr *surreal.Surreal) {
	pgSurrealSync := pg_surreal_sync.NewPGSurrealSync(l, pg, sr)
	roomAudience := room_audience.NewRoomAudience(l, sr)

	go func() {
		pgSurrealSync.Start()
	}()

	go func() {
		roomAudience.Start()
	}()
}
