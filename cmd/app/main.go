package main

import (
	"log"

	"github.com/btcs-longnp/isling-be/config"
	"github.com/btcs-longnp/isling-be/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
