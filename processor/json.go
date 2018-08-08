package processor

import (
	"encoding/json"

	"github.com/TylerBrock/colorjson"
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
		jsonFormatter := colorjson.NewFormatter()
		jsonFormatter.Indent = 2
		return jsonFormatter.Marshal
	}
}
