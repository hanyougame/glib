package redislock

import (
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

func New(rdb redis.UniversalClient) *redislock.Client {
	return redislock.New(rdb)
}
