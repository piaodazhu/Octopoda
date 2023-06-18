package output

import (
	"fmt"
	"octl/config"

	"os"

	"github.com/hokaccha/go-prettyjson"
)

func PrintJSON(message interface{}) {
	if config.GlobalConfig.OutputPretty {
		switch msg := message.(type) {
		case string:
			s, _ := prettyjson.Format([]byte(msg))
			fmt.Println(string(s))
		case []byte:
			s, _ := prettyjson.Format(msg)
			fmt.Println(string(s))
		default:
			PrintFatalln("unsupported message type")
		}
	} else {
		switch msg := message.(type) {
		case string:
			fmt.Println(msg)
		case []byte:
			fmt.Println(string(msg))
		default:
			PrintFatalln("unsupported message type")
		}
	}
}

func PrintFatalf(format string, args ...interface{}) {
	fmt.Printf("\033[1;31;40mFatal Error: %s\033[0m\n", fmt.Sprintf(format, args...))
	os.Exit(1)
}

func PrintFatalln(args ...interface{}) {
	fmt.Printf("\033[1;31;40mFatal Error: %s\033[0m\n", fmt.Sprint(args...))
	os.Exit(1)
}

func PrintInfof(format string, args ...interface{}) {
	fmt.Printf("\033[1;32;40mInfo: %s\033[0m\n", fmt.Sprintf(format, args...))
}

func PrintInfoln(args ...interface{}) {
	fmt.Printf("\033[1;32;40mInfo: %s\033[0m\n", fmt.Sprint(args...))
}

func PrintWarningf(format string, args ...interface{}) {
	fmt.Printf("\033[1;33;40mWarning: %s\033[0m\n", fmt.Sprintf(format, args...))
}

func PrintWarningln(args ...interface{}) {
	fmt.Printf("\033[1;33;40mWarning: %s\033[0m\n", fmt.Sprint(args...))
}
