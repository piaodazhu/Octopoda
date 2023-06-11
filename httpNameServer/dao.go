package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type NameEntryDao struct {
	DefaultTTL int
}

var rdb *redis.Client
var ctx context.Context
var nedao *NameEntryDao

func DaoInit() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	ctx = context.TODO()
	nedao = &NameEntryDao{DefaultTTL: 5000}
	return rdb.Ping(ctx).Err()
}

func GetNameEntryDao() *NameEntryDao {
	return nedao
}

func (n *NameEntryDao) DaoSet(key string, entry NameEntry, ttl int) error {
	if ttl == 0 {
		ttl = n.DefaultTTL
	}

	err := rdb.HSet(ctx, key, entry).Err()
	if err != nil {
		return err
	}

	if ttl > 0 {
		err = rdb.Expire(ctx, key, time.Millisecond*time.Duration(ttl)).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *NameEntryDao) DaoGet(key string) (*NameEntry, error) {
	res := &NameEntry{}
	err := rdb.HGetAll(ctx, key).Scan(res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NameEntryDao) DaoDel(key string) error {
	return rdb.Del(ctx, key).Err()
}

func (n *NameEntryDao) DaoList(pattern string) ([]string, error) {
	return rdb.Keys(ctx, pattern).Result()
}
