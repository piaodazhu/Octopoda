package output

import (
	"fmt"
	"octl/config"

	"github.com/hokaccha/go-prettyjson"
)

func PrintJSON(message interface{}) {
	if config.GlobalConfig.Server.OutputPretty {
		switch msg := message.(type) {
		case string:
			s, _ := prettyjson.Format([]byte(msg))
			fmt.Println(string(s))
		case []byte:
			s, _ := prettyjson.Format(msg)
			fmt.Println(string(s))
		default:
			panic("unsupported message type")
		}
	} else {
		switch msg := message.(type) {
		case string:
			fmt.Println(msg)
		case []byte:
			fmt.Println(string(msg))
		default:
			panic("unsupported message type")
		}
	}
}
