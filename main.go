package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"os"

	"github.com/azusa0127/gsp/processor"
	"golang.org/x/crypto/ssh/terminal"
)

var jsonpathQuery = flag.String("q", "", "Query in jsonpath")
var jsonCompressMode = flag.Bool("c", false, "Compress Mode")

var jsonFlag = flag.Bool("json", true, "Flag to process JSON stream")
var base64encodeFlag = flag.Bool("base64e", false, "Flag to encode the result string with base64")
var base64decodeFlag = flag.Bool("base64d", false, "Flag to decode the result string with base64")

var inputProcessor = flag.String("i", "text", "InputProcessor [json|base64|text]")
var outputProcessor = flag.String("o", "text", "OutputProcessor [json|base64|text]")

func main() {
	flag.Parse()
	var src io.Reader
	args := flag.Args()
	switch len(args) {
	case 2:
		src = bytes.NewBufferString(args[0])
		jsonpathQuery = &args[1]
	case 1:
		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			src = bytes.NewBufferString(args[0])
		} else {
			src = os.Stdin
			jsonpathQuery = &args[0]
		}
	case 0:
		src = os.Stdin
	default:
		log.Fatalln("Invalid arguments")
	}

	var in, out processor.Processor
	switch {
	case *jsonFlag:
		in = processor.NewJSONProcessor(*jsonpathQuery, *jsonCompressMode)
		out = in
	case *base64decodeFlag:
		in = processor.NewBase64DecodeProcessor()
		out = processor.TextProcessorPTR
	case *base64encodeFlag:
		in = processor.NewBase64EncodeProcessor()
		out = processor.TextProcessorPTR
	default:
		switch *inputProcessor {
		case "json":
			in = processor.NewJSONProcessor(*jsonpathQuery, *jsonCompressMode)
		case "base64e":
			in = processor.NewBase64EncodeProcessor()
		case "base64d":
			in = processor.NewBase64DecodeProcessor()
		case "text":
			in = processor.TextProcessorPTR
		default:
			log.Fatalln("Invalid input processor - " + *inputProcessor)
		}

		switch *outputProcessor {
		case "json":
			out = processor.NewJSONProcessor(*jsonpathQuery, *jsonCompressMode)
		case "base64e":
			out = processor.NewBase64EncodeProcessor()
		case "base64d":
			out = processor.NewBase64DecodeProcessor()
		case "text":
			out = processor.TextProcessorPTR
		default:
			log.Fatalln("Invalid output processor - " + *outputProcessor)
		}
	}

	var bufChan = make(chan interface{})
	var dst = os.Stdout
	go func() {
		if err := in.ParseStream(bufio.NewScanner(src), bufChan); err != nil {
			log.Fatalln(err)
		}
	}()
	if err := out.PushStream(bufChan, dst); err != nil {
		log.Fatalln(err)
	}
}
