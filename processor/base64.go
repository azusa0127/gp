package processor

import (
	"bufio"
	"encoding/base64"
	"io"
)

// Base64encodeProcessor is the Processor implementation for Base64 encoding
type Base64encodeProcessor struct {
	s   *bufio.Scanner
	w   io.Writer
	enc io.WriteCloser
}

// NewBase64encodeProcessor creates and initializes a Base64encodeProcessor
func NewBase64encodeProcessor(src io.Reader, dst io.Writer) *Base64encodeProcessor {
	return &Base64encodeProcessor{
		s:   bufio.NewScanner(src),
		w:   dst,
		enc: base64.NewEncoder(base64.StdEncoding, dst),
	}
}

// Process processes the stream and returns possible fatal error
func (b *Base64encodeProcessor) Process() error {
	var err error
	for b.s.Scan() {
		_, err = b.enc.Write(b.s.Bytes())
		if err != nil {
			return err
		}
		_, err = b.w.Write(LineBreakBytes)
		if err != nil {
			return err
		}
	}
	if err = b.s.Err(); err != nil {
		return err
	}
	return b.enc.Close()
}
