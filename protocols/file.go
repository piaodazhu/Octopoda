package protocols

type FileParams struct {
	PackName    string
	TargetPath  string
	FileBuf     string
	ForceCreate bool
}

type FileDistribParams struct {
	LocalFile   string
	TargetPath  string
	TargetNodes []string
}

type FilePullParams struct {
	PackName   string
	PathType   string
	TargetPath string
	FileBuf    string
}

type FileSpreadParams struct {
	TargetPath  string
	FileOrDir   string
	TargetNodes []string
}

