package protocols

type CommandParams struct {
	Command    string
	Background bool
	DelayTime  int
	ExecTs     int64
	NeedAlign  bool
}


type ScriptParams struct {
	FileName   string
	TargetPath string
	FileBuf    string
	DelayTime  int
	ExecTs     int64
	NeedAlign  bool
}

