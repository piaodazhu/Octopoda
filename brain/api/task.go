package api

import (
	"brain/message"
	"brain/rdb"

	"github.com/gin-gonic/gin"
)

func TaskState(ctx *gin.Context) {
	taskid := ctx.Query("taskid")
	rmsg := message.Result{
		Rmsg: "OK",
	}
	rcode := 200
	if len(taskid) == 0 {
		rmsg.Rmsg = "Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}
	code, result := rdb.TaskQuery(taskid)
	switch code {
	case rdb.TaskFinished:
		rdb.TaskDelete(taskid)
		ctx.Data(200, "application/json", []byte(result))
		return
	case rdb.TaskFailed:
		rdb.TaskDelete(taskid)
		ctx.Data(204, "application/json", []byte(result))
		return
	case rdb.TaskProcessing:
		rmsg.Rmsg = "TaskProcessing"
		rcode = 202
	case rdb.TaskNotFound:
		rmsg.Rmsg = "TaskNotFound"
		rcode = 404
	case rdb.TaskDbError:
		rmsg.Rmsg = "TaskDbError"
		rcode = 500
	}
	ctx.JSON(rcode, rmsg)
}
