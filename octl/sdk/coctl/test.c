
#include <stdio.h>
#include <stdlib.h>
#include "coctl.h"

#define EBUF_LEN 256
int main() {
	char ebuf[EBUF_LEN];
	int elen;
	elen = octl_init("../../octl_test.yaml", ebuf, EBUF_LEN);
	if (elen > 0) {
		printf("octl_init: %.*s\n", elen, ebuf);
		return 1;
	}

	octl_node_info ninfo;
	elen = octl_get_node_info("pi4", &ninfo, ebuf, EBUF_LEN);
	if (elen > 0) {
		printf("octl_get_node_info: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("node %s addr=%s state=%d connstate=%s\n", ninfo.Name, ninfo.Addr, ninfo.State, ninfo.ConnState);

	octl_node_status nstatus;
	elen = octl_get_node_status("yang", &nstatus, ebuf, EBUF_LEN);
	if (elen > 0) {
		printf("octl_get_node_status: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("node %s platform=%s cpu=%f mem=%lld\n", nstatus.Name, nstatus.Platform, nstatus.CpuLoadShort, nstatus.MemUsed);

	octl_brain_info binfo;
	octl_node_info all[32];
	int total_num = 32;
	elen = octl_get_nodes_info_list(NULL, 0, &binfo, all, &total_num, ebuf, EBUF_LEN);
	if (elen > 0) {
		printf("octl_get_nodes_info_list: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("brain %s addr=%s version=%s nodenum=%d\n",binfo.Name, binfo.Addr, binfo.Version, total_num);

	octl_execution_result results[2];
	char *targets[2] = {"pi4", "pi5"};
	int total_results = 2;
	elen = octl_run("{uname -a}", targets, 2, results, &total_results, ebuf, EBUF_LEN);
	if (elen > 0) {
		printf("octl_run: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("result of %s code=%d output=%s", results[0].Name, results[0].Code, results[0].Result);
	printf("result of %s code=%d output=%s", results[1].Name, results[1].Code, results[1].Result);

	total_results = 2;
	elen = octl_xrun("{uname -a}", targets, 2, 1, results, &total_results, ebuf, EBUF_LEN);
	if (elen > 0) {
		printf("octl_run: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("xrun result of %s code=%d output=%s\n", results[0].Name, results[0].Code, results[0].Result);
	printf("xrun result of %s code=%d output=%s\n", results[1].Name, results[1].Code, results[1].Result);

	char *gnames[32];
	int total_groups = 32;
	elen = octl_get_groups_list(gnames, &total_groups, ebuf, EBUF_LEN);
	if (elen > 0) {
		printf("octl_get_groups_list: %.*s\n", elen, ebuf);
		return 1;
	}
	printf("total groups count=%d\n", total_groups);
	int i, j;
	for (i = 0; i < total_groups; i++) {
		printf("group[%d]: %s\n  ", i, gnames[i]);
		
		char *mnames[32];
		int total_members = 32;
		elen = octl_get_group(gnames[i], mnames, &total_members, ebuf, EBUF_LEN);
		if (elen > 0) {
			printf("octl_get_group: %.*s\n", elen, ebuf);
			continue;
		}
		
		for (j = 0; j < total_members; j++) {
			printf("%s, ", mnames[j]);
		}
		printf("\n");
	}


	printf("PASS\n");
	return 0;
}
