package rdb

import (
	"brain/logger"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	TaskFinished = iota
	TaskFailed
	TaskProcessing
	TaskNotFound
	TaskDbError
)

// key format: "taskid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// value format: "{code}{results...}"
func TaskIdGen() string {
	return uuid.New().String()
}

func TaskNew(taskid string, timeout_second int) bool {
	if timeout_second == 0 {
		timeout_second = 86400
	}
	err := rdb.Set(context.TODO(), "taskid:"+taskid, TaskProcessing, time.Second*time.Duration(timeout_second)).Err()
	if err != nil {
		logger.Brain.Println(err)
		return false
	}
	return true
}

func TaskMarkDone(taskid string, result interface{}, timeout_second int) bool {
	res_serialized, _ := json.Marshal(result)
	value := fmt.Sprintf("%d%s", TaskFinished, string(res_serialized))
	return taskMark(taskid, value, timeout_second)
}

func TaskMarkFailed(taskid string, result interface{}, timeout_second int) bool {
	res_serialized, _ := json.Marshal(result)
	value := fmt.Sprintf("%d%s", TaskFailed, string(res_serialized))
	return taskMark(taskid, value, timeout_second)
}

func taskMark(taskid string, value string, timeout_second int) bool {
	// get remaining ttl
	rem := rdb.TTL(context.TODO(), "taskid:"+taskid).Val().Nanoseconds()
	timeout := time.Second*time.Duration(timeout_second)

	// if rem == -2, the key is not exists, return false
	// if rem > timeout_second, we keep the finished task record as it was originally set
	// else, we keep the finished task for timeout_second seconds.
	if rem == -2 {
		return false
	} else if rem > int64(time.Second*time.Duration(timeout_second)) {
		timeout = time.Duration(rem)
	} 
	
	err := rdb.Set(context.TODO(), "taskid:"+taskid, value, timeout).Err()
	if err != nil {
		logger.Brain.Println(err)
		return false
	}
	return true
}


func TaskDelete(taskid string) bool {
	err := rdb.Del(context.TODO(), "taskid:"+taskid).Err()
	if err != nil {
		logger.Brain.Println(err)
		return false
	}
	return true
}

func TaskQuery(taskid string) (int, string) {
	value, err := rdb.Get(context.TODO(), "taskid:"+taskid).Result()
	if err == redis.Nil {
		return TaskNotFound, ""
	} else if err != nil || len(value) == 0 {
		logger.Brain.Println(err)
		return TaskDbError, ""
	}
	
	return int(value[0] - '0'), value[1:]
}
