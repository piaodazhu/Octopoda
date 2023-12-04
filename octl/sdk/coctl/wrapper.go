package main

/*
struct node_info {
    char* Name;
    char* Version;
    char* Addr;
    int State;
    int ConnState;
    long long Delay;
    long long OnlineTs;
    long long OfflineTs;
    long long ActiveTs;
    long long BrainTs;
};

struct brain_info {
    char* Name;
    char* Version;
    char* Addr;
};

struct node_status {
    char* Name;
    char* Platform;
    int CpuCores;
    long long LocalTime;
    double CpuLoadShort;
    double CpuLoadLong;
    long long MemUsed;
    long long MemTotal;
    long long DiskUsed;
    long long DiskTotal;
};

struct execution_result {
    char* Name;
    int Code;
    char* CommunicationErrorMsg;
    char* ProcessErrorMsg;
    char* Result;
};
*/
import "C"
import (
	"context"
	"fmt"
	"time"

	"github.com/piaodazhu/Octopoda/octl/sdk"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

//export Init
func Init(configFile string) (C.int, *C.char) {
	err := sdk.Init(configFile)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	return 0, nil
}

//export NodeInfo
func NodeInfo(names []string, brain *C.struct_brain_info, results []C.struct_node_info, size *C.int) (C.int, *C.char) {
	res, err := sdk.NodeInfo(names)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	if len(results) < len(res.InfoList) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res.InfoList)))
	}

	*size = C.int(len(res.InfoList))
	*brain = C.struct_brain_info{
		Name:    C.CString(res.BrainName),
		Version: C.CString(res.BrainVersion),
		Addr:    C.CString(res.BrainAddr),
	}
	for i := range res.InfoList {
		results[i] = C.struct_node_info{
			Name:      C.CString(res.InfoList[i].Name),
			Version:   C.CString(res.InfoList[i].Version),
			Addr:      C.CString(res.InfoList[i].Addr),
			State:     C.int(res.InfoList[i].State),
			ConnState: C.int(res.InfoList[i].ConnState),
			Delay:     C.longlong(res.InfoList[i].Delay),
			OnlineTs:  C.longlong(res.InfoList[i].OnlineTs),
			OfflineTs: C.longlong(res.InfoList[i].OfflineTs),
			ActiveTs:  C.longlong(res.InfoList[i].ActiveTs),
			BrainTs:   C.longlong(res.InfoList[i].BrainTs),
		}
	}
	return 0, nil
}

//export NodeStatus
func NodeStatus(names []string, results []C.struct_node_status, size *C.int) (C.int, *C.char) {
	res, err := sdk.NodeStatus(names)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	if len(results) < len(res.StatusList) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res.StatusList)))
	}

	*size = C.int(len(res.StatusList))
	for i := range res.StatusList {
		results[i] = C.struct_node_status{
			Name:         C.CString(res.StatusList[i].Name),
			Platform:     C.CString(res.StatusList[i].Platform),
			CpuCores:     C.int(res.StatusList[i].CpuCores),
			LocalTime:    C.longlong(res.StatusList[i].LocalTime.Unix()),
			CpuLoadShort: C.double(res.StatusList[i].CpuLoadShort),
			CpuLoadLong:  C.double(res.StatusList[i].CpuLoadLong),
			MemUsed:      C.longlong(res.StatusList[i].MemUsed),
			MemTotal:     C.longlong(res.StatusList[i].MemTotal),
			DiskUsed:     C.longlong(res.StatusList[i].DiskUsed),
			DiskTotal:    C.longlong(res.StatusList[i].DiskTotal),
		}
	}
	return 0, nil
}

//export UploadFile
func UploadFile(localFileOrDir string, remoteTargetPath string, isForce bool, names []string, results []C.struct_execution_result, size *C.int) (C.int, *C.char) {
	res, err := sdk.UploadFile(localFileOrDir, remoteTargetPath, isForce, names)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	if len(results) < len(res) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res)))
	}

	*size = C.int(len(res))
	for i := range res {
		results[i] = C.struct_execution_result{
			Name:                  C.CString(res[i].Name),
			Code:                  C.int(res[i].Code),
			CommunicationErrorMsg: C.CString(res[i].CommunicationErrorMsg),
			ProcessErrorMsg:       C.CString(res[i].ProcessErrorMsg),
			Result:                C.CString(res[i].Result),
		}
	}
	return 0, nil
}

