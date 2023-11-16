package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/piaodazhu/Octopoda/httpns/config"
	"github.com/piaodazhu/Octopoda/protocols"

	"github.com/bluele/gcache"
	"github.com/redis/go-redis/v9"
)

type BaseDao struct {
	DefaultPrefix string
	DefaultTTL    int
	Cache         gcache.Cache
}

type NameEntryDao struct {
	BaseDao
}

var rdb *redis.Client
var ctx context.Context
var namedao *NameEntryDao

func DaoInit() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.GlobalConfig.Redis.Ip, config.GlobalConfig.Redis.Port),
		Password: config.GlobalConfig.Redis.Password,
		DB:       config.GlobalConfig.Redis.Db,
	})
	ctx = context.TODO()

	namedao = &NameEntryDao{
		BaseDao: BaseDao{
			DefaultTTL:    5000,
			DefaultPrefix: "NameEntry:",
			Cache:         gcache.New(256).Build(),
		},
	}
	return rdb.Ping(ctx).Err()
}

func GetNameEntryDao() *NameEntryDao {
	return namedao
}

func (b *BaseDao) list(pattern string) ([]string, error) {
	pattern = b.DefaultPrefix + "*" + pattern
	keys, err := rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	res := make([]string, len(keys))
	for i := range keys {
		res[i] = keys[i][len(b.DefaultPrefix):]
	}
	return res, nil
}

func (b *BaseDao) del(key string) error {
	key = b.DefaultPrefix + key
	defer b.Cache.Remove(key)
	if err := rdb.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

func (b *BaseDao) set(key string, value string, ttl int) error {
	key = b.DefaultPrefix + key
	defer b.Cache.Remove(key)
	if ttl == 0 {
		ttl = b.DefaultTTL
	} else if ttl < 0 {
		ttl = 0
	}
	if err := rdb.Set(ctx, key, value, time.Millisecond*time.Duration(ttl)).Err(); err != nil {
		return err
	}
	return nil
}

func (b *BaseDao) get(key string) (string, error) {
	key = b.DefaultPrefix + key
	val, err := b.Cache.Get(key)
	if err == nil {
		value := val.(string)
		if len(value) == 0 {
			return "", errors.New("value not found")
		}
		return val.(string), nil
	}

	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		b.Cache.SetWithExpire(key, "", time.Second)
		return "", err
	}
	b.Cache.SetWithExpire(key, value, time.Second)
	return value, nil
}

// NameEntryDao

func (n *NameEntryDao) Set(key string, entry protocols.NameServiceEntry, ttl int) error {
	raw, _ := json.Marshal(entry)
	return n.set(key, string(raw), ttl)
}

func (n *NameEntryDao) Get(key string) (*protocols.NameServiceEntry, error) {
	value, err := n.get(key)
	if err != nil {
		return nil, err
	}
	res := &protocols.NameServiceEntry{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NameEntryDao) Del(key string) error {
	return n.del(key)
}

func (n *NameEntryDao) List(pattern string) ([]string, error) {
	return n.list(pattern)
}
