package processor

import (
	"encoding/json"

	"github.com/hokaccha/go-prettyjson"
)

func NewJSONUnmarshalFunction() UnmarshalFunction {
	return json.Unmarshal
}

func NewJSONMarshalFunction(compressMode, noColorMode bool) MarshalFunction {
	switch {
	case compressMode:
		return json.Marshal
	case noColorMode:
		return func(v interface{}) ([]byte, error) { return json.MarshalIndent(v, "", "  ") }
	default:
		return prettyjson.NewFormatter().Marshal
	}
}
