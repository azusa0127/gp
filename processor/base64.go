package processor

import (
	"bufio"
	"encoding/base64"
	"io"
)

// Base64Processor is the Processor implementation for Base64 encoding
type Base64Processor struct {
	s    *bufio.Scanner
	w    io.Writer
	proc func(src []byte) ([]byte, error)
}

// NewBase64Processor creates and initializes a Base64Processor
func NewBase64Processor(src io.Reader, dst io.Writer, encode bool) *Base64Processor {
	var proc func(src []byte) ([]byte, error)
	if encode {
		proc = func(src []byte) ([]byte, error) {
			buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
			base64.StdEncoding.Encode(buf, src)
			return buf, nil
		}
	} else {
		proc = func(src []byte) ([]byte, error) {
			buf := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
			_, err := base64.StdEncoding.Decode(buf, src)
			return buf, err
		}
	}
	return &Base64Processor{
		s:    bufio.NewScanner(src),
		w:    dst,
		proc: proc,
	}
}

// Process processes the stream and returns possible fatal error
func (b *Base64Processor) Process() error {
	var err error
	for b.s.Scan() {
		var buf []byte
		if buf, err = b.proc(b.s.Bytes()); err != nil {
			return err
		}
		if _, err = b.w.Write(buf); err != nil {
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
