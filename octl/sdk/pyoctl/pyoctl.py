import ctypes
from typing import List, Tuple

class node_info(ctypes.Structure):
    _fields_ = [("name", ctypes.c_char_p), ("version", ctypes.c_char_p), ("address", ctypes.c_char_p), ("state", ctypes.c_int), ("conn_state", ctypes.c_char_p), ("online_ts", ctypes.c_int64), ("offline_ts", ctypes.c_int64), ("active_ts", ctypes.c_int64), ("brain_ts", ctypes.c_int64)]

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
        self.state = cstruct
        self.conn_state = bytes.decode(cstruct.conn_state)
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

class OctlClient:
    def __init__(self, lib_so_or_dll: str, config_yaml :str):
        # 构造函数，初始化对象的属性
        self.lib = ctypes.CDLL(lib_so_or_dll)
        self.ebuf = ctypes.create_string_buffer(256)
        self.ebuflen = ctypes.c_int(256)
        
        self.lib.octl_init.restype = ctypes.c_int
        self.lib.octl_init.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]

        ret = self.lib.octl_init(str.encode(config_yaml), self.ebuf, self.ebuflen)
        if ret != 0:
            raise Exception(bytes.decode(self.ebuf[:ret]))

        self.lib.octl_get_node_info.restype = ctypes.c_int
        self.lib.octl_get_node_info.argtypes = [ctypes.c_char_p, ctypes.POINTER(node_info), ctypes.c_char_p, ctypes.c_int]
        self.lib.octl_get_nodes_info_list.restype = ctypes.c_int
        self.lib.octl_get_nodes_info_list.argtypes = [ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(brain_info), ctypes.POINTER(node_info), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.c_int]
        self.lib.octl_get_node_status.restype = ctypes.c_int
        self.lib.octl_get_node_status.argtypes = [ctypes.c_char_p, ctypes.POINTER(node_status), ctypes.c_char_p, ctypes.c_int]
        self.lib.octl_get_nodes_status_list.restype = ctypes.c_int
        self.lib.octl_get_nodes_status_list.argtypes = [ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(node_status), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.c_int]
        self.lib.octl_distribute_file.restype = ctypes.c_int
        self.lib.octl_distribute_file.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(execution_result), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.c_int]
        self.lib.octl_pull_file.restype = ctypes.c_int
        self.lib.octl_pull_file.argtypes = [ctypes.c_int8, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.POINTER(execution_result), ctypes.c_char_p, ctypes.c_int]
        self.lib.octl_run.restype = ctypes.c_int
        self.lib.octl_run.argtypes = [ctypes.c_char_p, ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.POINTER(execution_result), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.c_int]
        self.lib.octl_xrun.restype = ctypes.c_int
        self.lib.octl_xrun.argtypes = [ctypes.c_char_p, ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.c_int, ctypes.POINTER(execution_result), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.c_int]
        self.lib.octl_get_groups_list.restype = ctypes.c_int
        self.lib.octl_get_groups_list.argtypes = [ctypes.POINTER(ctypes.c_char_p), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.c_int]
        self.lib.octl_get_group.restype = ctypes.c_int
        self.lib.octl_get_group.argtypes = [ctypes.c_char_p, ctypes.POINTER(ctypes.c_char_p), ctypes.POINTER(ctypes.c_int), ctypes.c_char_p, ctypes.c_int]
        self.lib.octl_set_group.restype = ctypes.c_int
        self.lib.octl_set_group.argtypes = [ctypes.c_char_p, ctypes.c_int, ctypes.POINTER(ctypes.c_char_p), ctypes.c_int, ctypes.c_char_p, ctypes.c_int]
        self.lib.octl_del_group.restype = ctypes.c_int
        self.lib.octl_del_group.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
        
        self.lib.octl_clear_node_info.restype = None
        self.lib.octl_clear_node_info.argtypes = [ctypes.POINTER(node_info)]
        self.lib.octl_clear_node_status.restype = None
        self.lib.octl_clear_node_status.argtypes = [ctypes.POINTER(node_status)]
        self.lib.octl_clear_brain_info.restype = None
        self.lib.octl_clear_brain_info.argtypes = [ctypes.POINTER(brain_info)]
        self.lib.octl_clear_execution_result.restype = None
        self.lib.octl_clear_execution_result.argtypes = [ctypes.POINTER(execution_result)]
        self.lib.octl_clear_nodes_info_list.restype = None
        self.lib.octl_clear_nodes_info_list.argtypes = [ctypes.POINTER(node_info), ctypes.c_int]
        self.lib.octl_clear_nodes_status_list.restype = None
        self.lib.octl_clear_nodes_status_list.argtypes = [ctypes.POINTER(node_status), ctypes.c_int]
        self.lib.octl_clear_execution_results_list.restype = None
        self.lib.octl_clear_execution_results_list.argtypes = [ctypes.POINTER(execution_result), ctypes.c_int]
        self.lib.octl_clear_name_list.restype = None
        self.lib.octl_clear_name_list.argtypes = [ctypes.POINTER(ctypes.c_char_p), ctypes.c_int]


    def get_node_info(self, name: str) -> NodeInfo:
        ninfo = node_info()
        ret = self.lib.octl_get_node_info(str.encode(name), ninfo, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_node_info(ninfo)
            raise Exception(bytes.decode(self.ebuf[:ret]))
        
        retobj = NodeInfo(ninfo)
        self.lib.octl_clear_node_info(ninfo)
        return retobj

    def get_nodes_info_list(self, names: List[str]) -> Tuple[BrainInfo, List[NodeInfo]]:
        binfo = brain_info()
        output_len = 1024
        if len(names) != 0:
            output_len = len(names)
        
        names_input = (ctypes.c_char_p * len(names))()
        ninfos_output = (node_info * output_len)()
        ninfolen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])

        ret = self.lib.octl_get_nodes_info_list(names_input, len(names), binfo, ninfos_output, ninfolen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_brain_info(binfo)
            self.lib.octl_clear_nodes_info_list(ninfos_output, ninfolen_output)
            raise Exception(bytes.decode(self.ebuf[:ret]))

        ninfo_list = []
        for idx in range(ninfolen_output.value):
            ninfo_list.append(NodeInfo(ninfos_output[idx]))

        binfo_obj = BrainInfo(binfo)

        self.lib.octl_clear_brain_info(binfo)
        self.lib.octl_clear_nodes_info_list(ninfos_output, ninfolen_output)
        return binfo_obj, ninfo_list
    
    def get_node_status(self, name: str) -> NodeStatus:
        nstatus = node_status()
        ret = self.lib.octl_get_node_status(str.encode(name), nstatus, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_node_status(nstatus)
            raise Exception(bytes.decode(self.ebuf[:ret]))
        retobj = NodeStatus(nstatus)

        self.lib.octl_clear_node_status(nstatus)
        return retobj

    def get_nodes_status_list(self, names: List[str]) -> List[NodeStatus]:
        output_len = 1024
        if len(names) != 0:
            output_len = len(names)
        
        names_input = (ctypes.c_char_p * len(names))()
        nstatus_output = (node_status * output_len)()
        nstatuslen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])

        ret = self.lib.octl_get_nodes_status_list(names_input, len(names), nstatus_output, nstatuslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_nodes_status_list(nstatus_output, nstatuslen_output)
            raise Exception(bytes.decode(self.ebuf[:ret]))

        nstatus_list = []
        for idx in range(nstatuslen_output.value):
            nstatus_list.append(NodeStatus(nstatus_output[idx]))

        self.lib.octl_clear_nodes_status_list(nstatus_output, nstatuslen_output)
        return nstatus_list

    def distribute_file(self, local_file_or_dir: str, target_path: str, names: List[str]) -> List[ExecutionResult]:
        if len(names) == 0:
            return
        output_len = len(names)
        names_input = (ctypes.c_char_p * len(names))()
        results_output = (execution_result * output_len)()
        resultslen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])
        
        ret = self.lib.octl_distribute_file(str.encode(local_file_or_dir), str.encode(target_path), names_input, len(names), results_output, resultslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_execution_results_list(results_output, resultslen_output)
            raise Exception(bytes.decode(self.ebuf[:ret]))

        result_list = []
        for idx in range(resultslen_output.value):
                result_list.append(ExecutionResult(results_output[idx]))
        
        self.lib.octl_clear_execution_results_list(results_output, resultslen_output)
        return result_list

    def pull_file(self, ftype: str, name: str, remote_file_or_dir: str, local_dir: str) -> ExecutionResult:
        type_input = ctypes.c_int8(0)
        if ftype == 'fstore':
            type_input = ctypes.c_int8(0) 
        elif ftype == 'log':
            type_input = ctypes.c_int8(1)
        elif ftype == 'nodeapp':
            type_input = ctypes.c_int8(2)
        else:
            type_input = ctypes.c_int8(0)
        
        result_output = execution_result()
        ret = self.lib.octl_pull_file(type_input, str.encode(name), str.encode(remote_file_or_dir), str.encode(local_dir), result_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_execution_result(result_output)
            raise Exception(bytes.decode(self.ebuf[:ret]))
        
        ret_obj = ExecutionResult(result_output)
        self.lib.octl_clear_execution_result(result_output)
        return ret_obj

    def run(self, cmd_expr: str, names: List[str]) -> List[ExecutionResult]:
        if len(names) == 0:
            return
        output_len = len(names)
        names_input = (ctypes.c_char_p * len(names))()
        results_output = (execution_result * output_len)()
        resultslen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])
        
        ret = self.lib.octl_run(str.encode(cmd_expr), names_input, len(names), results_output, resultslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_execution_results_list(results_output, resultslen_output)
            raise Exception(bytes.decode(self.ebuf[:ret]))

        result_list = []
        for idx in range(resultslen_output.value):
                result_list.append(ExecutionResult(results_output[idx]))
        
        self.lib.octl_clear_execution_results_list(results_output, resultslen_output)
        return result_list

    def xrun(self, cmd_expr: str, names: List[str], delay: int) -> List[ExecutionResult]:
        if len(names) == 0:
            return
        output_len = len(names)
        names_input = (ctypes.c_char_p * len(names))()
        results_output = (execution_result * output_len)()
        resultslen_output = ctypes.c_int(output_len)
        for idx in range(len(names)):
            names_input[idx] = str.encode(names[idx])
        
        print("args", cmd_expr, names, delay, output_len)
        ret = self.lib.octl_xrun(str.encode(cmd_expr), names_input, len(names), ctypes.c_int(delay), results_output, resultslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_execution_results_list(results_output, resultslen_output)
            raise Exception(bytes.decode(self.ebuf[:ret]))

        result_list = []
        for idx in range(resultslen_output.value):
                result_list.append(ExecutionResult(results_output[idx]))
        
        self.lib.octl_clear_execution_results_list(results_output, resultslen_output)
        return result_list

    def get_groups_list(self) -> List[str]:
        output_len = 1024
        names_output = (ctypes.c_char_p * output_len)()
        nameslen_output = ctypes.c_int(output_len)
        ret = self.lib.octl_get_groups_list(names_output, nameslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_name_list(names_output, nameslen_output)
            raise Exception(bytes.decode(self.ebuf[:ret]))
        
        names_list = []
        for idx in range(nameslen_output.value):
            names_list.append(bytes.decode(names_output[idx]))

        self.lib.octl_clear_name_list(names_output, nameslen_output)
        return names_list
    
    def get_group(self, name: str) -> List[str]:
        output_len = 1024
        names_output = (ctypes.c_char_p * output_len)()
        nameslen_output = ctypes.c_int(output_len)
        ret = self.lib.octl_get_group(str.encode(name), names_output, nameslen_output, self.ebuf, self.ebuflen)
        if ret != 0:
            self.lib.octl_clear_name_list(names_output, nameslen_output)
            raise Exception(bytes.decode(self.ebuf[:ret]))
        
        names_list = []
        for idx in range(nameslen_output.value):
            names_list.append(bytes.decode(names_output[idx]))

        self.lib.octl_clear_name_list(names_output, nameslen_output)
        return names_list

    def set_group(self, name: str, skipCheck: bool, members: List[str]) -> None:
        if len(members) == 0:
            return
        
        members_input = (ctypes.c_char_p * len(members))()
        for idx in range(len(members)):
            members_input[idx] = str.encode(members[idx])
        
        skipcheck_input = ctypes.c_int(0)
        if skipCheck:
            skipcheck_input = ctypes.c_int(1)
        ret = self.lib.octl_set_group(str.encode(name), skipcheck_input, members_input, len(members), self.ebuf, self.ebuflen)
        if ret != 0:
            raise Exception(bytes.decode(self.ebuf[:ret]))

    def del_group(self, name: str) -> None:
        ret = self.lib.octl_del_group(str.encode(name), self.ebuf, self.ebuflen)
        if ret != 0:
            raise Exception(bytes.decode(self.ebuf[:ret]))

    def __str__(self):
        # 定义对象的字符串表示，可用于打印对象
        return f"octl python client"
