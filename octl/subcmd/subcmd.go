package subcmd

import "octl/shell"

func Apply(arglist []string) {

}

func Get(arglist []string) {

}

func Status(arglist []string) {

}

func Log(arglist []string) {

}

func Version(arglist []string) {

}

func Reset(arglist []string) {

}

func Shell(arglist []string) {
	if len(arglist) != 1 {
		return
	}
	shell.SSH(arglist[0])
}

func Upload(arglist []string) {

}

func Spread(arglist []string) {

}

func Purge(arglist []string) {

}
