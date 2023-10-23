package output

import (
	"fmt"
	"octl/config"
	"os"
	"strings"

	"github.com/tidwall/pretty"
)

var replacer *strings.Replacer

func init() {
	replacer = strings.NewReplacer(
		"\\n", "\n",
		"\\t", "\t",
		"\\r", "\r",
		"\\u003c", "\u003c", // <
		"\\u003e", "\u003e", // >
		"\\u0026", "\u0026", // &
		"\\u001b", "\u001b", // esc
		`\"`, `"`,
		`\'`, `'`,
	)
}

func PrintJSON(message interface{}) {
	if config.GlobalConfig.OutputPretty {
		switch msg := message.(type) {
		case string:
			s := pretty.Color(pretty.Pretty([]byte(msg)), nil)
			fmt.Println(replacer.Replace(string(s)))
		case []byte:
			s := pretty.Color(pretty.Pretty(msg), nil)
			fmt.Println(replacer.Replace(string(s)))
		default:
			PrintFatalln("unsupported message type")
		}
	} else {
		switch msg := message.(type) {
		case string:
			fmt.Print(msg)
		case []byte:
			fmt.Print(string(msg))
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

func PrintPanicf(format string, args ...interface{}) {
	panic(fmt.Errorf("\033[1;31m%s\033[0m", fmt.Sprintf(format, args...)))
}

func PrintPanicln(args ...interface{}) {
	panic(fmt.Errorf("\033[1;31m%s\033[0m", fmt.Sprint(args...)))
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
