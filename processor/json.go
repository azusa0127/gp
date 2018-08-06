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
	s              *bufio.Scanner
	w              io.Writer
	jsonpathEvalFn gval.Evaluable
}

// NewJSONProcessor creates and initializes an JSONProcessor
func NewJSONProcessor(src io.Reader, dst io.Writer, jsonpathQuery string) *JSONProcessor {
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
		s:              bufio.NewScanner(src),
		w:              dst,
		jsonpathEvalFn: evalFn,
	}
}

// Process processes the stream and returns possible fatal error
func (j *JSONProcessor) Process() error {
	var err error
	var v interface{}
	var buf []byte
	for j.s.Scan() {
		if err = json.Unmarshal(j.s.Bytes(), &v); err != nil {
			return err
		}
		if v, err = j.jsonpathEvalFn(context.Background(), v); err != nil {
			return err
		}
		if buf, err = jsonFormatter.Marshal(v); err != nil {
			return err
		}
		if _, err = j.w.Write(buf); err != nil {
			return err
		}
		if _, err = j.w.Write(LineBreakBytes); err != nil {
			return err
		}
	}
	if err = j.s.Err(); err != nil {
		return err
	}
	return nil
}
