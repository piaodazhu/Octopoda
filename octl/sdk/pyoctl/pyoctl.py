import ctypes
from typing import List, Tuple, Dict
from enum import Enum
import json

class node_info(ctypes.Structure):
    _fields_ = [("name", ctypes.c_char_p), ("version", ctypes.c_char_p), ("address", ctypes.c_char_p), ("state", ctypes.c_int), ("conn_state", ctypes.c_int), ("delay", ctypes.c_int64), ("online_ts", ctypes.c_int64), ("offline_ts", ctypes.c_int64), ("active_ts", ctypes.c_int64), ("brain_ts", ctypes.c_int64)]

class brain_info(ctypes.Structure):
    _fields_ = [("name", ctypes.c_char_p), ("version", ctypes.c_char_p), ("address", ctypes.c_char_p)]

class node_status(ctypes.Structure):
    _fields_ = [("name", ctypes.c_char_p), ("platform", ctypes.c_char_p), ("cpu_cores", ctypes.c_int), ("local_time", ctypes.c_int64), ("cpuload_short", ctypes.c_double), ("cpuload_long", ctypes.c_double), ("memory_used", ctypes.c_int64), ("memory_total", ctypes.c_int64), ("disk_used", ctypes.c_int64), ("disk_total", ctypes.c_int64)]

class execution_result(ctypes.Structure):
    _fields_ = [("name", ctypes.c_char_p), ("code", ctypes.c_int), ("communication_error_msg", ctypes.c_char_p), ("process_error_msg", ctypes.c_char_p), ("result", ctypes.c_char_p)]

class NodeInfo:
    def __init__(self, cstruct: node_info):
        self.name = bytes.decode(cstruct.name)
        self.version = bytes.decode(cstruct.version)
        self.address = bytes.decode(cstruct.address) 
        self.state = cstruct.state
        self.conn_state = cstruct.conn_state
        self.delay = cstruct.delay
        self.online_ts = cstruct.online_ts
        self.offline_ts = cstruct.offline_ts
        self.active_ts = cstruct.active_ts
        self.brain_ts = cstruct.brain_ts
    def __str__(self) -> str:
        return f"tentacle {self.name} @ {self.address}"

class BrainInfo:
    def __init__(self, cstruct: brain_info):
        self.name = bytes.decode(cstruct.name)
        self.version = bytes.decode(cstruct.version)
        self.address = bytes.decode(cstruct.address)
    def __str__(self) -> str:
        return f"brain {self.name} @ {self.address}"

class NodeStatus:
    def __init__(self, cstruct: node_status):
        self.name = bytes.decode(cstruct.name)
        self.platform = bytes.decode(cstruct.platform)
        self.cpu_cores = cstruct.cpu_cores
        self.local_time = cstruct.local_time
        self.cpuload_short = cstruct.cpuload_short
        self.cpuload_long = cstruct.cpuload_long
        self.memory_used = cstruct.memory_used
        self.memory_total = cstruct.memory_total
        self.disk_used = cstruct.disk_used
        self.disk_total = cstruct.disk_total
    def __str__(self) -> str:
        return f"status {self.name} ({self.platform}) : cpu={self.cpuload_short}, mem={self.memory_used}, disk={self.disk_used}"

class ExecutionResult:
    def __init__(self, cstruct: node_status):
        self.name = bytes.decode(cstruct.name)
        self.code = cstruct.code
        self.communication_error_msg = bytes.decode(cstruct.communication_error_msg) 
        self.process_error_msg = bytes.decode(cstruct.process_error_msg)
        self.result = bytes.decode(cstruct.result)
    def __str__(self) -> str:
        if self.code == 0:
            return f"result {self.name} : [OK]"
        elif self.code == 1:
            return f"result {self.name} : [CommunicationError], msg={self.communication_error_msg}"
        elif self.code == 2:
            return f"result {self.name} : [ProcessError], msg={self.process_error_msg}"
        else:
            return f"result {self.name} : [Unknown], msg={self.result}"

class OctlErrorCode(Enum):
    OctlReadConfigError = 1
    OctlInitClientError = 2
    OctlHttpRequestError = 3
    OctlHttpStatusError = 4
    OctlMessageParseError = 5
    OctlNodeParseError = 6
    OctlFileOperationError = 7
    OctlGitOperationError = 8
    OctlTaskWaitingError = 9
    OctlArgumentError = 10
    OctlSdkNotInitializedError = 11
    OctlSdkPanicRecoverError = 12
    OctlSdkBufferError = 13
    OctlUnknownError = 255

