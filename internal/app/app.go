// Package app configures and runs application.
package app

import (
	"fmt"
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
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, l, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	sur, err := surreal.New(
		cfg.Surreal.URL,
		cfg.Surreal.NS,
		cfg.Surreal.DB,
		cfg.Surreal.User,
		cfg.Surreal.Password,
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - surreal.New: %w", err))
	}
	defer sur.Close()

	// Cache
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1e8,
		BufferItems: 64,
	})
	if err != nil {
		l.Fatal("app - Run - cache.New: %w", err)
	}

	// Msg bus
	msgBus := make(map[string]chan string)

	msgBus["accountCreated"] = make(chan string)
	defer close(msgBus["accountCreated"])

	msgBus["roomCreated"] = make(chan string)
	defer close(msgBus["accountCreated"])

	msgBus["roomDeleted"] = make(chan string)
	defer close(msgBus["roomDeleted"])

	msgBus["userActivityOnItem"] = make(chan string)
	defer close(msgBus["userActivityOnItem"])

	// Setup facade
	facade.Setup(l, cfg, cache)

	// HTTP Server
	handler := echo.New()
	configHTTPServer(cfg, handler)
	stopModules := useModules(l, cache, cfg, handler, pg, sur, &msgBus)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	l.Info("app - Run - server started")

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	stopModules()

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
