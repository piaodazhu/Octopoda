#ifndef _COCTL_CLIENT_H
#define _COCTL_CLIENT_H
#include "wrapper.h"

#if defined(WIN32) || defined(_WIN32) || defined(_WIN32_) || defined(WIN64) || defined(_WIN64) || defined(_WIN64_)
#define WIN_DLL_EXPORT
#endif

enum ERRORTYPE {
	OK,
	OctlReadConfigError,
	OctlWriteConfigError,
	OctlInitClientError,
	OctlWorkgroupAuthError,
	OctlHttpRequestError,
	OctlHttpStatusError,
	OctlMessageParseError,
	OctlNodeParseError,
	OctlFileOperationError,
	OctlGitOperationError,
	OctlTaskWaitingError,
	OctlArgumentError,
	OctlSdkNotInitializedError,
	OctlSdkPanicRecoverError,
	OctlSdkBufferError,
	OctlContextCancelError,
};


typedef struct node_info octl_node_info;
typedef struct brain_info octl_brain_info;
typedef struct node_status octl_node_status;
typedef struct execution_result octl_execution_result;

#ifdef __cplusplus
extern "C" {
#endif

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_init(char* config,
		char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_get_node_info(char** names, int input_size, 
		octl_brain_info *output_obj, octl_node_info *output_list, int *output_size,
		char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif 
int octl_get_node_status(char** names, int input_size, 
		octl_node_status *output_list, int *output_size,
		char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif 
int octl_upload_file(char* local_file_or_dir, char* remote_target_path, int is_force,
		char** names, int input_size, octl_execution_result *output_list, int *output_size,
		char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_download_file(char* remote_file_or_dir, char* local_dir, char* name, 
		octl_execution_result *output_obj,
		char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_run_command(char *cmd_str, int need_align, char **names, int input_size, 
	octl_execution_result *output_list, int *output_size,
	char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_run_script(char *script_file, int need_align, char **names, int input_size, 
	octl_execution_result *output_list, int *output_size,
	char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_run_command_background(char *cmd_str, int need_align, char **names, int input_size, 
	octl_execution_result *output_list, int *output_size,
	char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_xrun_command(char *cmd_str, int need_align, int delay, 
	char **names, int input_size, 
	octl_execution_result *output_list, int *output_size,
	char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_xrun_script(char *script_file, int need_align, int delay, 
	char **names, int input_size, 
	octl_execution_result *output_list, int *output_size,
	char *errbuf, int *errbuflen);

// #ifdef WIN_DLL_EXPORT
// __declspec(dllexport) 
// #endif
// int octl_get_groups_list(char **output_list, int *output_size,
// 	char *errbuf, int *errbuflen);

// #ifdef WIN_DLL_EXPORT
// __declspec(dllexport) 
// #endif
// int octl_get_group(char *group_name, char **output_list, int *output_size,
// 	char *errbuf, int *errbuflen);

// #ifdef WIN_DLL_EXPORT
// __declspec(dllexport) 
// #endif
// int octl_set_group(char *group_name, int skipCheck, char **names, int input_size,
// 		char *errbuf, int *errbuflen);

// #ifdef WIN_DLL_EXPORT
// __declspec(dllexport) 
// #endif
// int octl_del_group(char *group_name,
// 	char *errbuf, int *errbuflen);

// #ifdef WIN_DLL_EXPORT
// __declspec(dllexport) 
// #endif
// int octl_prune_nodes(char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_get_scenarios_info_list(char **output_list, int *output_size,
	char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_get_scenario_info(char *name, char *output_buf, int *output_size,
	char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_get_scenario_version(char *name, int offset, int limit, 
	char *output_buf, int *output_size,
	char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_get_nodeapp_info(char *name, char *app, char *scenario, 
	char *output_buf, int *output_size,
	char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_get_nodeapps_info_list(char *name, char **output_list, int *output_size,
	char *errbuf, int *errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_apply_scenario(char *name, char *target, char *message, 
	int timeout, char **log_list, int *log_size,
	char *errbuf, int *errbuflen);


#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif 
void octl_clear_node_info(octl_node_info *obj);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
void octl_clear_node_status(octl_node_status *obj);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
void octl_clear_brain_info(octl_brain_info *obj);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
void octl_clear_execution_result(octl_execution_result *obj);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
void octl_clear_node_info_list(octl_node_info *list, int n);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
void octl_clear_node_status_list(octl_node_status *list, int n);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
void octl_clear_execution_result_list(octl_execution_result *list, int n);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
void octl_clear_string_list(char **list, int n);


#ifdef __cplusplus
}
#endif
#endif