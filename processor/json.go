package processor

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/PaesslerAG/jsonpath"
	"github.com/TylerBrock/colorjson"
	"github.com/jmespath/go-jmespath"
)

// JSONProcessor is the Processor implementation for JSON input and output
type JSONProcessor struct {
	evalFn    func(v interface{}) (interface{}, error)
	marshalFn func(v interface{}) ([]byte, error)
}

const (
	// JMESPathEngine uses JMESPath (http://jmespath.org/) to parse query
	JMESPathEngine = "jmespath"
	// JSONPathEngine uses JSONPath (http://goessner.net/articles/JsonPath/index.html) to parse query
	JSONPathEngine = "jsonpath"
)

// NewJSONProcessor creates and initializes an JSONProcessor
func NewJSONProcessor(queryEngine, query string, compressMode bool) *JSONProcessor {
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
			if query != "" {
				switch queryEngine {
				case JMESPathEngine:
					return jmespath.MustCompile(query).Search
				case JSONPathEngine:
					f, err := jsonpath.New(query)
					if err != nil {
						panic(`jsonpath: Compile(` + strconv.Quote(query) + `): ` + err.Error())
					}
					return func(v interface{}) (interface{}, error) {
						return f(context.Background(), v)
					}
				default:
					panic(`invalid queryEngine - ` + queryEngine)
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
