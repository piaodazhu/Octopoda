package buildinfo

import "fmt"

type BuildInfo struct {
	BuildVersion string
	BuildTime    string
	BuildName    string
	CommitID     string
}

var binfo BuildInfo 

func SetBuildInfo(buildinfo BuildInfo) {
	binfo = buildinfo
	if len(buildinfo.CommitID) > 8 {
		binfo.CommitID = buildinfo.CommitID[:8]
	}
}

func String() string {
	return fmt.Sprintf("%s: %s, time=%s, commitID=%s", binfo.BuildName, binfo.BuildVersion, binfo.BuildTime, binfo.CommitID)
}