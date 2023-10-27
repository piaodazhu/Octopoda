package rdb

import (
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client = nil

func InitRedis() {
	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.Redis.Ip)
	sb.WriteByte(':')
	sb.WriteString(fmt.Sprint(config.GlobalConfig.Redis.Port))
	rdb = redis.NewClient(&redis.Options{
		Addr:     sb.String(),
		DB:       config.GlobalConfig.Redis.Db,
		Password: config.GlobalConfig.Redis.Password,
	})

	if rdb.Ping(context.TODO()).Err() != nil {
		logger.Exceptions.Panic("Cannot connect to redis")
	}
}
