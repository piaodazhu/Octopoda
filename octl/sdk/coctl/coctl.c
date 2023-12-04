#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include "coctl.h"

int min(int a, int b) {
	if (a < b) {
		return a;
	}
	return b;
}

GoString make_GoString(char* s) {
	GoString gs = {.p = s, .n = strlen(s)};
	return gs;
}

GoSlice make_GoSlice(void *data, int n) {
	GoSlice gs = {.cap = n, .len = n, .data = data};
	return gs;
}

int
octl_init(char* config,
		char *errbuf, int *errbuflen) {
	struct Init_return ret = Init(make_GoString(config));
	int code = ret.r0;
	char *emsg = ret.r1;
	if (code == 0) {
		return 0;
	}
	*errbuflen = min(*errbuflen, strlen(emsg));
	memcpy(errbuf, emsg, *errbuflen);
	free(emsg);
	return code;
}

int
octl_get_node_info(char** names, int input_size, 
		octl_brain_info *output_obj, octl_node_info *output_list, int *output_size,
		char *errbuf, int *errbuflen) {
	GoString *name_strs = malloc(sizeof(GoString) * input_size);
	int i;
	for (i = 0; i < input_size; i++) {
		name_strs[i] = make_GoString(names[i]);
	}
	struct NodeInfo_return ret = NodeInfo(make_GoSlice(name_strs, input_size), output_obj, make_GoSlice(output_list, *output_size), output_size);
	free(name_strs);
	int code = ret.r0;
	char *emsg = ret.r1;
	if (code == 0) {
		return 0;
	}
	*errbuflen = min(*errbuflen, strlen(emsg));
	memcpy(errbuf, emsg, *errbuflen);
	free(emsg);
	return code;
}

int
octl_get_node_status(char** names, int input_size, 
		octl_node_status *output_list, int *output_size,
		char *errbuf, int *errbuflen) {
	GoString *name_strs = malloc(sizeof(GoString) * input_size);
	int i;
	for (i = 0; i < input_size; i++) {
		name_strs[i] = make_GoString(names[i]);
	}
	struct NodeStatus_return ret = NodeStatus(make_GoSlice(name_strs, input_size), make_GoSlice(output_list, *output_size), output_size);
	free(name_strs);
	int code = ret.r0;
	char *emsg = ret.r1;
	if (code == 0) {
		return 0;
	}
	*errbuflen = min(*errbuflen, strlen(emsg));
	memcpy(errbuf, emsg, *errbuflen);
	free(emsg);
	return code;
}

int
octl_upload_file(char* local_file_or_dir, char* remote_target_path, int is_force,
		char** names, int input_size, octl_execution_result *output_list, int *output_size,
		char *errbuf, int *errbuflen) {
	GoString *name_strs = malloc(sizeof(GoString) * input_size);
	int i;
	for (i = 0; i < input_size; i++) {
		name_strs[i] = make_GoString(names[i]);
	}
	struct UploadFile_return ret = UploadFile(make_GoString(local_file_or_dir), make_GoString(remote_target_path), (GoInt)is_force, make_GoSlice(name_strs, input_size), make_GoSlice(output_list, *output_size), output_size);
	free(name_strs);
	int code = ret.r0;
	char *emsg = ret.r1;
	if (code == 0) {
		return 0;
	}
	*errbuflen = min(*errbuflen, strlen(emsg));
	memcpy(errbuf, emsg, *errbuflen);
	free(emsg);
	return code;
}

int
octl_download_file(char* remote_file_or_dir, char* local_dir, char* name, 
		octl_execution_result *output_obj,
		char *errbuf, int *errbuflen) {
	struct DownloadFile_return ret = DownloadFile(make_GoString(remote_file_or_dir), make_GoString(local_dir), make_GoString(name), output_obj);
	int code = ret.r0;
	char *emsg = ret.r1;
	if (code == 0) {
		return 0;
	}
	*errbuflen = min(*errbuflen, strlen(emsg));
	memcpy(errbuf, emsg, *errbuflen);
	free(emsg);
	return code;
}

