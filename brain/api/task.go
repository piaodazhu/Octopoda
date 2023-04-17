package api

import (
	"brain/rdb"

	"github.com/gin-gonic/gin"
)

func TaskState(ctx *gin.Context) {
	taskid := ctx.Query("taskid")
	rmsg := RMSG{}
	rcode := 200
	if len(taskid) == 0 {
		rmsg.Msg = "Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}
	code, result := rdb.TaskQuery(taskid)
	switch code {
	case rdb.TaskFinished:
		ctx.Data(200, "application/json", []byte(result))
		return
	case rdb.TaskProcessing:
		rmsg.Msg = "TaskProcessing"
		rcode = 202
	case rdb.TaskNotFound:
		rmsg.Msg = "TaskNotFound"
		rcode = 404
	case rdb.TaskDbError:
		rmsg.Msg = "TaskDbError"
		rcode = 500
	}
	ctx.JSON(rcode, rmsg)
}
