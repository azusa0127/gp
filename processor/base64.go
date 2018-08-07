package processor

import (
	"encoding/base64"
)

// Base64Processor is the Processor implementation for Base64 encoding
type Base64Processor struct {
	process func(buf, src []byte) ([]byte, error)
}

// NewBase64EncodeProcessor creates and initializes a Base64Processor for encoding
func NewBase64EncodeProcessor() *Base64Processor {
	return &Base64Processor{
		process: func(buf, src []byte) ([]byte, error) {
			n := base64.StdEncoding.EncodedLen(len(src))
			if cap(buf) < n {
				buf = make([]byte, n)
			}
			base64.StdEncoding.Encode(buf, src)
			return buf[:n], nil
		},
	}
}

// NewBase64DecodeProcessor creates and initializes a Base64Processor for decoding
func NewBase64DecodeProcessor() *Base64Processor {
	return &Base64Processor{
		process: func(buf, src []byte) ([]byte, error) {
			n := base64.StdEncoding.DecodedLen(len(src))
			if cap(buf) < n {
				buf = make([]byte, n)
			}
			var err error
			n, err = base64.StdEncoding.Decode(buf, src)
			return buf[:n], err
		},
	}
}

func (b *Base64Processor) Process(buf, src []byte) ([]byte, error) {
	return b.process(buf, src)
}