int octl_run(char *cmd, int runtype, int need_align, int delay, 
		char **names, int input_size, 
		octl_execution_result *output_list, int *output_size,
		char *errbuf, int *errbuflen) {
	GoString *name_strs = malloc(sizeof(GoString) * input_size);
	int i;
	for (i = 0; i < input_size; i++) {
		name_strs[i] = make_GoString(names[i]);
	}
	
	int code;
	char *emsg;
	if (delay < 0) { // must be run 
		if (runtype == 0) {
			struct RunCommand_return ret = RunCommand(make_GoString(cmd), (GoInt)(need_align), make_GoSlice(name_strs, input_size), make_GoSlice(output_list, *output_size), output_size);
			code = ret.r0;
			emsg = ret.r1;
		} else if (runtype == 1) {
			struct RunScript_return ret = RunScript(make_GoString(cmd), (GoInt)(need_align), make_GoSlice(name_strs, input_size), make_GoSlice(output_list, *output_size), output_size);
			code = ret.r0;
			emsg = ret.r1;
		} else {
			struct RunCommandBackground_return ret = RunCommandBackground(make_GoString(cmd), (GoInt)(need_align), make_GoSlice(name_strs, input_size), make_GoSlice(output_list, *output_size), output_size);
			code = ret.r0;
			emsg = ret.r1;
		}
	} else { // must be xrun
		if (runtype == 0) {
			struct XRunCommand_return ret = XRunCommand(make_GoString(cmd), (GoInt)(need_align), (GoInt)(delay), make_GoSlice(name_strs, input_size), make_GoSlice(output_list, *output_size), output_size);
			code = ret.r0;
			emsg = ret.r1;
		} else {
			struct XRunScript_return ret = XRunScript(make_GoString(cmd), (GoInt)(need_align), (GoInt)(delay), make_GoSlice(name_strs, input_size), make_GoSlice(output_list, *output_size), output_size);
			code = ret.r0;
			emsg = ret.r1;
		}
	}
	free(name_strs);
	if (code == 0) {
		return 0;
	}
	*errbuflen = min(*errbuflen, strlen(emsg));
	memcpy(errbuf, emsg, *errbuflen);
	free(emsg);
	return code;
}

int
octl_run_command(char *cmd_str, int need_align, char **names, int input_size, 
		octl_execution_result *output_list, int *output_size,
		char *errbuf, int *errbuflen) {
	return octl_run(cmd_str, 0, need_align, -1, names, input_size, output_list, output_size, errbuf, errbuflen);
}

int octl_run_script(char *script_file, int need_align, char **names, int input_size, 
		octl_execution_result *output_list, int *output_size,
		char *errbuf, int *errbuflen) {
	return octl_run(script_file, 1, need_align, -1, names, input_size, output_list, output_size, errbuf, errbuflen);
}

int octl_run_command_background(char *cmd_str, int need_align, char **names, int input_size, 
		octl_execution_result *output_list, int *output_size,
		char *errbuf, int *errbuflen) {
	return octl_run(cmd_str, 2, need_align, -1, names, input_size, output_list, output_size, errbuf, errbuflen);
}

int octl_xrun_command(char *cmd_str, int need_align, int delay, 
		char **names, int input_size, 
		octl_execution_result *output_list, int *output_size,
		char *errbuf, int *errbuflen) {
	return octl_run(cmd_str, 0, need_align, delay, names, input_size, output_list, output_size, errbuf, errbuflen);
}

int octl_xrun_script(char *script_file, int need_align, int delay, 
		char **names, int input_size, 
		octl_execution_result *output_list, int *output_size,
		char *errbuf, int *errbuflen) {
	return octl_run(script_file, 1, need_align, delay, names, input_size, output_list, output_size, errbuf, errbuflen);
}

// int
// octl_get_groups_list(char **output_list, int *output_size,
// 		char *errbuf, int *errbuflen) {
// 	struct GroupGetAll_return ret = GroupGetAll(make_GoSlice(output_list, *output_size), output_size);
// 	int code = ret.r0;
// 	char *emsg = ret.r1;
// 	if (code == 0) {
// 		return 0;
// 	}
// 	*errbuflen = min(*errbuflen, strlen(emsg));
// 	memcpy(errbuf, emsg, *errbuflen);
// 	free(emsg);
// 	return code;
// }

