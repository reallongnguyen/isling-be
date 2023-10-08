package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

// Interface -.
type Interface interface {
	Trace(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
	Fatal(message string, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(level string, prettier bool) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	case "trace":
		l = zerolog.TraceLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	logger := zlog.
		With().
		Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
		Logger()

	if prettier {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return &Logger{
		logger: &logger,
	}
}

// Debug -.
func (l *Logger) Trace(message string, args ...interface{}) {
	l.logger.Trace().Msgf(message, args...)
}

// Debug -.
func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Debug().Msgf(message, args...)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Info().Msgf(message, args...)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Warn().Msgf(message, args...)
}

// Error -.
func (l *Logger) Error(message string, args ...interface{}) {
	l.logger.Error().Msgf(message, args...)
}

// Fatal -.
func (l *Logger) Fatal(message string, args ...interface{}) {
	l.logger.Fatal().Msgf(message, args...)
}
