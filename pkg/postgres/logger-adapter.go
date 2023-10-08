package postgres

import (
	"context"
	"isling-be/pkg/logger"

	"github.com/jackc/pgx/v4"
)

type Logger struct {
	logger.Interface
}

func (r *Logger) Log(_ context.Context, _ pgx.LogLevel, _ string, data map[string]interface{}) {
	sql, ok := data["sql"]
	if !ok {
		return
	}

	r.Trace(sql.(string))
}
