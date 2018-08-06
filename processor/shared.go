package processor

import "github.com/TylerBrock/colorjson"

// Processor is the interface of stream processors
type Processor interface {
	// Process processes the stream and returns possible fatal error
	Process() error
}

var jsonFormatter = colorjson.NewFormatter()

// LineBreak is the line break symbol in byte
var LineBreak = byte('\n')

// LineBreakBytes is the line break symbol in byte array
var LineBreakBytes = []byte("\n")

func init() {
	jsonFormatter.Indent = 2
}
