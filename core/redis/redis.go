package redis

import (
	"cookie_supply_management/core/config"
	"github.com/redis/go-redis/v9"
)

func New(conf config.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		DB:       conf.Database,
		Password: conf.Password,
	})
}
