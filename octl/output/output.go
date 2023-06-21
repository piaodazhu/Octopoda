package output

import (
	"fmt"
	"octl/config"

	"os"

	"github.com/tidwall/pretty"
)

func PrintJSON(message interface{}) {
	if config.GlobalConfig.OutputPretty {
		switch msg := message.(type) {
		case string:
			s := pretty.Color(pretty.Pretty([]byte(msg)), nil)
			fmt.Println(string(s))
		case []byte:
			s := pretty.Color(pretty.Pretty(msg), nil)
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
	fmt.Printf("\033[1;31mFatal Error: %s\033[0m\n", fmt.Sprintf(format, args...))
	os.Exit(1)
}

func PrintFatalln(args ...interface{}) {
	fmt.Printf("\033[1;31mFatal Error: %s\033[0m\n", fmt.Sprint(args...))
	os.Exit(1)
}

func PrintInfof(format string, args ...interface{}) {
	fmt.Printf("\033[1;32mInfo: %s\033[0m\n", fmt.Sprintf(format, args...))
}

func PrintInfoln(args ...interface{}) {
	fmt.Printf("\033[1;32mInfo: %s\033[0m\n", fmt.Sprint(args...))
}

func PrintWarningf(format string, args ...interface{}) {
	fmt.Printf("\033[1;33mWarning: %s\033[0m\n", fmt.Sprintf(format, args...))
}

func PrintWarningln(args ...interface{}) {
	fmt.Printf("\033[1;33mWarning: %s\033[0m\n", fmt.Sprint(args...))
}
