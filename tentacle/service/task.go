package service

import (
	"net"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
	"github.com/piaodazhu/Octopoda/tentacle/task"
)

// Must return a task struct
func TaskWaitResult(conn net.Conn, serialNum uint32, raw []byte) {
	result, found := task.TaskManager.WaitTask(string(raw))
	if !found {
		logger.Exceptions.Print("TaskWaitState input invalid taskId: ", string(raw))
	}

	if result == nil {
		result = &protocols.Result{
			Rcode: -2,
			Rmsg:  "Canceled",
		}
	}
	serialized, _ := config.Jsoner.Marshal(result)
	err := protocols.SendMessageUnique(conn, protocols.TypeWaitTaskResponse, serialNum, serialized)
	if err != nil {
		logger.Comm.Println("TaskWaitState service error")
	}
}

func TaskQueryState(conn net.Conn, serialNum uint32, raw []byte) {
	task, found := task.TaskManager.QueryTask(string(raw))
	if !found {
		logger.Exceptions.Print("TaskQueryState input invalid taskId: ", string(raw))
	}
	serialized, _ := config.Jsoner.Marshal(&task)
	err := protocols.SendMessageUnique(conn, protocols.TypeQueryTaskResponse, serialNum, serialized)
	if err != nil {
		logger.Comm.Println("TaskQueryState service error")
	}
}

func TaskListAll(conn net.Conn, serialNum uint32, raw []byte) {
	tasks := task.TaskManager.ListTasks()
	serialized, _ := config.Jsoner.Marshal(&tasks)
	err := protocols.SendMessageUnique(conn, protocols.TypeListTasksResponse, serialNum, serialized)
	if err != nil {
		logger.Comm.Println("TaskListAll service error")
	}
}

func TaskCancel(conn net.Conn, serialNum uint32, raw []byte) {
	tasks := task.TaskManager.CancelTask(string(raw))
	serialized, _ := config.Jsoner.Marshal(&tasks)
	err := protocols.SendMessageUnique(conn, protocols.TypeCancelTaskResponse, serialNum, serialized)
	if err != nil {
		logger.Comm.Println("TaskCancel service error")
	}
}
