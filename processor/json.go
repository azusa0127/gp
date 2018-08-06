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
	for j.s.Scan() {
		var buf interface{}
		if err = json.Unmarshal(j.s.Bytes(), &buf); err != nil {
			return err
		}
		buf, err = j.jsonpathEvalFn(context.Background(), buf)
		if err != nil {
			return err
		}
		b, err := jsonFormatter.Marshal(buf)
		if err != nil {
			return err
		}
		_, err = j.w.Write(b)
		if err != nil {
			return err
		}
		_, err = j.w.Write(LineBreakBytes)
		if err != nil {
			return err
		}
	}
	if err = j.s.Err(); err != nil {
		return err
	}
	return nil
}
