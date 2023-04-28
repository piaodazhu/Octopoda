package config

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
)

type JSONER interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}
var Jsoner JSONER

type StdJsoner struct{}
func (StdJsoner) Marshal(v any) ([]byte, error) {return json.Marshal(v)}
func (StdJsoner) Unmarshal(data []byte, v any) error {return json.Unmarshal(data, v)}
func setStdJsoner() {
	Jsoner = StdJsoner{}
}

func setFastJsoner() {
	Jsoner = jsoniter.ConfigCompatibleWithStandardLibrary
}