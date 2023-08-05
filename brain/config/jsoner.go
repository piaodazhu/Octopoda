package config

import (
	"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type JSONER interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

var Jsoner JSONER

type StdJsoner struct{}

func (StdJsoner) Marshal(v any) ([]byte, error)      { return json.Marshal(v) }
func (StdJsoner) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }
func setStdJsoner() {
	Jsoner = StdJsoner{}
	fmt.Println(Jsoner)
}

func setFastJsoner() {
	Jsoner = jsoniter.ConfigFastest
	fmt.Println(Jsoner)
}