//export DownloadFile
func DownloadFile(remoteFileOrDir string, localTargetPath string, name string, result *C.struct_execution_result) (C.int, *C.char) {
	res, err := sdk.DownloadFile(remoteFileOrDir, localTargetPath, name)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	*result = C.struct_execution_result{
		Name:                  C.CString(res.Name),
		Code:                  C.int(res.Code),
		CommunicationErrorMsg: C.CString(res.CommunicationErrorMsg),
		ProcessErrorMsg:       C.CString(res.ProcessErrorMsg),
		Result:                C.CString(res.Result),
	}
	return 0, nil
}

//export RunCommand
func RunCommand(cmd string, needAlign bool, names []string, results []C.struct_execution_result, size *C.int) (C.int, *C.char) {
	res, err := sdk.RunCommand(cmd, needAlign, names)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	if len(results) < len(res) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res)))
	}

	*size = C.int(len(res))
	for i := range res {
		results[i] = C.struct_execution_result{
			Name:                  C.CString(res[i].Name),
			Code:                  C.int(res[i].Code),
			CommunicationErrorMsg: C.CString(res[i].CommunicationErrorMsg),
			ProcessErrorMsg:       C.CString(res[i].ProcessErrorMsg),
			Result:                C.CString(res[i].Result),
		}
	}
	return 0, nil
}

//export RunScript
func RunScript(cmd string, needAlign bool, names []string, results []C.struct_execution_result, size *C.int) (C.int, *C.char) {
	res, err := sdk.RunScript(cmd, needAlign, names)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	if len(results) < len(res) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res)))
	}

	*size = C.int(len(res))
	for i := range res {
		results[i] = C.struct_execution_result{
			Name:                  C.CString(res[i].Name),
			Code:                  C.int(res[i].Code),
			CommunicationErrorMsg: C.CString(res[i].CommunicationErrorMsg),
			ProcessErrorMsg:       C.CString(res[i].ProcessErrorMsg),
			Result:                C.CString(res[i].Result),
		}
	}
	return 0, nil
}

//export RunCommandBackground
func RunCommandBackground(cmd string, needAlign bool, names []string, results []C.struct_execution_result, size *C.int) (C.int, *C.char) {
	res, err := sdk.RunCommandBackground(cmd, needAlign, names)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	if len(results) < len(res) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res)))
	}

	*size = C.int(len(res))
	for i := range res {
		results[i] = C.struct_execution_result{
			Name:                  C.CString(res[i].Name),
			Code:                  C.int(res[i].Code),
			CommunicationErrorMsg: C.CString(res[i].CommunicationErrorMsg),
			ProcessErrorMsg:       C.CString(res[i].ProcessErrorMsg),
			Result:                C.CString(res[i].Result),
		}
	}
	return 0, nil
}

//export XRunCommand
func XRunCommand(cmd string, needAlign bool, delay int, names []string, results []C.struct_execution_result, size *C.int) (C.int, *C.char) {
	res, err := sdk.XRunCommand(cmd, needAlign, delay, names)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	if len(results) < len(res) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res)))
	}

	*size = C.int(len(res))
	for i := range res {
		results[i] = C.struct_execution_result{
			Name:                  C.CString(res[i].Name),
			Code:                  C.int(res[i].Code),
			CommunicationErrorMsg: C.CString(res[i].CommunicationErrorMsg),
			ProcessErrorMsg:       C.CString(res[i].ProcessErrorMsg),
			Result:                C.CString(res[i].Result),
		}
	}
	return 0, nil
}

//export XRunScript
func XRunScript(cmd string, needAlign bool, delay int, names []string, results []C.struct_execution_result, size *C.int) (C.int, *C.char) {
	res, err := sdk.XRunScript(cmd, needAlign, delay, names)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	if len(results) < len(res) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res)))
	}

	*size = C.int(len(res))
	for i := range res {
		results[i] = C.struct_execution_result{
			Name:                  C.CString(res[i].Name),
			Code:                  C.int(res[i].Code),
			CommunicationErrorMsg: C.CString(res[i].CommunicationErrorMsg),
			ProcessErrorMsg:       C.CString(res[i].ProcessErrorMsg),
			Result:                C.CString(res[i].Result),
		}
	}
	return 0, nil
}

