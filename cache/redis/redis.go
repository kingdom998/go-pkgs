package redis

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type Config struct {
	Addr         string
	Password     string
	Db           int64
	PoolSize     int32
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DialTimeout  time.Duration
}

func NewClient(conf *Config, logger log.Logger) *redis.Client {
	helper := log.NewHelper(log.With(logger, "module", "service/data/redis"))
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		DB:           int(conf.Db),
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		DialTimeout:  conf.DialTimeout,
		PoolSize:     int(conf.PoolSize),
	})
	cmd := rdb.Ping(context.Background())
	if cmd.Err() != nil {
		helper.Fatalf("ping redis error: %v", cmd.Err())
	}

	return rdb
}
