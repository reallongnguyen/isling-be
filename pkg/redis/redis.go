package redis

import (
	"isling-be/pkg/logger"

	go_redis "github.com/redis/go-redis/v9"
)

const (
	_maxRetries = 10
)

type Redis struct {
	*go_redis.Client
}

func NewRedis(url string, log logger.Interface) *Redis {
	opt, err := go_redis.ParseURL(url)
	if err != nil {
		log.Fatal("parse redis url: %w", err)
	}

	opt.MaxRetries = _maxRetries

	client := go_redis.NewClient(opt)

	return &Redis{
		Client: client,
	}
}

func (r *Redis) Close() {
	if r.Client != nil {
		r.Client.Close()
	}
}