// int
// octl_get_group(char *group_name, char **output_list, int *output_size,
// 		char *errbuf, int *errbuflen) {
// 	struct GroupGet_return ret = GroupGet(make_GoString(group_name), make_GoSlice(output_list, *output_size), output_size);
// 	int code = ret.r0;
// 	char *emsg = ret.r1;
// 	if (code == 0) {
// 		return 0;
// 	}
// 	*errbuflen = min(*errbuflen, strlen(emsg));
// 	memcpy(errbuf, emsg, *errbuflen);
// 	free(emsg);
// 	return code;
// }

// int
// octl_set_group(char *group_name, int skipCheck, char **names, int input_size,
// 		char *errbuf, int *errbuflen) {
// 	GoString *name_strs = malloc(sizeof(GoString) * input_size);
// 	int i;
// 	for (i = 0; i < input_size; i++) {
// 		name_strs[i] = make_GoString(names[i]);
// 	}
// 	struct GroupSet_return ret = GroupSet(make_GoString(group_name), skipCheck, make_GoSlice(name_strs, input_size));
// 	int code = ret.r0;
// 	char *emsg = ret.r1;
// 	if (code == 0) {
// 		return 0;
// 	}
// 	*errbuflen = min(*errbuflen, strlen(emsg));
// 	memcpy(errbuf, emsg, *errbuflen);
// 	free(emsg);
// 	return code;
// }

// int
// octl_del_group(char *group_name,
// 		char *errbuf, int *errbuflen) {
// 	struct GroupDel_return ret = GroupDel(make_GoString(group_name));
// 	int code = ret.r0;
// 	char *emsg = ret.r1;
// 	if (code == 0) {
// 		return 0;
// 	}
// 	*errbuflen = min(*errbuflen, strlen(emsg));
// 	memcpy(errbuf, emsg, *errbuflen);
// 	free(emsg);
// 	return code;
// }


// int octl_prune_nodes(char *errbuf, int *errbuflen) {
// 	struct Prune_return ret = Prune();
// 	int code = ret.r0;
// 	char *emsg = ret.r1;
// 	if (code == 0) {
// 		return 0;
// 	}
// 	*errbuflen = min(*errbuflen, strlen(emsg));
// 	memcpy(errbuf, emsg, *errbuflen);
// 	free(emsg);
// 	return code;
// }


int octl_get_scenarios_info_list(char **output_list, int *output_size,
		char *errbuf, int *errbuflen) {
	struct ScenariosInfo_return ret = ScenariosInfo(make_GoSlice(output_list, *output_size), output_size);
	int code = ret.r0;
	char *emsg = ret.r1;
	if (code == 0) {
		return 0;
	}
	*errbuflen = min(*errbuflen, strlen(emsg));
	memcpy(errbuf, emsg, *errbuflen);
	free(emsg);
	return code;
}


int octl_get_scenario_info(char *name, char *output_buf, int *output_size,
		char *errbuf, int *errbuflen) {
	struct ScenarioInfo_return ret = ScenarioInfo(make_GoString(name));
	char *scenBuf = ret.r0;
	int scenLen = strlen(scenBuf);
	int code = ret.r1;
	char *emsg = ret.r2;
	if (code == 0) {
		if (scenLen <= *output_size) {
			*output_size = scenLen;
			memcpy(output_buf, scenBuf, scenLen);
			free(scenBuf);
			return 0;
		}
		code = OctlSdkBufferError;
		emsg = "buf size not enough";
	}
	*errbuflen = min(*errbuflen, strlen(emsg));
	memcpy(errbuf, emsg, *errbuflen);
	if (code != OctlSdkBufferError)
		free(emsg);
	free(scenBuf);
	return code;
}


int octl_get_scenario_version(char *name, int offset, int limit, char *output_buf, int *output_size,
		char *errbuf, int *errbuflen) {
	struct ScenarioVersion_return ret = ScenarioVersion(make_GoString(name), (GoInt)(offset), (GoInt)(limit));
	char *verBuf = ret.r0;
	int verLen = strlen(verBuf);
	int code = ret.r1;
	char *emsg = ret.r2;
	if (code == 0) {
		if (verLen <= *output_size) {
			*output_size = verLen;
			memcpy(output_buf, verBuf, verLen);
			free(verBuf);
			return 0;
		}
		code = OctlSdkBufferError;
		emsg = "buf size not enough";
	}
	*errbuflen = min(*errbuflen, strlen(emsg));
	memcpy(errbuf, emsg, *errbuflen);
	if (code != OctlSdkBufferError)
		free(emsg);
	free(verBuf);
	return code;
}


