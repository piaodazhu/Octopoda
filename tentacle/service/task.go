package service

import (
	"net"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
	"tentacle/task"
)

// Must return a task struct
func TaskWaitResult(conn net.Conn, serialNum uint32, raw []byte) {
	result, found := task.TaskManager.WaitTask(string(raw))
	if !found {
		logger.Exceptions.Print("TaskWaitState input invalid taskId: ", string(raw))
		result = &message.Result{
			Rcode: -2,
			Rmsg: "Canceled",
		}
	}
	serialized, _ := config.Jsoner.Marshal(result)
	err := message.SendMessageUnique(conn, message.TypeWaitTaskResponse, serialNum, serialized)
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
	err := message.SendMessageUnique(conn, message.TypeQueryTaskResponse, serialNum, serialized)
	if err != nil {
		logger.Comm.Println("TaskQueryState service error")
	}
}

func TaskListAll(conn net.Conn, serialNum uint32, raw []byte) {
	tasks := task.TaskManager.ListTasks()
	serialized, _ := config.Jsoner.Marshal(&tasks)
	err := message.SendMessageUnique(conn, message.TypeListTasksResponse, serialNum, serialized)
	if err != nil {
		logger.Comm.Println("TaskListAll service error")
	}
}
