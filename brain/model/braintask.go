package model

import (
	"sync"

	"github.com/google/uuid"
)

// key format: "taskid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
// value format: "{code}{results...}"
func TaskIdGen() string {
	return uuid.New().String()
}

type brainSubTask struct {
	subTaskId string
	isDone    bool
	result    interface{}
}

type brainTask struct {
	taskId   string
	isDone   bool
	size     int
	finish   int
	subTasks []brainSubTask

	lock sync.Mutex
}

type brainTaskManager struct {
	tasks sync.Map
}

var BrainTaskManager brainTaskManager

func (m *brainTaskManager) CreateTask(subtask_num int) string {
	for {
		taskid := TaskIdGen()
		task := &brainTask{
			taskId:   taskid,
			isDone:   false,
			size:     subtask_num,
			finish:   0,
			subTasks: nil,
			lock:     sync.Mutex{},
		}
		if _, crashed := m.tasks.LoadOrStore(taskid, task); !crashed {
			return taskid
		}
		// taskid必须唯一，否则重新生成
	}
}

func (m *brainTaskManager) AddSubTask(task_id string, subtask_id string, result interface{}) bool {
	value, found := m.tasks.Load(task_id)
	if !found {
		return false
	}
	bt := value.(*brainTask)
	bt.lock.Lock()
	defer bt.lock.Unlock()
	bt.subTasks = append(bt.subTasks, brainSubTask{
		subTaskId: subtask_id,
		isDone:    false,
		result:    result,
	})
	return true
}

func (m *brainTaskManager) AddFailedSubTask(task_id string, subtask_id string, result interface{}) bool {
	value, found := m.tasks.Load(task_id)
	if !found {
		return false
	}
	bt := value.(*brainTask)
	bt.lock.Lock()
	defer bt.lock.Unlock()
	bt.subTasks = append(bt.subTasks, brainSubTask{
		subTaskId: subtask_id,
		isDone:    true,
		result:    result,
	})
	bt.finish++
	if bt.finish == bt.size {
		bt.isDone = true
	}
	return true
}

func (m *brainTaskManager) DoneSubTask(task_id string, subtask_id string, result interface{}) bool {
	value, found := m.tasks.Load(task_id)
	if !found {
		return false
	}
	bt := value.(*brainTask)
	bt.lock.Lock()
	defer bt.lock.Unlock()
	for i := range bt.subTasks {
		if bt.subTasks[i].subTaskId == subtask_id {
			bt.subTasks[i].isDone = true
			bt.subTasks[i].result = result
			bt.finish++
		}
	}
	if bt.finish == bt.size {
		bt.isDone = true
	}
	return true
}

type PendingSubTask struct {
	SubTaskId string
	Result    interface{}
}

func (m *brainTaskManager) PendingList(task_id string) []PendingSubTask {
	value, found := m.tasks.Load(task_id)
	if !found {
		return nil
	}
	bt := value.(*brainTask)
	bt.lock.Lock()
	defer bt.lock.Unlock()

	plist := []PendingSubTask{}
	for i := range bt.subTasks {
		if !bt.subTasks[i].isDone {
			plist = append(plist, PendingSubTask{
				SubTaskId: bt.subTasks[i].subTaskId,
				Result:    bt.subTasks[i].result,
			})
		}
	}
	return plist
}

func (m *brainTaskManager) DeleteTask(task_id string) []interface{} {
	value, found := m.tasks.LoadAndDelete(task_id)
	if !found {
		return nil
	}
	bt := value.(*brainTask)
	bt.lock.Lock()
	defer bt.lock.Unlock()

	rlist := []interface{}{}
	for i := range bt.subTasks {
		rlist = append(rlist, bt.subTasks[i].result)
	}
	return rlist
}

func (m *brainTaskManager) IsTaskDone(task_id string) bool {
	value, found := m.tasks.Load(task_id)
	if !found {
		return true
	}
	bt := value.(*brainTask)
	return bt.isDone
}
