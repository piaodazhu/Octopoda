package rdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"

	"github.com/bluele/gcache"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client = nil
var cache gcache.Cache

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

	cache = gcache.New(256).LRU().Expiration(time.Minute).Build()
}

func GetString(key string) (result string, found bool, err error) {
	cachedValue, err := cache.Get(key)
	if err == nil {
		return cachedValue.(string), true, nil
	}

	result, err = rdb.Get(context.TODO(), key).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			found = false
		}
		return
	}
	found = true
	cache.Set(key, result)
	return
}

func SetStringXX(key, value string) (bool, error) {
	cache.Remove(key)
	defer time.AfterFunc(time.Millisecond*200, func() { cache.Remove(key) })
	return rdb.SetXX(context.TODO(), key, value, 0).Result()
}

func SetStringNX(key, value string) (bool, error) {
	cache.Remove(key)
	defer time.AfterFunc(time.Millisecond*200, func() { cache.Remove(key) })
	return rdb.SetNX(context.TODO(), key, value, 0).Result()
}

func DelKey(key string) error {
	cache.Remove(key)
	defer time.AfterFunc(time.Millisecond*200, func() { cache.Remove(key) })
	return rdb.Del(context.TODO(), key).Err()
}

func GetSMembers(key string) ([]string, error) {
	cachedValue, err := cache.Get(key)
	if err == nil {
		return cachedValue.([]string), nil
	}

	result, err := rdb.SMembers(context.TODO(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	cache.Set(key, result)
	return result, nil
}

func AddSMembers(key string, members ...string) error {
	// members must be valid names
	cache.Remove(key)
	defer time.AfterFunc(time.Millisecond*200, func() { cache.Remove(key) })
	return rdb.SAdd(context.TODO(), key, members).Err()
}

func RemoveSMembers(key string, members ...string) (int64, error) {
	// members must be subset of origin set
	cache.Remove(key)
	defer time.AfterFunc(time.Millisecond*200, func() { cache.Remove(key) })
	return rdb.SRem(context.TODO(), key, members).Result()
}

func CountSMembers(key string) (int64, error) {
	// members must be subset of origin set
	cachedValue, err := cache.Get(key)
	if err == nil {
		return cachedValue.(int64), nil
	}
	return rdb.SCard(context.TODO(), key).Result()
}
