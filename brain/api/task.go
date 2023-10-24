package api

import (
	"brain/model"
	"protocols"

	"github.com/gin-gonic/gin"
)

func TaskState(ctx *gin.Context) {
	taskid := ctx.Query("taskid")
	rmsg := protocols.Result{
		Rmsg: "OK",
	}
	rcode := 200
	if len(taskid) == 0 {
		rmsg.Rmsg = "Wrong Args"
		ctx.JSON(400, rmsg)
		return
	}

	// ONLY FOR HACK. NO CANCEL!!!
	if !model.BrainTaskManager.IsTaskDone(taskid) {
		rmsg.Rmsg = "TaskProcessing"
		rcode = 202
		ctx.JSON(rcode, rmsg)
		return
	}
	rlist := model.BrainTaskManager.DeleteTask(taskid)
	results := []*BasicNodeResults{}
	for i := range rlist {
		results = append(results, rlist[i].(*BasicNodeResults))
	}
	ctx.JSON(rcode, results)
}