// //export GroupGetAll
// func GroupGetAll(results []*C.char, size *C.int) (C.int, *C.char) {
// 	res, err := sdk.GroupGetAll()
// 	if err != nil {
// 		return C.int(err.Code()), C.CString(err.Error())
// 	}
// 	if len(results) < len(res) {
// 		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res)))
// 	}

// 	*size = C.int(len(res))
// 	for i, name := range res {
// 		results[i] = C.CString(name)
// 	}
// 	// copy(results, res)
// 	return 0, nil
// }

// //export GroupGet
// func GroupGet(name string, results []*C.char, size *C.int) (C.int, *C.char) {
// 	res, err := sdk.GroupGet(name)
// 	if err != nil {
// 		return C.int(err.Code()), C.CString(err.Error())
// 	}
// 	if len(results) < len(res.Nodes) {
// 		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res.Nodes)))
// 	}

// 	*size = C.int(len(res.Nodes))
// 	for i, name := range res.Nodes {
// 		results[i] = C.CString(name)
// 	}
// 	// copy(results, res.Nodes)
// 	return 0, nil
// }

// //export GroupSet
// func GroupSet(name string, nocheck bool, names []string) (C.int, *C.char) {
// 	err := sdk.GroupSet(name, nocheck, names)
// 	if err != nil {
// 		return C.int(err.Code()), C.CString(err.Error())
// 	}
// 	return 0, nil
// }

// //export GroupDel
// func GroupDel(name string) (C.int, *C.char) {
// 	err := sdk.GroupDel(name)
// 	if err != nil {
// 		return C.int(err.Code()), C.CString(err.Error())
// 	}
// 	return 0, nil
// }

// //export Prune
// func Prune() (C.int, *C.char) {
// 	err := sdk.Prune()
// 	if err != nil {
// 		return C.int(err.Code()), C.CString(err.Error())
// 	}
// 	return 0, nil
// }

//export ScenarioInfo
func ScenarioInfo(name string) (*C.char, C.int, *C.char) {
	res, err := sdk.ScenarioInfo(name)
	if err != nil {
		return nil, C.int(err.Code()), C.CString(err.Error())
	}
	return C.CString(string(res)), 0, nil
}

//export ScenariosInfo
func ScenariosInfo(results []*C.char, size *C.int) (C.int, *C.char) {
	res, err := sdk.ScenariosInfo()
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	if len(results) < len(res) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res)))
	}

	*size = C.int(len(res))
	for i := range res {
		results[i] = C.CString(string(res[i]))
	}
	return 0, nil
}

//export ScenarioVersion
func ScenarioVersion(name string, offset, limit int) (*C.char, C.int, *C.char) {
	res, err := sdk.ScenarioVersion(name, offset, limit)
	if err != nil {
		return nil, C.int(err.Code()), C.CString(err.Error())
	}
	return C.CString(string(res)), 0, nil
}

//export NodeAppInfo
func NodeAppInfo(name, app, scenario string) (*C.char, C.int, *C.char) {
	res, err := sdk.NodeAppInfo(name, app, scenario)
	if err != nil {
		return nil, C.int(err.Code()), C.CString(err.Error())
	}
	return C.CString(string(res)), 0, nil
}

//export NodeAppsInfo
func NodeAppsInfo(name string, results []*C.char, size *C.int) (C.int, *C.char) {
	res, err := sdk.NodeAppsInfo(name)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	if len(results) < len(res) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res)))
	}

	*size = C.int(len(res))
	for i := range res {
		results[i] = C.CString(string(res[i]))
	}
	return 0, nil
}

//export Apply
func Apply(deployment, target, message string, timeout int, logs []*C.char, size *C.int) (C.int, *C.char) {
	var ctx context.Context
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
		defer cancel()
	} else {
		ctx = context.Background()
	}
	res, err := sdk.Apply(ctx, deployment, target, message)
	if len(logs) < len(res) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(logs), len(res)))
	}
	*size = C.int(len(res))
	for i := range res {
		logs[i] = C.CString(string(res[i]))
	}

	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	return 0, nil
}

func main() {}