class OctlException(Exception):
    def __init__(self, code:int, emsg:str):
        if code in [member.value for member in OctlErrorCode]:
            self.code = OctlErrorCode(code)
        else:
            self.code = OctlErrorCode.OctlUnknownError
        self.emsg = emsg
        super().__init__(emsg)
    def __str__(self) -> str:
        return f"OctlException <code={self.code}> <msg={self.emsg}>"

class OctlClient:
    def __init__(self, lib_so_or_dll: str, config_yaml :str):
        # 构造函数，初始化对象的属性
        self.lib = ctypes.CDLL(lib_so_or_dll)
        self.ebuf = ctypes.create_string_buffer(256)
        self.ebuflen = ctypes.c_int(256)
        
        self.lib.octl_init.restype = ctypes.c_int
        self.lib.octl_init.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]

        ret = self.lib.octl_init(str.encode(config_yaml), self.ebuf, self.ebuflen)
        if ret != 0:
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))

        self.lib.octl_get_node_info.restype = ctypes.c_int
        self.lib.octl_get_node_info.argtypes = [ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(brain_info), ctypes.POINTER(node_info), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_get_node_status.restype = ctypes.c_int
        self.lib.octl_get_node_status.argtypes = [ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(node_status), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_upload_file.restype = ctypes.c_int
        self.lib.octl_upload_file.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int, ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(execution_result), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_download_file.restype = ctypes.c_int
        self.lib.octl_download_file.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.POINTER(execution_result), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_run_command.restype = ctypes.c_int
        self.lib.octl_run_command.argtypes = [ctypes.c_char_p, ctypes.c_int, ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(execution_result), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_run_script.restype = ctypes.c_int
        self.lib.octl_run_script.argtypes = [ctypes.c_char_p, ctypes.c_int, ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(execution_result), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_run_command_background.restype = ctypes.c_int
        self.lib.octl_run_command_background.argtypes = [ctypes.c_char_p, ctypes.c_int, ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(execution_result), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_xrun_command.restype = ctypes.c_int
        self.lib.octl_xrun_command.argtypes = [ctypes.c_char_p, ctypes.c_int, ctypes.c_int, ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(execution_result), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_xrun_script.restype = ctypes.c_int
        self.lib.octl_xrun_script.argtypes = [ctypes.c_char_p, ctypes.c_int, ctypes.c_int, ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(execution_result), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        # self.lib.octl_get_groups_list.restype = ctypes.c_int
        # self.lib.octl_get_groups_list.argtypes = [ctypes.POINTER(ctypes.c_char_p), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        # self.lib.octl_get_group.restype = ctypes.c_int
        # self.lib.octl_get_group.argtypes = [ctypes.c_char_p, ctypes.POINTER(ctypes.c_char_p), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        # self.lib.octl_set_group.restype = ctypes.c_int
        # self.lib.octl_set_group.argtypes = [ctypes.c_char_p, ctypes.c_int, ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        # self.lib.octl_del_group.restype = ctypes.c_int
        # self.lib.octl_del_group.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        # self.lib.octl_prune_nodes.restype = ctypes.c_int
        # self.lib.octl_prune_nodes.argtypes = [ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_get_scenario_info.restype = ctypes.c_int
        self.lib.octl_get_scenario_info.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_get_scenarios_info_list.restype = ctypes.c_int
        self.lib.octl_get_scenarios_info_list.argtypes = [ctypes.POINTER(ctypes.c_char_p), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_get_scenario_version.restype = ctypes.c_int
        self.lib.octl_get_scenario_version.argtypes = [ctypes.c_char_p, ctypes.c_int,  ctypes.c_int, ctypes.c_char_p, ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_get_nodeapp_info.restype = ctypes.c_int
        self.lib.octl_get_nodeapp_info.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_get_nodeapps_info_list.restype = ctypes.c_int
        self.lib.octl_get_nodeapps_info_list.argtypes = [ctypes.c_char_p, ctypes.POINTER(ctypes.c_char_p), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]
        self.lib.octl_apply_scenario.restype = ctypes.c_int
        self.lib.octl_apply_scenario.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int, ctypes.POINTER(ctypes.c_char_p), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.POINTER(ctypes.c_int)]

        self.lib.octl_clear_node_info.restype = None
        self.lib.octl_clear_node_info.argtypes = [ctypes.POINTER(node_info)]
        self.lib.octl_clear_node_status.restype = None
        self.lib.octl_clear_node_status.argtypes = [ctypes.POINTER(node_status)]
        self.lib.octl_clear_brain_info.restype = None
        self.lib.octl_clear_brain_info.argtypes = [ctypes.POINTER(brain_info)]
        self.lib.octl_clear_execution_result.restype = None
        self.lib.octl_clear_execution_result.argtypes = [ctypes.POINTER(execution_result)]
        self.lib.octl_clear_node_info_list.restype = None
        self.lib.octl_clear_node_info_list.argtypes = [ctypes.POINTER(node_info), ctypes.c_int]
        self.lib.octl_clear_node_status_list.restype = None
        self.lib.octl_clear_node_status_list.argtypes = [ctypes.POINTER(node_status), ctypes.c_int]
        self.lib.octl_clear_execution_result_list.restype = None
        self.lib.octl_clear_execution_result_list.argtypes = [ctypes.POINTER(execution_result), ctypes.c_int]
        self.lib.octl_clear_string_list.restype = None
        self.lib.octl_clear_string_list.argtypes = [ctypes.POINTER(ctypes.c_char_p), ctypes.c_int]

    def get_node_info(self, names: List[str]) -> Tuple[BrainInfo, List[NodeInfo]]:
        binfo = brain_info()
        output_len = 1024
        if len(names) != 0:
            output_len = len(names)
        
        names_input = (ctypes.c_char_p * len(names))()
        ninfos_output = (node_info * output_len)()
        ninfolen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])

        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_get_node_info(names_input, len(names), binfo, ninfos_output, ninfolen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_brain_info(binfo)
            self.lib.octl_clear_node_info_list(ninfos_output, ninfolen_output)
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))

        ninfo_list = []
        for idx in range(ninfolen_output.value):
            ninfo_list.append(NodeInfo(ninfos_output[idx]))

        binfo_obj = BrainInfo(binfo)

        self.lib.octl_clear_brain_info(binfo)
        self.lib.octl_clear_node_info_list(ninfos_output, ninfolen_output)
        return binfo_obj, ninfo_list
    
    def get_node_status(self, names: List[str]) -> List[NodeStatus]:
        output_len = 1024
        if len(names) != 0:
            output_len = len(names)
        
        names_input = (ctypes.c_char_p * len(names))()
        nstatus_output = (node_status * output_len)()
        nstatuslen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])

        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_get_node_status(names_input, len(names), nstatus_output, nstatuslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_node_status_list(nstatus_output, nstatuslen_output)
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))

        nstatus_list = []
        for idx in range(nstatuslen_output.value):
            nstatus_list.append(NodeStatus(nstatus_output[idx]))

        self.lib.octl_clear_node_status_list(nstatus_output, nstatuslen_output)
        return nstatus_list

    def upload_file(self, local_file_or_dir: str, remote_target_path: str, names: List[str], is_force: bool = True) -> List[ExecutionResult]:
        if len(names) == 0:
            return
        output_len = len(names)
        names_input = (ctypes.c_char_p * len(names))()
        results_output = (execution_result * output_len)()
        resultslen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])
        
        force_flag = ctypes.c_int(0)
        if is_force:
            force_flag = ctypes.c_int(1)

        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_upload_file(str.encode(local_file_or_dir), str.encode(remote_target_path), force_flag, names_input, len(names), results_output, resultslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_execution_result_list(results_output, resultslen_output)
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))

        result_list = []
        for idx in range(resultslen_output.value):
                result_list.append(ExecutionResult(results_output[idx]))
        
        self.lib.octl_clear_execution_result_list(results_output, resultslen_output)
        return result_list

    def download_file(self, remote_file_or_dir: str, local_dir: str, name: str) -> ExecutionResult:
        result_output = execution_result()
        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_download_file(str.encode(remote_file_or_dir), str.encode(local_dir), str.encode(name), result_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_execution_result(result_output)
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))
        
        ret_obj = ExecutionResult(result_output)
        self.lib.octl_clear_execution_result(result_output)
        return ret_obj

    def run_command(self, cmd_str: str, names: List[str], should_align: bool = False) -> List[ExecutionResult]:
        if len(names) == 0:
            return
        align_flag = ctypes.c_int(0)
        if should_align:
            align_flag = ctypes.c_int(1)
        output_len = len(names)
        names_input = (ctypes.c_char_p * len(names))()
        results_output = (execution_result * output_len)()
        resultslen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])
        
        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_run_command(str.encode(cmd_str), align_flag, names_input, len(names), results_output, resultslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_execution_result_list(results_output, resultslen_output)
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))

        result_list = []
        for idx in range(resultslen_output.value):
                result_list.append(ExecutionResult(results_output[idx]))
        
        self.lib.octl_clear_execution_result_list(results_output, resultslen_output)
        return result_list

    def run_script(self, script_file: str, names: List[str], should_align: bool = False) -> List[ExecutionResult]:
        if len(names) == 0:
            return
        align_flag = ctypes.c_int(0)
        if should_align:
            align_flag = ctypes.c_int(1)
        output_len = len(names)
        names_input = (ctypes.c_char_p * len(names))()
        results_output = (execution_result * output_len)()
        resultslen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])
        
        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_run_script(str.encode(script_file), align_flag, names_input, len(names), results_output, resultslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_execution_result_list(results_output, resultslen_output)
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))

        result_list = []
        for idx in range(resultslen_output.value):
                result_list.append(ExecutionResult(results_output[idx]))
        
        self.lib.octl_clear_execution_result_list(results_output, resultslen_output)
        return result_list

    def run_command_background(self, cmd_str: str, names: List[str], should_align: bool = False) -> List[ExecutionResult]:
        if len(names) == 0:
            return
        align_flag = ctypes.c_int(0)
        if should_align:
            align_flag = ctypes.c_int(1)
        output_len = len(names)
        names_input = (ctypes.c_char_p * len(names))()
        results_output = (execution_result * output_len)()
        resultslen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])
        
        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_run_command_background(str.encode(cmd_str), align_flag, names_input, len(names), results_output, resultslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_execution_result_list(results_output, resultslen_output)
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))

        result_list = []
        for idx in range(resultslen_output.value):
                result_list.append(ExecutionResult(results_output[idx]))
        
        self.lib.octl_clear_execution_result_list(results_output, resultslen_output)
        return result_list

    def xrun_command(self, cmd_str: str, names: List[str], delay: int = 0, should_align: bool = False) -> List[ExecutionResult]:
        if len(names) == 0:
            return
        align_flag = ctypes.c_int(0)
        if should_align:
            align_flag = ctypes.c_int(1)
        output_len = len(names)
        names_input = (ctypes.c_char_p * len(names))()
        results_output = (execution_result * output_len)()
        resultslen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])
        
        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_xrun_command(str.encode(cmd_str), align_flag, ctypes.c_int(delay), names_input, len(names), results_output, resultslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_execution_result_list(results_output, resultslen_output)
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))

        result_list = []
        for idx in range(resultslen_output.value):
                result_list.append(ExecutionResult(results_output[idx]))
        
        self.lib.octl_clear_execution_result_list(results_output, resultslen_output)
        return result_list

    def xrun_script(self, script_file: str, names: List[str], delay: int = 0, should_align: bool = False) -> List[ExecutionResult]:
        if len(names) == 0:
            return
        align_flag = ctypes.c_int(0)
        if should_align:
            align_flag = ctypes.c_int(1)
        output_len = len(names)
        names_input = (ctypes.c_char_p * len(names))()
        results_output = (execution_result * output_len)()
        resultslen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])
        
        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_xrun_script(str.encode(script_file), align_flag, ctypes.c_int(delay), names_input, len(names), results_output, resultslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_execution_result_list(results_output, resultslen_output)
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))

        result_list = []
        for idx in range(resultslen_output.value):
                result_list.append(ExecutionResult(results_output[idx]))
        
        self.lib.octl_clear_execution_result_list(results_output, resultslen_output)
        return result_list

    # def get_groups_list(self) -> List[str]:
    #     output_len = 1024
    #     names_output = (ctypes.c_char_p * output_len)()
    #     nameslen_output = ctypes.c_int(output_len)
    #     self.ebuflen = ctypes.c_int(256)
    #     ret = self.lib.octl_get_groups_list(names_output, nameslen_output, self.ebuf, self.ebuflen)
    #     if ret != 0:
    #         self.lib.octl_clear_string_list(names_output, nameslen_output)
    #         raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))
        
    #     names_list = []
    #     for idx in range(nameslen_output.value):
    #         names_list.append(bytes.decode(names_output[idx]))

    #     self.lib.octl_clear_string_list(names_output, nameslen_output)
    #     return names_list
    
    # def get_group(self, name: str) -> List[str]:
    #     output_len = 1024
    #     names_output = (ctypes.c_char_p * output_len)()
    #     nameslen_output = ctypes.c_int(output_len)
    #     self.ebuflen = ctypes.c_int(256)
    #     ret = self.lib.octl_get_group(str.encode(name), names_output, nameslen_output, self.ebuf, self.ebuflen)
    #     if ret != 0:
    #         self.lib.octl_clear_string_list(names_output, nameslen_output)
    #         raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))
        
    #     names_list = []
    #     for idx in range(nameslen_output.value):
    #         names_list.append(bytes.decode(names_output[idx]))

    #     self.lib.octl_clear_string_list(names_output, nameslen_output)
    #     return names_list

    # def set_group(self, name: str, skipCheck: bool, members: List[str]) -> None:
    #     if len(members) == 0:
    #         return
        
    #     members_input = (ctypes.c_char_p * len(members))()
    #     for idx in range(len(members)):
    #         members_input[idx] = str.encode(members[idx])
        
    #     skipcheck_input = ctypes.c_int(0)
    #     if skipCheck:
    #         skipcheck_input = ctypes.c_int(1)
    #     self.ebuflen = ctypes.c_int(256)
    #     ret = self.lib.octl_set_group(str.encode(name), skipcheck_input, members_input, len(members), self.ebuf, self.ebuflen)
    #     if ret != 0:
    #         raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))

    # def del_group(self, name: str) -> None:
    #     self.ebuflen = ctypes.c_int(256)
    #     ret = self.lib.octl_del_group(str.encode(name), self.ebuf, self.ebuflen)
    #     if ret != 0:
    #         raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))

    # def prune_nodes(self) -> None:
    #     self.ebuflen = ctypes.c_int(256)
    #     ret = self.lib.octl_prune_nodes(self.ebuf, self.ebuflen)
    #     if ret != 0:
    #         raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))
        
    def get_scenario_info(self, name: str) -> Dict:
        rawbuf = ctypes.create_string_buffer(4096)
        rawlen = ctypes.c_int(4096)
        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_get_scenario_info(str.encode(name), rawbuf, rawlen, self.ebuf, self.ebuflen)
        if ret != 0:
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))
        return json.loads(bytes.decode(rawbuf[:rawlen.value]))
        
    def get_scenarios_info_list(self) -> List[Dict]:
        output_len = 1024
        scens_output = (ctypes.c_char_p * output_len)()
        scenslen_output = ctypes.c_int(output_len)
        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_get_scenarios_info_list(scens_output, scenslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_string_list(scens_output, scenslen_output)
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))
        
        scens_list = []
        for idx in range(scenslen_output.value):
            scens_list.append(json.loads(bytes.decode(scens_output[idx])))

        self.lib.octl_clear_string_list(scens_output, scenslen_output)
        return scens_list
        
    def get_scenario_version(self, name: str, offset: int = 0, limit: int = 3) -> Dict:
        rawbuf = ctypes.create_string_buffer(4096)
        rawlen = ctypes.c_int(4096)
        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_get_scenario_version(str.encode(name), ctypes.c_int(offset), ctypes.c_int(limit), rawbuf, rawlen, self.ebuf, self.ebuflen)
        if ret != 0:
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))
        return json.loads(bytes.decode(rawbuf[:rawlen.value]))
        
    def get_nodeapps_info_list(self, name: str) -> List[Dict]:
        output_len = 1024
        nodeapps_output = (ctypes.c_char_p * output_len)()
        nodeappslen_output = ctypes.c_int(output_len)
        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_get_nodeapps_info_list(str.encode(name), nodeapps_output, nodeappslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_string_list(nodeapps_output, nodeappslen_output)
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))
        
        nodeapps_list = []
        for idx in range(nodeappslen_output.value):
            nodeapps_list.append(json.loads(bytes.decode(nodeapps_output[idx])))

        self.lib.octl_clear_string_list(nodeapps_output, nodeappslen_output)
        return nodeapps_list
        
    def get_nodeapp_info(self, name: str, app: str, scenario: str) -> Dict:
        rawbuf = ctypes.create_string_buffer(4096)
        rawlen = ctypes.c_int(4096)
        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_get_nodeapp_info(str.encode(name), str.encode(app), str.encode(scenario), rawbuf, rawlen, self.ebuf, self.ebuflen)
        if ret != 0:
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))
        return json.loads(bytes.decode(rawbuf[:rawlen.value]))
    
    def apply_scenario(self, name: str, target: str, message: str, timeout_s: int = 0) -> List[str]:
        output_len = 1024
        logs_output = (ctypes.c_char_p * output_len)()
        logslen_output = ctypes.c_int(output_len)
        self.ebuflen = ctypes.c_int(256)
        ret = self.lib.octl_apply_scenario(str.encode(name), str.encode(target), str.encode(message), ctypes.c_int(timeout_s), logs_output, logslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_string_list(logs_output, logslen_output)
            raise OctlException(ret, bytes.decode(self.ebuf[:self.ebuflen.value]))
        
        logs_list = []
        for idx in range(logslen_output.value):
            logs_list.append(bytes.decode(logs_output[idx]))

        self.lib.octl_clear_string_list(logs_output, logslen_output)
        return logs_list
        

    def __str__(self):
        # 定义对象的字符串表示，可用于打印对象
        return f"octl python client"
