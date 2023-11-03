package rdb

import (
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"context"
)

func GroupAdd(name string, nodes []string) bool {
	value, _ := config.Jsoner.Marshal(nodes)
	pipe := rdb.TxPipeline()
	pipe.Del(context.Background(), name)
	pipe.Set(context.TODO(), "group:"+name, string(value), 0)

	_, err := pipe.Exec(context.TODO())
	if err != nil {
		logger.Exceptions.Println(err)
		return false
	}
	return true
}

func GroupGet(name string) ([]string, bool) {
	var nodes []string
	value, err := rdb.Get(context.TODO(), "group:"+name).Result()
	if err != nil {
		logger.Exceptions.Println(err)
		return nil, false
	}

	err = config.Jsoner.Unmarshal([]byte(value), &nodes)
	if err != nil {
		logger.Exceptions.Println(err)
		return nil, false
	}

	return nodes, true
}

func GroupDel(name string) bool {
	err := rdb.Del(context.TODO(), "group:"+name).Err()
	if err != nil {
		logger.Exceptions.Println(err)
		return false
	}
	return true
}

func GroupExist(name string) bool {
	cmd := rdb.Exists(context.TODO(), "group:"+name)
	return cmd.Val() != 0
}

func GroupGetAll() []string {
	keys := rdb.Keys(context.TODO(), "group:*").Val()
	groups := []string{}
	for _, key := range keys {
		groups = append(groups, key[len("group:"):])
	}
	return groups
}
