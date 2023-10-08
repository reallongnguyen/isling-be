package pg_surreal_sync

import (
	"context"
	"encoding/json"
	"isling-be/internal/tool/model"
	"isling-be/pkg/facade"
	"isling-be/pkg/postgres"
	"isling-be/pkg/surreal"
)

type PGSurrealSync struct {
	pg *postgres.Postgres
	sr *surreal.Surreal
}

var _ model.Tool = (*PGSurrealSync)(nil)

func NewPGSurrealSync(pg *postgres.Postgres, sr *surreal.Surreal) *PGSurrealSync {
	return &PGSurrealSync{
		pg: pg,
		sr: sr,
	}
}

func (r *PGSurrealSync) Start() error {
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

	for {
		noti, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			facade.Log().Error("PGsurrealSync: wait notification: %w", err)

			continue
		}

		payload := new(Payload)

		if err = json.Unmarshal([]byte(noti.Payload), payload); err != nil {
			continue
		}

		err = syncDataUC.Handle(payload)
		if err != nil {
			facade.Log().Error("PGsurrealSync: handle message: %w", err)
		}
	}
}
