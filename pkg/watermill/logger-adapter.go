package watermill

import (
	"isling-be/pkg/logger"

	"github.com/ThreeDotsLabs/watermill"
)

type LogAdapter struct {
	log    logger.Interface
	fields watermill.LogFields
}

var _ watermill.LoggerAdapter = (*LogAdapter)(nil)

func (r *LogAdapter) Info(msg string, logFields watermill.LogFields) {
	r.log.Info("%s %v", msg, logFields)
}

func (r *LogAdapter) Debug(msg string, logFields watermill.LogFields) {
	r.log.Debug("%s %v", msg, logFields)
}

func (r *LogAdapter) Trace(msg string, logFields watermill.LogFields) {
	r.log.Trace("%s %v", msg, logFields)
}

func (r *LogAdapter) Error(msg string, err error, logFields watermill.LogFields) {
	r.log.Error("%s %w %v", msg, err, logFields)
}

func (r *LogAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	r.fields.Add(fields)

	return r
}
