package processor

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/TylerBrock/colorjson"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
)

// JSONProcessor is the Processor implementation for JSON input and output
type JSONProcessor struct {
	jsonpathEvalFn gval.Evaluable
	marshalFn      func(v interface{}) ([]byte, error)
}

// NewJSONProcessor creates and initializes an JSONProcessor
func NewJSONProcessor(jsonpathQuery string, compressMode bool) *JSONProcessor {
	var evalFn gval.Evaluable
	if jsonpathQuery != "" {
		var err error
		evalFn, err = jsonpath.New(jsonpathQuery)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		evalFn = func(c context.Context, v interface{}) (interface{}, error) { return v, nil }
	}

	var marshalFn func(v interface{}) ([]byte, error)
	if compressMode {
		marshalFn = json.Marshal
	} else {
		jsonFormatter := colorjson.NewFormatter()
		jsonFormatter.Indent = 2
		marshalFn = jsonFormatter.Marshal
	}

	return &JSONProcessor{
		jsonpathEvalFn: evalFn,
		marshalFn:      marshalFn,
	}
}

func (j *JSONProcessor) ParseStream(s *bufio.Scanner, c chan<- interface{}) (err error) {
	defer close(c)
	var v interface{}
	for s.Scan() {
		if err = json.Unmarshal(s.Bytes(), v); err != nil {
			return
		}
		if v, err = j.jsonpathEvalFn(context.Background(), v); err != nil {
			return
		}
		c <- v
	}
	return s.Err()
}

func (j *JSONProcessor) PushStream(c <-chan interface{}, w io.Writer) (err error) {
	var buf []byte
	for v := range c {
		if _, err = j.marshalFn(v); err != nil {
			return
		}
		if _, err = w.Write(buf); err != nil {
			return
		}
		if _, err = w.Write(LineBreakBytes); err != nil {
			return
		}
	}
	return nil
}
