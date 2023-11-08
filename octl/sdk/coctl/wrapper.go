package main

/*
struct node_info {
    char* Name;
    char* Version;
    char* Addr;
    int State;
    char* ConnState;
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
	"fmt"

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
func NodeInfo(name string, result *C.struct_node_info) (C.int, *C.char) {
	res, err := sdk.NodeInfo(name)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	*result = C.struct_node_info{
		Name:      C.CString(res.Name),
		Version:   C.CString(res.Version),
		Addr:      C.CString(res.Addr),
		State:     C.int(res.State),
		ConnState: C.CString(res.ConnState),
		OnlineTs:  C.longlong(res.OnlineTs),
		OfflineTs: C.longlong(res.OfflineTs),
		ActiveTs:  C.longlong(res.ActiveTs),
		BrainTs:   C.longlong(res.BrainTs),
	}
	return 0, nil
}

//export NodesInfo
func NodesInfo(names []string, brain *C.struct_brain_info, results []C.struct_node_info, size *C.int) (C.int, *C.char) {
	res, err := sdk.NodesInfo(names)
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
			ConnState: C.CString(res.InfoList[i].ConnState),
			OnlineTs:  C.longlong(res.InfoList[i].OnlineTs),
			OfflineTs: C.longlong(res.InfoList[i].OfflineTs),
			ActiveTs:  C.longlong(res.InfoList[i].ActiveTs),
			BrainTs:   C.longlong(res.InfoList[i].BrainTs),
		}
	}
	return 0, nil
}

//export NodeStatus
func NodeStatus(name string, result *C.struct_node_status) (C.int, *C.char) {
	res, err := sdk.NodeStatus(name)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}

	*result = C.struct_node_status{
		Name:         C.CString(res.Name),
		Platform:     C.CString(res.Platform),
		CpuCores:     C.int(res.CpuCores),
		LocalTime:    C.longlong(res.LocalTime.Unix()),
		CpuLoadShort: C.double(res.CpuLoadShort),
		CpuLoadLong:  C.double(res.CpuLoadLong),
		MemUsed:      C.longlong(res.MemUsed),
		MemTotal:     C.longlong(res.MemTotal),
		DiskUsed:     C.longlong(res.DiskUsed),
		DiskTotal:    C.longlong(res.DiskTotal),
	}
	return 0, nil
}

//export NodesStatus
func NodesStatus(names []string, results []C.struct_node_status, size *C.int) (C.int, *C.char) {
	res, err := sdk.NodesStatus(names)
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

//export DistribFile
func DistribFile(localFileOrDir string, targetPath string, names []string, results []C.struct_execution_result, size *C.int) (C.int, *C.char) {
	res, err := sdk.DistribFile(localFileOrDir, targetPath, names)
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

//export PullFile
func PullFile(pathtype, node, fileOrDir, targetdir string, result *C.struct_execution_result) (C.int, *C.char) {
	res, err := sdk.PullFile(pathtype, node, fileOrDir, targetdir)
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

//export Run
func Run(runtask string, names []string, results []C.struct_execution_result, size *C.int) (C.int, *C.char) {
	res, err := sdk.Run(runtask, names)
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

//export XRun
func XRun(runtask string, names []string, delay int, results []C.struct_execution_result, size *C.int) (C.int, *C.char) {
	res, err := sdk.XRun(runtask, names, delay)
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

//export GroupGetAll
func GroupGetAll(results []*C.char, size *C.int) (C.int, *C.char) {
	res, err := sdk.GroupGetAll()
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	if len(results) < len(res) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res)))
	}

	*size = C.int(len(res))
	for i, name := range res {
		results[i] = C.CString(name)
	}
	// copy(results, res)
	return 0, nil
}

//export GroupGet
func GroupGet(name string, results []*C.char, size *C.int) (C.int, *C.char) {
	res, err := sdk.GroupGet(name)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	if len(results) < len(res.Nodes) {
		return C.int(errs.OctlSdkBufferError), C.CString(fmt.Sprintf("receiver buffer length is only %d but require %d", len(results), len(res.Nodes)))
	}

	*size = C.int(len(res.Nodes))
	for i, name := range res.Nodes {
		results[i] = C.CString(name)
	}
	// copy(results, res.Nodes)
	return 0, nil
}

//export GroupSet
func GroupSet(name string, nocheck bool, names []string) (C.int, *C.char) {
	err := sdk.GroupSet(name, nocheck, names)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	return 0, nil
}

//export GroupDel
func GroupDel(name string) (C.int, *C.char) {
	err := sdk.GroupDel(name)
	if err != nil {
		return C.int(err.Code()), C.CString(err.Error())
	}
	return 0, nil
}

func main() {}