int octl_get_nodeapp_info(char *name, char *app, char *scenario, 
		char *output_buf, int *output_size,
		char *errbuf, int *errbuflen) {
	struct NodeAppInfo_return ret = NodeAppInfo(make_GoString(name), make_GoString(app), make_GoString(scenario));
	char *nodeappBuf = ret.r0;
	int nodeappLen = strlen(nodeappBuf);
	int code = ret.r1;
	char *emsg = ret.r2;
	if (code == 0) {
		if (nodeappLen <= *output_size) {
			*output_size = nodeappLen;
			memcpy(output_buf, nodeappBuf, nodeappLen);
			free(nodeappBuf);
			return 0;
		}
		code = OctlSdkBufferError;
		emsg = "buf size not enough";
	}
	*errbuflen = min(*errbuflen, strlen(emsg));
	memcpy(errbuf, emsg, *errbuflen);
	if (code != OctlSdkBufferError)
		free(emsg);
	free(nodeappBuf);
	return code;
}


int octl_get_nodeapps_info_list(char *name, char **output_list, int *output_size,
		char *errbuf, int *errbuflen) {
	struct NodeAppsInfo_return ret = NodeAppsInfo(make_GoString(name), make_GoSlice(output_list, *output_size), output_size);
	int code = ret.r0;
	char *emsg = ret.r1;
	if (code == 0) {
		return 0;
	}
	*errbuflen = min(*errbuflen, strlen(emsg));
	memcpy(errbuf, emsg, *errbuflen);
	free(emsg);
	return code;
}


int octl_apply_scenario(char *name, char *target, char *message, 
		int timeout, char **log_list, int *log_size,
		char *errbuf, int *errbuflen) {
	struct Apply_return ret = Apply(make_GoString(name), make_GoString(target), make_GoString(message), timeout, make_GoSlice(log_list, *log_size), log_size);
	int code = ret.r0;
	char *emsg = ret.r1;
	if (code == 0) {
		return 0;
	}
	*errbuflen = min(*errbuflen, strlen(emsg));
	memcpy(errbuf, emsg, *errbuflen);
	free(emsg);
	return code;
}

void
octl_clear_node_info(octl_node_info *obj) {
	if (obj->Name != NULL)
		free(obj->Name);
	if (obj->Version != NULL)
		free(obj->Version);
	if (obj->Addr != NULL)
		free(obj->Addr);
	memset(obj, 0, sizeof(octl_node_info));
}

void
octl_clear_node_status(octl_node_status *obj) {
	if (obj->Name != NULL)
		free(obj->Name);
	if (obj->Platform != NULL)
		free(obj->Platform);
	memset(obj, 0, sizeof(octl_node_status));
}

void
octl_clear_brain_info(octl_brain_info *obj) {
	if (obj->Name != NULL)
		free(obj->Name);
	if (obj->Version != NULL)
		free(obj->Version);
	if (obj->Addr != NULL)
		free(obj->Addr);
	memset(obj, 0, sizeof(octl_brain_info));
}

void
octl_clear_execution_result(octl_execution_result *obj) {
	if (obj->Name != NULL)
		free(obj->Name);
	if (obj->CommunicationErrorMsg != NULL)
		free(obj->CommunicationErrorMsg);
	if (obj->ProcessErrorMsg != NULL)
		free(obj->ProcessErrorMsg);
	if (obj->Result != NULL)
		free(obj->Result);
	memset(obj, 0, sizeof(octl_execution_result));
}

void
octl_clear_node_info_list(octl_node_info *list, int n) {
	int i;
	for (i = 0; i < n; i++)
		octl_clear_node_info(&list[i]);
}

void
octl_clear_string_list(char **list, int n) {
	int i;
	for (i = 0; i < n; i++)
		free(list[i]);
}

void
octl_clear_node_status_list(octl_node_status *list, int n) {
	int i;
	for (i = 0; i < n; i++)
		octl_clear_node_status(&list[i]);
}

void
octl_clear_execution_result_list(octl_execution_result *list, int n) {
	int i;
	for (i = 0; i < n; i++)
		octl_clear_execution_result(&list[i]);
}
