
#include <stdio.h>
#include <stdlib.h>
#include "coctl.h"

#define EBUF_LEN 256
int main() {
	char ebuf[EBUF_LEN];
	int elen = EBUF_LEN;
	int ret;
	ret = octl_init("../../octl_test.yaml", ebuf, &elen);
	if (ret > 0) {
		printf("octl_init: %.*s\n", elen, ebuf);
		return 1;
	}

	octl_node_info ninfo;
	elen = EBUF_LEN;
	ret = octl_get_node_info("pi4", &ninfo, ebuf, &elen);
	if (ret > 0) {
		printf("octl_get_node_info: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("node %s addr=%s state=%d connstate=%s\n", ninfo.Name, ninfo.Addr, ninfo.State, ninfo.ConnState);
	octl_clear_node_info(&ninfo);

	octl_node_status nstatus;
	elen = EBUF_LEN;
	ret = octl_get_node_status("yang", &nstatus, ebuf, &elen);
	if (ret > 0) {
		printf("octl_get_node_status: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("node %s platform=%s cpu=%f mem=%lld\n", nstatus.Name, nstatus.Platform, nstatus.CpuLoadShort, nstatus.MemUsed);
	octl_clear_node_status(&nstatus);

	octl_brain_info binfo;
	octl_node_info all[32];
	int total_num = 32;
	elen = EBUF_LEN;
	ret = octl_get_nodes_info_list(NULL, 0, &binfo, all, &total_num, ebuf, &elen);
	if (ret > 0) {
		printf("octl_get_nodes_info_list: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("brain %s addr=%s version=%s nodenum=%d\n",binfo.Name, binfo.Addr, binfo.Version, total_num);
	octl_clear_brain_info(&binfo);
	octl_clear_nodes_info_list(all, total_num);

	octl_execution_result results[2];
	char *targets[2] = {"pi4", "pi5"};
	int total_results = 2;
	elen = EBUF_LEN;
	ret = octl_run("{uname -a}", targets, 2, 0, results, &total_results, ebuf, &elen);
	if (ret > 0) {
		printf("octl_run: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("result of %s code=%d output=%s", results[0].Name, results[0].Code, results[0].Result);
	printf("result of %s code=%d output=%s", results[1].Name, results[1].Code, results[1].Result);
	octl_clear_execution_results_list(results, total_results);

	total_results = 2;
	elen = EBUF_LEN;
	ret = octl_xrun("{uname -a}", targets, 2, 1, 0, results, &total_results, ebuf, &elen);
	if (ret > 0) {
		printf("octl_run: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("xrun result of %s code=%d output=%s\n", results[0].Name, results[0].Code, results[0].Result);
	printf("xrun result of %s code=%d output=%s\n", results[1].Name, results[1].Code, results[1].Result);
	octl_clear_execution_results_list(results, total_results);

	char *gnames[32];
	int total_groups = 32;
	elen = EBUF_LEN;
	ret = octl_get_groups_list(gnames, &total_groups, ebuf, &elen);
	if (ret > 0) {
		printf("octl_get_groups_list: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("total groups count=%d\n", total_groups);
	
	int i, j;
	for (i = 0; i < total_groups; i++) {
		printf("group[%d]: %s\n  ", i, gnames[i]);
		
		char *mnames[32];
		int total_members = 32;
		elen = EBUF_LEN;
		ret = octl_get_group(gnames[i], mnames, &total_members, ebuf, &elen);
		if (ret > 0) {
			printf("octl_get_group: %.*s\n", elen, ebuf);
			continue;
		}
		
		for (j = 0; j < total_members; j++) {
			printf("%s, ", mnames[j]);
		}
		printf("\n");

		octl_clear_string_list(mnames, total_members);
	}

	octl_clear_string_list(gnames, total_groups);

	char *scenarios[64];
	int total_scenarios = 64;
	elen = EBUF_LEN;
	ret = octl_get_scenarios_info_list(scenarios, &total_scenarios, ebuf, &elen);
	if (ret > 0) {
		printf("octl_get_scenarios_info_list: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("total scenarios count=%d\n", total_scenarios);
	for (i = 0; i < total_scenarios; i++) {
		printf("scenario[%d]: %s\n", i, scenarios[i]);
	}

	char rawBuf[1024];
	int rawLen = 1024;
	ret = octl_get_scenario_info("s1", rawBuf, &rawLen, ebuf, &elen);
	if (ret > 0) {
		printf("octl_get_scenario_info: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("scenario s1: %.*s\n", rawLen, rawBuf);

	rawLen = 1024;
	ret = octl_get_scenario_version("s1", rawBuf, &rawLen, ebuf, &elen);
	if (ret > 0) {
		printf("octl_get_scenario_version: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("version s1: %.*s\n", rawLen, rawBuf);

	char *nodeapps[64];
	int total_nodeapps = 64;
	elen = EBUF_LEN;
	ret = octl_get_nodeapps_info_list("pi4", nodeapps, &total_nodeapps, ebuf, &elen);
	if (ret > 0) {
		printf("octl_get_nodeapps_info_list: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("total nodeapps count=%d\n", total_nodeapps);
	for (i = 0; i < total_nodeapps; i++) {
		printf("nodeapps[%d]: %s\n", i, nodeapps[i]);
	}

	char *logs[64];
	int total_logs = 64;
	elen = EBUF_LEN;
	ret = octl_apply_scenario("../../s1", "start", "bycoctl", 10, logs, &total_logs, ebuf, &elen);
	if (ret > 0) {
		printf("octl_apply_scenario: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("total logs count=%d\n", total_logs);

	printf("PASS\n");
	return 0;
}
