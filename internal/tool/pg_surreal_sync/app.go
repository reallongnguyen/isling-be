package pg_surreal_sync

import (
	"context"
	"encoding/json"
	"isling-be/internal/tool/entity"
	"isling-be/pkg/facade"
	"isling-be/pkg/postgres"
	"isling-be/pkg/surreal"
	"time"
)

type PGSurrealSync struct {
	pg *postgres.Postgres
	sr *surreal.Surreal
}

var _ entity.Tool = (*PGSurrealSync)(nil)

func NewPGSurrealSync(pg *postgres.Postgres, sr *surreal.Surreal) *PGSurrealSync {
	return &PGSurrealSync{
		pg: pg,
		sr: sr,
	}
}

func (r *PGSurrealSync) Start() {
	go func() {
		for {
			r.Sync()

			time.Sleep(time.Second)
		}
	}()
}

func (r *PGSurrealSync) Sync() error {
	conn, err := r.pg.Pool.Acquire(context.Background())
	if err != nil {
		facade.Log().Error("PGsurrealSync: acquire connection: %w", err)

		return err
	}

	defer conn.Release()

	_, err = conn.Exec(context.Background(), "LISTEN table_update")
	if err != nil {
		facade.Log().Error("PGsurrealSync: listen event: %w", err)

		return err
	}

	syncDataUC := NewSyncDataUsecase(r.sr)

	facade.Log().Info("PGsurrealSync: start listen postgres then sync to surreal")

	for {
		noti, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			facade.Log().Error("PGsurrealSync: wait notification: %w", err)

			return err
		}

		payload := new(Payload)

		if err = json.Unmarshal([]byte(noti.Payload), payload); err != nil {
			continue
		}

		facade.Log().Trace("PGsurrealSync: receive msg: %v", *payload)

		err = syncDataUC.Handle(payload)
		if err != nil {
			facade.Log().Error("PGsurrealSync: handle message: %w", err)
		}
	}
}
