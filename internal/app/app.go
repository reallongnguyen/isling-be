// Package app configures and runs application.
package app

import (
	"os"
	"os/signal"
	"syscall"

	"isling-be/config"
	"isling-be/pkg/facade"
	"isling-be/pkg/httpserver"
	"isling-be/pkg/logger"
	"isling-be/pkg/postgres"
	"isling-be/pkg/surreal"

	"github.com/dgraph-io/ristretto"
	"github.com/labstack/echo/v4"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	logPrettier := cfg.App.ENV == "development"
	l := logger.New(cfg.Log.Level, logPrettier)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, l, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal("app - Run - postgres.New: %v", err)
	}
	defer pg.Close()

	sur, err := surreal.New(
		cfg.Surreal.URL,
		cfg.Surreal.NS,
		cfg.Surreal.DB,
		cfg.Surreal.User,
		cfg.Surreal.Password,
		l,
	)
	if err != nil {
		l.Fatal("app - Run - surreal.New: %v", err)
	}
	defer sur.Close()

	// Cache
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1e8,
		BufferItems: 64,
	})
	if err != nil {
		l.Fatal("app - Run - cache.New: %v", err)
	}

	// Setup facade
	facade.Setup(l, cfg, cache)

	// HTTP Server
	handler := echo.New()
	configHTTPServer(handler)
	stopModules := useModules(handler, pg, sur)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	l.Info("app - Run - server started")

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify: %v", err)
	}

	// Shutdown
	stopModules()

	err = httpServer.Shutdown()
	if err != nil {
		l.Error("app - Run - httpServer.Shutdown: %v", err)
	}
}
