package rdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"

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

func GetString(key string) (result string, found bool, err error) {
	result, err = rdb.Get(context.TODO(), key).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			found = false
		}
		return
	}
	found = true
	return
}

func SetStringXX(key, value string) (bool, error) {
	return rdb.SetXX(context.TODO(), key, value, 0).Result()
}

func SetStringNX(key, value string) (bool, error) {
	return rdb.SetNX(context.TODO(), key, value, 0).Result()
}

func DelKey(key string) error {
	return rdb.Del(context.TODO(), key).Err()
}

func GetSMembers(key string) ([]string, error) {
	result, err := rdb.SMembers(context.TODO(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func AddSMembers(key string, members ...string) error {
	// members must be valid names
	return rdb.SAdd(context.TODO(), key, members).Err()
}

func RemoveSMembers(key string, members ...string) (int64, error) {
	// members must be subset of origin set
	return rdb.SRem(context.TODO(), key, members).Result()
}

func CountSMembers(key string) (int64, error) {
	// members must be subset of origin set
	return rdb.SCard(context.TODO(), key).Result()
}