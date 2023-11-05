package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/piaodazhu/Octopoda/octl/config"

	"github.com/tidwall/pretty"
)

var enableColor bool
var enablePrint bool
var enableSpinner bool

func EnablePrint() {
	enablePrint = true
}

func EnableColor() {
	enablePrint = true
	enableColor = true
}

func EnableSpinner() {
	enableSpinner = true
}

func IsSpinnerEnabled() bool {
	return enableSpinner
}

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

	enableColor = false
	enablePrint = false
	enableSpinner = false
}

func PrintJSON(message interface{}) {
	if !enablePrint {
		return
	}
	if config.GlobalConfig.OutputPretty {
		switch msg := message.(type) {
		case string:
			s := pretty.Pretty([]byte(msg))
			if enableColor {
				s = pretty.Color(s, nil)
			}
			fmt.Println(replacer.Replace(string(s)))
		case []byte:
			s := pretty.Pretty(msg)
			if enableColor {
				s = pretty.Color(s, nil)
			}
			fmt.Println(replacer.Replace(string(s)))
		default:
			raw, _ := json.Marshal(msg)
			s := pretty.Pretty(raw)
			if enableColor {
				s = pretty.Color(s, nil)
			}
			fmt.Println(replacer.Replace(string(s)))
		}
	} else {
		switch msg := message.(type) {
		case string:
			fmt.Print(msg)
		case []byte:
			fmt.Print(string(msg))
		default:
			raw, _ := json.Marshal(msg)
			fmt.Println(string(raw))
		}
	}
}

const (
	template_YELLOW = "\033[1;33m%s\033[0m"
	template_GREEN  = "\033[1;32m%s\033[0m"
	template_RED    = "\033[1;31m%s\033[0m"
)

func printFatal(message string) {
	msg := fmt.Sprintf("[FATAL] %s", message)
	if enableColor {
		msg = fmt.Sprintf(template_RED, msg)
	}
	fmt.Println(msg)
}

func printWarn(message string) {
	msg := fmt.Sprintf("[WARN] %s", message)
	if enableColor {
		msg = fmt.Sprintf(template_YELLOW, msg)
	}
	fmt.Println(msg)
}

func printInfo(message string) {
	msg := fmt.Sprintf("[INFO] %s", message)
	if enableColor {
		msg = fmt.Sprintf(template_GREEN, msg)
	}
	fmt.Println(msg)
}


func PrintFatalf(format string, args ...interface{}) {
	if !enablePrint {
		return
	}
	printFatal(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func PrintFatalln(args ...interface{}) {
	if !enablePrint {
		return
	}
	printFatal(fmt.Sprint(args...))
	os.Exit(1)
}

func PrintInfof(format string, args ...interface{}) {
	if !enablePrint {
		return
	}
	printInfo(fmt.Sprintf(format, args...))
}

func PrintInfoln(args ...interface{}) {
	if !enablePrint {
		return
	}
	printInfo(fmt.Sprint(args...))
}

func PrintWarningf(format string, args ...interface{}) {
	if !enablePrint {
		return
	}
	printWarn(fmt.Sprintf(format, args...))
}

func PrintWarningln(args ...interface{}) {
	if !enablePrint {
		return
	}
	printWarn(fmt.Sprint(args...))
}


func PrintPanicf(format string, args ...interface{}) {
	err := fmt.Errorf(format, args...)
	panic(err)
}

func PrintPanicln(args ...interface{}) {
	err := errors.New(fmt.Sprint(args...))
	panic(err)
}
