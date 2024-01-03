package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/protocols"
)

func TaskState(ctx *gin.Context) {
	taskid := ctx.Query("taskid")
	rmsg := protocols.Result{
		Rmsg: "OK",
	}
	rcode := http.StatusOK
	if len(taskid) == 0 {
		rmsg.Rmsg = "Wrong Args"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	// ONLY FOR HACK. NO CANCEL!!!
	if !model.BrainTaskManager.IsTaskDone(taskid) {
		rmsg.Rmsg = "TaskProcessing"
		rcode = http.StatusAccepted
		ctx.JSON(rcode, rmsg)
		return
	}
	rlist := model.BrainTaskManager.DeleteTask(taskid)
	results := []*protocols.ExecutionResult{}
	for i := range rlist {
		results = append(results, rlist[i].(*protocols.ExecutionResult))
	}
	ctx.JSON(rcode, results)
}
