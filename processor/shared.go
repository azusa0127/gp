package processor

import (
	"bufio"
	"io"
)

// Processor is the interface of stream processors
type Processor interface {
	// ParseStream parses source stream and sent the parsed results to PushStream argument channel
	ParseStream(src *bufio.Scanner, c chan<- interface{}) error
	// PushStream takes the parsed results from ParseStream and push the output conversion out
	PushStream(c <-chan interface{}, dst io.Writer) error
}

// LineBreakBytes is the line break symbol in byte array
var LineBreakBytes = []byte("\n")

var TextProcessorPTR = &PlainTextProcessor{}

type PlainTextProcessor struct{}

func (p *PlainTextProcessor) ParseStream(s *bufio.Scanner, c chan<- interface{}) error {
	defer close(c)
	for s.Scan() {
		c <- s.Bytes()
	}
	return s.Err()
}

func (p *PlainTextProcessor) PushStream(c <-chan interface{}, w io.Writer) (err error) {
	for buf := range c {
		if _, err = w.Write(buf.([]byte)); err != nil {
			return
		}
		if _, err = w.Write(LineBreakBytes); err != nil {
			return
		}
	}
	return nil
}
