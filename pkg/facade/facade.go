package facade

import (
	"isling-be/config"
	"isling-be/pkg/logger"
	"isling-be/pkg/redis"
	"isling-be/pkg/watermill"

	"github.com/dgraph-io/ristretto"
)

type Facade struct {
	log     logger.Interface
	config  *config.Config
	cache   *ristretto.Cache
	pubsub  *watermill.Watermill
	redis   *redis.Redis
	isSetup bool
}

var pool Facade

func Setup(log logger.Interface, cfg *config.Config, cache *ristretto.Cache) {
	if pool.isSetup {
		return
	}

	pool.isSetup = true

	pool.log = log
	pool.config = cfg
	pool.cache = cache
	pool.pubsub = watermill.NewWatermill(log)
	pool.redis = redis.NewRedis(cfg.Redis.URL, log)

	// Because router feature is not used yet, no need RunRouter
	// pool.messageBus.RunRouter()
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

func Pubsub() *watermill.Watermill {
	return pool.pubsub
}

func Redis() *redis.Redis {
	return pool.redis
}

func Close() {
	if !pool.isSetup {
		return
	}

	pool.redis.Close()
	pool.pubsub.Close()
}
