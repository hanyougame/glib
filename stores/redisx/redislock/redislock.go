package redislock

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

// New 获取锁
func New(rdb redis.UniversalClient) *redsync.Redsync {
	return redsync.New(goredis.NewPool(rdb))
}
