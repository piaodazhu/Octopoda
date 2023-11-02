package api

import (
	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/protocols"
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
	results := []*protocols.ExecutionResults{}
	for i := range rlist {
		results = append(results, rlist[i].(*protocols.ExecutionResults))
	}
	ctx.JSON(rcode, results)
}
