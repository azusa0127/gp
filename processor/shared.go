package processor

import (
	"bufio"
	"io"

	"github.com/TylerBrock/colorjson"
)

// Processor is the interface of stream processors
type Processor interface {
	// Process processes the stream and returns possible fatal error
	Process(src *bufio.Scanner, dst io.Writer) error
}

var jsonFormatter = colorjson.NewFormatter()

// LineBreakBytes is the line break symbol in byte array
var LineBreakBytes = []byte("\n")

func init() {
	jsonFormatter.Indent = 2
}
