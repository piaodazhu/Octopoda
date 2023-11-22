/* Code generated by cmd/cgo; DO NOT EDIT. */

/* package command-line-arguments */


#line 1 "cgo-builtin-export-prolog"

#include <stddef.h> /* for ptrdiff_t below */

#ifndef GO_CGO_EXPORT_PROLOGUE_H
#define GO_CGO_EXPORT_PROLOGUE_H

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef struct { const char *p; ptrdiff_t n; } _GoString_;
#endif

#endif

/* Start of preamble from import "C" comments.  */


#line 3 "wrapper.go"

struct node_info {
    char* Name;
    char* Version;
    char* Addr;
    int State;
    char* ConnState;
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

#line 1 "cgo-generated-wrapper"


/* End of preamble from import "C" comments.  */


/* Start of boilerplate cgo prologue.  */
#line 1 "cgo-gcc-export-header-prolog"

#ifndef GO_CGO_PROLOGUE_H
#define GO_CGO_PROLOGUE_H

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef __SIZE_TYPE__ GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;
typedef float _Complex GoComplex64;
typedef double _Complex GoComplex128;

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt.
*/
typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef _GoString_ GoString;
#endif
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;

#endif

/* End of boilerplate cgo prologue.  */

#ifdef __cplusplus
extern "C" {
#endif


/* Return type for Init */
struct Init_return {
	int r0;
	char* r1;
};
extern struct Init_return Init(GoString configFile);

/* Return type for NodeInfo */
struct NodeInfo_return {
	int r0;
	char* r1;
};
extern struct NodeInfo_return NodeInfo(GoString name, struct node_info* result);

/* Return type for NodesInfo */
struct NodesInfo_return {
	int r0;
	char* r1;
};
extern struct NodesInfo_return NodesInfo(GoSlice names, struct brain_info* brain, GoSlice results, int* size);

/* Return type for NodeStatus */
struct NodeStatus_return {
	int r0;
	char* r1;
};
extern struct NodeStatus_return NodeStatus(GoString name, struct node_status* result);

/* Return type for NodesStatus */
struct NodesStatus_return {
	int r0;
	char* r1;
};
extern struct NodesStatus_return NodesStatus(GoSlice names, GoSlice results, int* size);

/* Return type for DistribFile */
struct DistribFile_return {
	int r0;
	char* r1;
};
extern struct DistribFile_return DistribFile(GoString localFileOrDir, GoString targetPath, GoSlice names, GoSlice results, int* size);

/* Return type for PullFile */
struct PullFile_return {
	int r0;
	char* r1;
};
extern struct PullFile_return PullFile(GoString pathtype, GoString node, GoString fileOrDir, GoString targetdir, struct execution_result* result);

/* Return type for Run */
struct Run_return {
	int r0;
	char* r1;
};
extern struct Run_return Run(GoString runtask, GoSlice names, GoUint8 needAlign, GoSlice results, int* size);

/* Return type for XRun */
struct XRun_return {
	int r0;
	char* r1;
};
extern struct XRun_return XRun(GoString runtask, GoSlice names, GoInt delay, GoUint8 needAlign, GoSlice results, int* size);

/* Return type for GroupGetAll */
struct GroupGetAll_return {
	int r0;
	char* r1;
};
extern struct GroupGetAll_return GroupGetAll(GoSlice results, int* size);

/* Return type for GroupGet */
struct GroupGet_return {
	int r0;
	char* r1;
};
extern struct GroupGet_return GroupGet(GoString name, GoSlice results, int* size);

/* Return type for GroupSet */
struct GroupSet_return {
	int r0;
	char* r1;
};
extern struct GroupSet_return GroupSet(GoString name, GoUint8 nocheck, GoSlice names);

/* Return type for GroupDel */
struct GroupDel_return {
	int r0;
	char* r1;
};
extern struct GroupDel_return GroupDel(GoString name);

/* Return type for Prune */
struct Prune_return {
	int r0;
	char* r1;
};
extern struct Prune_return Prune();

/* Return type for ScenarioInfo */
struct ScenarioInfo_return {
	char* r0;
	int r1;
	char* r2;
};
extern struct ScenarioInfo_return ScenarioInfo(GoString name);

/* Return type for ScenariosInfo */
struct ScenariosInfo_return {
	int r0;
	char* r1;
};
extern struct ScenariosInfo_return ScenariosInfo(GoSlice results, int* size);

/* Return type for ScenarioVersion */
struct ScenarioVersion_return {
	char* r0;
	int r1;
	char* r2;
};
extern struct ScenarioVersion_return ScenarioVersion(GoString name);

/* Return type for NodeAppInfo */
struct NodeAppInfo_return {
	char* r0;
	int r1;
	char* r2;
};
extern struct NodeAppInfo_return NodeAppInfo(GoString name, GoString app, GoString scenario);

/* Return type for NodeAppsInfo */
struct NodeAppsInfo_return {
	int r0;
	char* r1;
};
extern struct NodeAppsInfo_return NodeAppsInfo(GoString name, GoSlice results, int* size);

/* Return type for Apply */
struct Apply_return {
	int r0;
	char* r1;
};
extern struct Apply_return Apply(GoString deployment, GoString target, GoString message, GoInt timeout, GoSlice logs, int* size);

#ifdef __cplusplus
}
#endif
