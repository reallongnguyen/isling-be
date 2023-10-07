package facade

import (
	"isling-be/config"
	"isling-be/pkg/logger"

	"github.com/dgraph-io/ristretto"
)

type Facade struct {
	log    logger.Interface
	config *config.Config
	cache  *ristretto.Cache
}

var pool Facade

func Setup(log logger.Interface, cfg *config.Config, cache *ristretto.Cache) {
	pool.log = log
	pool.config = cfg
	pool.cache = cache
}

func Log() logger.Interface {
	return pool.log
}

func Config() *config.Config {
	return pool.config
}

func Cache() *ristretto.Cache {
	return pool.cache
}
