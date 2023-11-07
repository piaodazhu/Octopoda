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
		char *errbuf, int errbuflen) {
	char* ret = Init(make_GoString(config));
	if (ret == NULL) {
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	free(ret);
	return len;
}

int
octl_get_node_info(char* name, octl_node_info *output_obj, 
		char *errbuf, int errbuflen) {
	char* ret = NodeInfo(make_GoString(name), output_obj);
	if (ret == NULL) {
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	free(ret);
	return len;
}

int
octl_get_nodes_info_list(char** names, int input_size, 
		octl_brain_info *output_obj, octl_node_info *output_list, int *output_size,
		char *errbuf, int errbuflen) {
	GoString *name_strs = malloc(sizeof(GoString) * input_size);
	int i;
	for (i = 0; i < input_size; i++) {
		name_strs[i] = make_GoString(names[i]);
	}
	char* ret = NodesInfo(make_GoSlice(name_strs, input_size), output_obj, make_GoSlice(output_list, *output_size), output_size);
	if (ret == NULL) {
		free(name_strs);
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	free(ret);
	free(name_strs);
	return len;
}

int
octl_get_node_status(char* name, octl_node_status *output_obj,
		char* errbuf, int errbuflen) {
	char* ret = NodeStatus(make_GoString(name), output_obj);
	if (ret == NULL) {
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	free(ret);
	return len;
}

int
octl_get_nodes_status_list(char** names, int input_size, 
		octl_node_status *output_list, int *output_size,
		char *errbuf, int errbuflen) {
	GoString *name_strs = malloc(sizeof(GoString) * input_size);
	int i;
	for (i = 0; i < input_size; i++) {
		name_strs[i] = make_GoString(names[i]);
	}
	char* ret = NodesStatus(make_GoSlice(name_strs, input_size), make_GoSlice(output_list, *output_size), output_size);
	if (ret == NULL) {
		free(name_strs);
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	free(name_strs);
	free(ret);
	return len;
}

int
octl_distribute_file(char* local_file_or_dir, char* target_path, 
		char** names, int input_size, octl_execution_result *output_list, int *output_size,
		char *errbuf, int errbuflen) {
	GoString *name_strs = malloc(sizeof(GoString) * input_size);
	int i;
	for (i = 0; i < input_size; i++) {
		name_strs[i] = make_GoString(names[i]);
	}
	char* ret = DistribFile(make_GoString(local_file_or_dir), make_GoString(target_path), make_GoSlice(name_strs, input_size), make_GoSlice(output_list, *output_size), output_size);
	if (ret == NULL) {
		free(name_strs);
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	free(name_strs);
	free(ret);
	return len;
}

int
octl_pull_file(enum PATHTYPE type, char* name, 
		char* remote_file_or_dir, char* local_dir, 
		octl_execution_result *output_obj,
		char *errbuf, int errbuflen) {
	GoString type_str;
	switch (type)
	{
	case FSTORE:
		type_str = make_GoString("store");
		break;
	case LOG:
		type_str = make_GoString("log");
		break;
	case NODEAPP:
		type_str = make_GoString("nodeapp");
		break;
	default:
		break;
	}
	
	char* ret = PullFile(type_str, make_GoString(name), make_GoString(remote_file_or_dir), make_GoString(local_dir), output_obj);
	if (ret == NULL) {
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	free(ret);
	return len;
}

int
octl_run(char *cmd_expr, char **names, int input_size, 
	octl_execution_result *output_list, int *output_size,
	char *errbuf, int errbuflen) {
	GoString *name_strs = malloc(sizeof(GoString) * input_size);
	int i;
	for (i = 0; i < input_size; i++) {
		name_strs[i] = make_GoString(names[i]);
	}
	char* ret = Run(make_GoString(cmd_expr), make_GoSlice(name_strs, input_size), make_GoSlice(output_list, 2), output_size);
	if (ret == NULL) {
		free(name_strs);
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	free(name_strs);
	free(ret);
	return len;
}

int
octl_xrun(char *cmd_expr, char **names, int input_size, int delay, 
	octl_execution_result *output_list, int *output_size,
	char *errbuf, int errbuflen) {
	GoString *name_strs = malloc(sizeof(GoString) * input_size);
	int i;
	for (i = 0; i < input_size; i++) {
		name_strs[i] = make_GoString(names[i]);
	}
	char* ret = XRun(make_GoString(cmd_expr), make_GoSlice(name_strs, input_size), (GoInt)(delay), make_GoSlice(output_list, *output_size), output_size);
	if (ret == NULL) {
		free(name_strs);
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	free(name_strs);
	free(ret);
	return len;
}

int
octl_get_groups_list(char **output_list, int *output_size,
		char *errbuf, int errbuflen) {
	// GoString *str_slice = malloc(sizeof(GoString) * (*output_size));
	char* ret = GroupGetAll(make_GoSlice(output_list, *output_size), output_size);
	if (ret == NULL) {
		// int i;
		// for (i = 0; i < *output_size; i++) {
		// 	output_list[i] = str_slice[i].p;
		// }
		// free(str_slice);
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	// free(str_slice);
	free(ret);
	return len;
}

int
octl_get_group(char *group_name, char **output_list, int *output_size,
		char *errbuf, int errbuflen) {
	// GoString *str_slice = malloc(sizeof(GoString) * (*output_size));
	char* ret = GroupGet(make_GoString(group_name), make_GoSlice(output_list, *output_size), output_size);
	if (ret == NULL) {
		// int i;
		// for (i = 0; i < *output_size; i++) {
		// 	output_list[i] = str_slice[i].p;
		// }
		// free(str_slice);
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	// free(str_slice);
	free(ret);
	return len;
}

int
octl_set_group(char *group_name, int skipCheck, char **names, int input_size,
	char *errbuf, int errbuflen) {
	GoString *name_strs = malloc(sizeof(GoString) * input_size);
	int i;
	for (i = 0; i < input_size; i++) {
		name_strs[i] = make_GoString(names[i]);
	}
	char* ret = GroupSet(make_GoString(group_name), skipCheck, make_GoSlice(name_strs, input_size));
	if (ret == NULL) {
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	free(ret);
	return len;
}

int
octl_del_group(char *group_name,
	char *errbuf, int errbuflen) {
	char* ret = GroupDel(make_GoString(group_name));
	if (ret == NULL) {
		return 0;
	}
	int len = min(errbuflen, strlen(ret));
	memcpy(errbuf, ret, len);
	free(ret);
	return len;
}

void
octl_clear_node_info(octl_node_info *obj) {
	if (obj->Name != NULL)
		free(obj->Name);
	if (obj->Version != NULL)
		free(obj->Version);
	if (obj->Addr != NULL)
		free(obj->Addr);
	if (obj->ConnState != NULL)
		free(obj->ConnState);
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
octl_clear_nodes_info_list(octl_node_info *list, int n) {
	int i;
	for (i = 0; i < n; i++)
		octl_clear_node_info(&list[i]);
}

void
octl_clear_name_list(char **list, int n) {
	int i;
	for (i = 0; i < n; i++)
		free(list[i]);
}

void
octl_clear_nodes_status_list(octl_node_status *list, int n) {
	int i;
	for (i = 0; i < n; i++)
		octl_clear_node_status(&list[i]);
}

void
octl_clear_execution_results_list(octl_execution_result *list, int n) {
	int i;
	for (i = 0; i < n; i++)
		octl_clear_execution_result(&list[i]);
}
