package tool

import (
	"isling-be/internal/tool/pg_surreal_sync"
	"isling-be/pkg/postgres"
	"isling-be/pkg/surreal"
)

func Register(pg *postgres.Postgres, sr *surreal.Surreal) {
	pgSurrealSync := pg_surreal_sync.NewPGSurrealSync(pg, sr)

	go func() {
		pgSurrealSync.Start()
	}()
}
