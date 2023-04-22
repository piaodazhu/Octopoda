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
			PrintFatal("unsupported message type")
		}
	} else {
		switch msg := message.(type) {
		case string:
			fmt.Println(msg)
		case []byte:
			fmt.Println(string(msg))
		default:
			PrintFatal("unsupported message type")
		}
	}
}

func PrintFatal(message string) {
	fmt.Printf("\033[1;31;40mFatal Error: %s\033[0m\n", message)
	os.Exit(1)
}
