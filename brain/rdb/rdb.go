package rdb

import (
	"brain/config"
	"brain/logger"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client = nil

func InitRedis() {
	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.Redis.Ip)
	sb.WriteByte(':')
	sb.WriteString(fmt.Sprint(config.GlobalConfig.Redis.Port))
	Rdb = redis.NewClient(&redis.Options{
		Addr:     sb.String(),
		DB:       config.GlobalConfig.Redis.Db,
		Password: config.GlobalConfig.Redis.Password,
	})

	if Rdb.Ping(context.TODO()).Err() != nil {
		logger.Tentacle.Panic("Cannot connect to redis")
	}
}

func StoreNode(nodename string, ip string, port uint16) error {
	var sb strings.Builder
	sb.WriteString(ip)
	sb.WriteByte(':')
	sb.WriteString(strconv.Itoa(int(port)))
	err := Rdb.HSet(context.TODO(), nodename, "addr", sb.String(), "stat", "active").Err()
	if err != nil {
		return err
	}

	err = Rdb.Expire(context.TODO(), nodename, time.Duration(config.GlobalConfig.TentacleFace.RecordTimeout)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

func UpdateNode(nodename string) bool {
	Rdb.ExpireXX(context.TODO(), nodename, time.Duration(config.GlobalConfig.TentacleFace.RecordTimeout)*time.Second)
	exists, _ := Rdb.Exists(context.TODO(), nodename).Result()
	if exists == 0 {
		return false
	}

	err := Rdb.HSet(context.TODO(), nodename, "stat", "active").Err()
	return err == nil
}

func GetNodeAddress(nodename string) string {
	cmd := Rdb.HGet(context.TODO(), nodename, "addr")
	if cmd.Err() != nil {
		return ""
	}
	res, _ := cmd.Result()
	return res
}

func GetNodeState(nodename string) string {
	cmd := Rdb.HGet(context.TODO(), nodename, "stat")
	if cmd.Err() != nil {
		return ""
	}
	res, _ := cmd.Result()
	return res
}
