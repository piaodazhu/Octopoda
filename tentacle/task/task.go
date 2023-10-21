package task

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"tentacle/message"
	"time"

	"github.com/google/uuid"
)

const (
	StateInvalid   = "Invalid"
	StatePending   = "Pending"
	StateFinished  = "Finish"
	StateCancelled = "Canceled"
)

type TaskInfo struct {
	TaskId      string
	State       string
	Brief       string
	CreatedTime time.Time
}

type Task struct {
	TaskInfo

	result chan *message.Result
	cancelFunc func()
	ctx        context.Context
}

type UserTaskFunc func() *message.Result
type UserCancelFunc func()

func emptyFunc() {}

type taskManager struct {
	PendingTasks, FinishTasks sync.Map
}

var TaskManager taskManager

// key format: "taskid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// value format: "{code}{results...}"
func (tm *taskManager) taskIdGen() string {
	return uuid.New().String()
}

func (tm *taskManager) CreateTask(brief string, utask UserTaskFunc, ucancel UserCancelFunc) (string, error) {
	if utask == nil {
		return "", errors.New("empty user task is not allowed")
	}
	if ucancel == nil {
		ucancel = emptyFunc
	}

	newId := tm.taskIdGen()
	newTask := &Task{
		TaskInfo: TaskInfo{
			TaskId:      newId,
			Brief:       brief,
			State:       StatePending,
			CreatedTime: time.Now(),
		},
		result: make(chan *message.Result, 1),
	}

	for { // 如果Id冲突，则必须重新生成
		if _, crashed := tm.PendingTasks.LoadOrStore(newId, newTask); !crashed {
			// finish里面也不能有冲突
			if _, crashed := tm.FinishTasks.Load(newId); !crashed {
				break
			} else {
				// Pending里面存储了，但是发现在load里面有重复，则删除已存储的，重新生成id
				tm.PendingTasks.Delete(newId)
			}
		}
		newId = tm.taskIdGen()
		newTask.TaskId = newId
	}

	ctx, ccancel := context.WithCancel(context.Background())
	newTask.ctx = ctx
	newTask.cancelFunc = func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recover from ", err)
			}
		}()
		ccancel()
		ucancel()
	}

	result := make(chan *message.Result, 1)
	go func() {
		res := utask() // utask必须是用户可cancel的，否则这个goroutine无法结束
		result <- res
		close(result)
	}()

	// 上下两个goroutine的关系备注：下面的goroutine退出要么utask执行完，要么ctx调用了取消，都一定会意味着上面的goroutine能正常退出
	go func() {
		var res *message.Result
		select {
		case res = <-result:
			// done
		case <-ctx.Done():
			// cancelled
			res = nil
		}
		value, _ := tm.PendingTasks.LoadAndDelete(newId) // 必定存在！
		task := value.(*Task)
		if res == nil {
			task.State = StateCancelled
		} else {
			task.State = StateFinished
		}
		task.result <- res
		tm.FinishTasks.Store(newId, task)
	}()

	return newId, nil
}

func (tm *taskManager) CancelTask(taskId string) bool {
	if value, found := tm.PendingTasks.Load(taskId); found {
		task := value.(*Task)
		task.cancelFunc()
	}
	return false
}

// 阻塞等待结果，如果任务被取消则返回false
func (tm *taskManager) WaitTask(taskId string) (*message.Result, bool) {
	var task *Task
	if value, found := tm.FinishTasks.Load(taskId); found {
		task = value.(*Task)
	} else if value, found := tm.PendingTasks.Load(taskId); found {
		task = value.(*Task)
	} else {
		return nil, false 
	}
	res := <- task.result
	tm.FinishTasks.Delete(taskId)
	
	return res, true
}


// -- READ-ONLY --

func (tm *taskManager) QueryTask(taskId string) (TaskInfo, bool) {
	if value, found := tm.FinishTasks.Load(taskId); found {
		task := value.(*Task)
		return task.TaskInfo, true
	}
	if value, found := tm.PendingTasks.Load(taskId); found {
		task := value.(*Task)
		return task.TaskInfo, true
	}

	// 无效的taskId
	tinfo := TaskInfo{}
	tinfo.TaskId = taskId
	tinfo.State = StateInvalid
	return tinfo, false
}

func (tm *taskManager) ListTasks() []TaskInfo {
	taskSet := map[string]TaskInfo{}
	tm.PendingTasks.Range(func(key, value any) bool {
		task := value.(*Task)
		taskSet[task.TaskId] = task.TaskInfo
		return true
	})
	// 先看pending，再看finish，如果期间有pending的task被完成，则会直接覆盖
	tm.FinishTasks.Range(func(key, value any) bool {
		task := value.(*Task)
		taskSet[task.TaskId] = task.TaskInfo
		return true
	})

	// 输出的列表按照taskid去重，按创建任务的时间排序
	taskList := []TaskInfo{}
	for _, taskState := range taskSet {
		taskList = append(taskList, taskState)
	}
	sort.Slice(taskList, func(i, j int) bool {
		return taskList[i].CreatedTime.Before(taskList[j].CreatedTime)
	})

	return taskList
}
