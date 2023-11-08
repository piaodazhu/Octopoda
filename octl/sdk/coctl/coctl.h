#ifndef _COCTL_CLIENT_H
#define _COCTL_CLIENT_H
#include "wrapper.h"

#if defined(WIN32) || defined(_WIN32) || defined(_WIN32_) || defined(WIN64) || defined(_WIN64) || defined(_WIN64_)
#define WIN_DLL_EXPORT
#endif

enum PATHTYPE {
	FSTORE,
	LOG,
	NODEAPP,
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
		char *errbuf, int errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_get_node_info(char* name, octl_node_info *output_obj, 
		char *errbuf, int errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_get_nodes_info_list(char** names, int input_size, 
		octl_brain_info *output_obj, octl_node_info *output_list, int *output_size,
		char *errbuf, int errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif 
int octl_get_node_status(char* name, octl_node_status *output_obj,
		char* errbuf, int errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif 
int octl_get_nodes_status_list(char** names, int input_size, 
		octl_node_status *output_list, int *output_size,
		char *errbuf, int errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif 
int octl_distribute_file(char* local_file_or_dir, char* target_path, 
		char** names, int input_size, octl_execution_result *output_list, int *output_size,
		char *errbuf, int errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_pull_file(enum PATHTYPE type, char* name, 
		char* remote_file_or_dir, char* local_dir, 
		octl_execution_result *output_obj,
		char *errbuf, int errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_run(char *cmd_expr, char **names, int input_size, 
	octl_execution_result *output_list, int *output_size,
	char *errbuf, int errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_xrun(char *cmd_expr, char **names, int input_size, int delay, 
	octl_execution_result *output_list, int *output_size,
	char *errbuf, int errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_get_groups_list(char **output_list, int *output_size,
	char *errbuf, int errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_get_group(char *group_name, char **output_list, int *output_size,
	char *errbuf, int errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_set_group(char *group_name, int skipCheck, char **names, int input_size,
		char *errbuf, int errbuflen);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
int octl_del_group(char *group_name,
	char *errbuf, int errbuflen);

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
void octl_clear_nodes_info_list(octl_node_info *list, int n);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
void octl_clear_nodes_status_list(octl_node_status *list, int n);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
void octl_clear_execution_results_list(octl_execution_result *list, int n);

#ifdef WIN_DLL_EXPORT
__declspec(dllexport) 
#endif
void octl_clear_name_list(char **list, int n);


#ifdef __cplusplus
}
#endif
#endif