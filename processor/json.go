package processor

import (
	"context"
	"encoding/json"
	"log"

	"github.com/TylerBrock/colorjson"

	"github.com/PaesslerAG/jsonpath"
)

// JSONProcessor is the Processor implementation for JSON input and output
type JSONProcessor struct {
	evalFn    func(v interface{}) (interface{}, error)
	marshalFn func(v interface{}) ([]byte, error)
}

// NewJSONProcessor creates and initializes an JSONProcessor
func NewJSONProcessor(jsonpathQuery string, compressMode bool) *JSONProcessor {
	return &JSONProcessor{
		marshalFn: func() func(v interface{}) ([]byte, error) {
			if compressMode {
				return json.Marshal
			}
			jsonFormatter := colorjson.NewFormatter()
			jsonFormatter.Indent = 2
			return jsonFormatter.Marshal
		}(),
		evalFn: func() func(v interface{}) (interface{}, error) {
			if jsonpathQuery != "" {
				f, err := jsonpath.New(jsonpathQuery)
				if err != nil {
					log.Fatal(err)
				}
				return func(v interface{}) (interface{}, error) {
					return f(context.Background(), v)
				}
			}
			return func(v interface{}) (interface{}, error) { return v, nil }
		}(),
	}
}

func (j *JSONProcessor) Unmarshal(s []byte) (v interface{}, err error) {
	if err = json.Unmarshal(s, &v); err != nil {
		return
	}
	return j.evalFn(v)
}

func (j *JSONProcessor) Marshal(v interface{}) ([]byte, error) {
	return j.marshalFn(v)
}
