package processor

import (
	"bufio"
	"encoding/base64"
	"io"
)

// Base64encodeProcessor is the Processor implementation for Base64 encoding
type Base64encodeProcessor struct {
	s *bufio.Scanner
	w io.Writer
}

// NewBase64encodeProcessor creates and initializes a Base64encodeProcessor
func NewBase64encodeProcessor(src io.Reader, dst io.Writer) *Base64encodeProcessor {
	return &Base64encodeProcessor{
		s: bufio.NewScanner(src),
		w: dst,
	}
}

// Process processes the stream and returns possible fatal error
func (b *Base64encodeProcessor) Process() error {
	var err error
	for b.s.Scan() {
		enc := base64.NewEncoder(base64.StdEncoding, b.w)
		if _, err = enc.Write(b.s.Bytes()); err != nil {
			return err
		}
		if err = enc.Close(); err != nil {
			return err
		}
		if _, err = b.w.Write(LineBreakBytes); err != nil {
			return err
		}
	}
	if err = b.s.Err(); err != nil {
		return err
	}
	return nil
}
