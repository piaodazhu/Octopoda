package protocols

type WorkgroupInfo struct {
	Path     string
	Password string
}

type WorkgroupMembers []string
type WorkgroupChildren []string

type WorkgroupMembersPostParams struct {
	Path    string
	IsAdd   bool
	Members WorkgroupMembers
}
