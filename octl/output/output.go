package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/protocols"

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

func printBytes(b []byte, isPretty bool) {
	if isPretty {
		b = pretty.Pretty(b)
	}
	if enableColor {
		b = pretty.Color(b, nil)
	}
	str := string(b)
	if isPretty {
		str = replacer.Replace(str)
	}
	fmt.Println(str)
}

func PrintString(info string) {
	if !enablePrint {
		return
	}
	fmt.Println(info)
}

func PrintJSON(message interface{}) {
	if !enablePrint {
		return
	}

	switch msg := message.(type) {
	case string:
		printBytes([]byte(msg), config.GlobalConfig.OutputPretty)
	case []byte:
		printBytes(msg, config.GlobalConfig.OutputPretty)
	case protocols.ExecutionResults:
		results := msg.ToText()
		for i := range results.Results {
			if results.Results[i].Code == 0 {
				results.Results[i].Result = "\n" + results.Results[i].Result
			}
		}
		raw, _ := json.Marshal(results)
		printBytes(raw, config.GlobalConfig.OutputPretty)
	default:
		raw, _ := json.Marshal(msg)
		printBytes(raw, config.GlobalConfig.OutputPretty)
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
