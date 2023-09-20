package main

import (
	"context"
	"encoding/json"
	"errors"
	"time"

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
type ConfigDao struct {
	BaseDao
}
type SshInfoDao struct {
	BaseDao
}

var rdb *redis.Client
var ctx context.Context
var namedao *NameEntryDao
var confdao *ConfigDao
var sshdao *SshInfoDao

func DaoInit() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	ctx = context.TODO()
	namedao = &NameEntryDao{
		BaseDao: BaseDao{
			DefaultTTL:    5000,
			DefaultPrefix: "NameEntry:",
			Cache:         gcache.New(256).Build(),
		},
	}
	confdao = &ConfigDao{
		BaseDao: BaseDao{
			DefaultTTL:    0,
			DefaultPrefix: "Config:",
			Cache:         gcache.New(256).Build(),
		},
	}
	sshdao = &SshInfoDao{
		BaseDao: BaseDao{
			DefaultTTL:    0,
			DefaultPrefix: "SshInfo:",
			Cache:         gcache.New(256).Build(),
		},
	}
	return rdb.Ping(ctx).Err()
}

func GetNameEntryDao() *NameEntryDao {
	return namedao
}
func GetNameConfigDao() *ConfigDao {
	return confdao
}
func GetSshInfoDao() *SshInfoDao {
	return sshdao
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
		b.Cache.Set(key, "")
		return "", err
	}
	b.Cache.Set(key, value)
	return value, nil
}

func (b *BaseDao) append(key string, value string) error {
	key = b.DefaultPrefix + key
	return rdb.LPush(ctx, key, value).Err()
}

func (b *BaseDao) getrange(key string, index, amount int) ([]string, error) {
	key = b.DefaultPrefix + key
	return rdb.LRange(ctx, key, int64(index), int64(index+amount)).Result()
}

// NameEntryDao

func (n *NameEntryDao) Set(key string, entry NameEntry, ttl int) error {
	raw, _ := json.Marshal(entry)
	return n.set(key, string(raw), ttl)
}

func (n *NameEntryDao) Get(key string) (*NameEntry, error) {
	value, err := n.get(key)
	if err != nil {
		return nil, err
	}
	res := &NameEntry{}
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

// ConfigDao

func (c *ConfigDao) Append(key string, config ConfigEntry) error {
	raw, _ := json.Marshal(config)
	return c.append(key, string(raw))
}

func (c *ConfigDao) GetRange(key string, index, amount int) ([]*ConfigEntry, error) {
	value, err := c.getrange(key, index, amount)
	if err != nil {
		return nil, err
	}
	res := []*ConfigEntry{}
	for _, conf := range value {
		var item ConfigEntry
		err = json.Unmarshal([]byte(conf), &item)
		if err != nil {
			return nil, err
		}
		res = append(res, &item)
	}
	return res, nil
}

func (c *ConfigDao) Del(key string) error {
	return c.del(key)
}

func (c *ConfigDao) List(pattern string) ([]string, error) {
	return c.list(pattern)
}

// SshInfoDao

func (s *SshInfoDao) Set(key string, ssh SshInfo, ttl int) error {
	raw, _ := json.Marshal(ssh)
	return s.set(key, string(raw), ttl)
}

func (s *SshInfoDao) Get(key string) (*SshInfo, error) {
	value, err := s.get(key)
	if err != nil {
		return nil, err
	}
	res := &SshInfo{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *SshInfoDao) Del(key string) error {
	return s.del(key)
}

func (s *SshInfoDao) List(pattern string) ([]string, error) {
	return s.list(pattern)
}
