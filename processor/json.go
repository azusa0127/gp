package processor

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
)

// JSONProcessor is the Processor implementation for JSON input and output
type JSONProcessor struct {
	jsonpathEvalFn gval.Evaluable
}

// NewJSONProcessor creates and initializes an JSONProcessor
func NewJSONProcessor(jsonpathQuery string) *JSONProcessor {
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

	return &JSONProcessor{
		jsonpathEvalFn: evalFn,
	}
}

// Process processes the stream and returns possible fatal error
func (j *JSONProcessor) Process(s *bufio.Scanner, w io.Writer) error {
	var err error
	var v interface{}
	var buf []byte
	for s.Scan() {
		if err = json.Unmarshal(s.Bytes(), &v); err != nil {
			return err
		}
		if v, err = j.jsonpathEvalFn(context.Background(), v); err != nil {
			return err
		}
		if buf, err = jsonFormatter.Marshal(v); err != nil {
			return err
		}
		if _, err = w.Write(buf); err != nil {
			return err
		}
		if _, err = w.Write(LineBreakBytes); err != nil {
			return err
		}
	}
	return s.Err()
}
