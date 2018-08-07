package processor

import (
	"bufio"
	"encoding/base64"
	"io"
)

// Base64Processor is the Processor implementation for Base64 encoding
type Base64Processor struct {
	proc func(buf, src []byte) ([]byte, error)
}

func encodeBase64(buf, src []byte) ([]byte, error) {
	n := base64.StdEncoding.EncodedLen(len(src))
	if cap(buf) < n {
		buf = make([]byte, n)
	}
	base64.StdEncoding.Encode(buf, src)
	return buf[:n], nil
}

func decodeBase64(buf, src []byte) ([]byte, error) {
	n := base64.StdEncoding.DecodedLen(len(src))
	if cap(buf) < n {
		buf = make([]byte, n)
	}
	var err error
	n, err = base64.StdEncoding.Decode(buf, src)
	return buf[:n], err
}

// NewBase64EncodeProcessor creates and initializes a Base64Processor for encoding
func NewBase64EncodeProcessor() *Base64Processor {
	return &Base64Processor{proc: encodeBase64}
}

// NewBase64DecodeProcessor creates and initializes a Base64Processor for decoding
func NewBase64DecodeProcessor() *Base64Processor {
	return &Base64Processor{proc: decodeBase64}
}

func (b *Base64Processor) ParseStream(s *bufio.Scanner, c chan<- interface{}) (err error) {
	defer close(c)
	var buf []byte
	for s.Scan() {
		if buf, err = b.proc(buf, s.Bytes()); err != nil {
			return
		}
		c <- buf
	}
	return s.Err()
}

func (b *Base64Processor) PushStream(c <-chan interface{}, w io.Writer) (err error) {
	var buf []byte
	for s := range c {
		if buf, err = b.proc(buf, s.([]byte)); err != nil {
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
