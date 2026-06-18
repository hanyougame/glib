package redisx

import (
	"context"
	"time"

	"github.com/hanyougame/glib/stores/redisx/config"
	"github.com/redis/go-redis/v9"
)

var Engine redis.UniversalClient

func Must(c config.Config) {
	Engine = NewEngine(c)
}

func NewEngine(c config.Config) (rdb redis.UniversalClient) {
	// 连接池配置，0 值交给 go-redis 的 Options.init() 填充默认值。
	dialTimeout := time.Duration(c.DialTimeoutMs) * time.Millisecond
	readTimeout := time.Duration(c.ReadTimeoutMs) * time.Millisecond
	writeTimeout := time.Duration(c.WriteTimeoutMs) * time.Millisecond
	poolTimeout := time.Duration(c.PoolTimeoutMs) * time.Millisecond

	if c.IsCluster {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        c.Addrs,
			Username:     c.Username,
			Password:     c.Password,
			PoolSize:     c.PoolSize,
			MinIdleConns: c.MinIdleConns,
			DialTimeout:  dialTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			PoolTimeout:  poolTimeout,
		})
	} else {
		rdb = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:        c.Addrs,
			Username:     c.Username,
			Password:     c.Password,
			MasterName:   c.MasterName,
			DB:           c.DB,
			PoolSize:     c.PoolSize,
			MinIdleConns: c.MinIdleConns,
			DialTimeout:  dialTimeout,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			PoolTimeout:  poolTimeout,
		})
	}

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	if c.Debug {
		rdb.AddHook(DebugHook{})
	}

	if c.Trace {
		rdb.AddHook(TraceHook{})
	}

	return rdb
}
