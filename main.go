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

var queryEngine = flag.String("qe", processor.JMESPathEngine, "Query engine selection [jmespath|jsonpath] default `jmespath`")
var queryString = flag.String("q", "", "Query string to be passed in specified query engine")
var jsonCompressMode = flag.Bool("c", false, "Compress Mode")

var jsonFlag = flag.Bool("json", true, "Flag to process JSON stream")
var base64encodeFlag = flag.Bool("base64e", false, "Flag to encode the result string with base64")
var base64decodeFlag = flag.Bool("base64d", false, "Flag to decode the result string with base64")

var inputProcessor = flag.String("i", "json", "InputProcessor [json]")
var outputProcessor = flag.String("o", "json", "OutputProcessor [json]")

var lineBreakBytes = []byte("\n")

func main() {
	flag.Parse()
	var err error
	var src io.Reader
	args := flag.Args()
	switch len(args) {
	case 2:
		src = bytes.NewBufferString(args[0])
		queryString = &args[1]
	case 1:
		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			src = bytes.NewBufferString(args[0])
		} else {
			src = os.Stdin
			queryString = &args[0]
		}
	case 0:
		src = os.Stdin
	default:
		log.Fatalln("Invalid arguments")
	}

	var p processor.Processor
	switch {
	case *jsonFlag:
		jp := processor.NewJSONProcessor(*queryEngine, *queryString, *jsonCompressMode)
		p = processor.NewMixedProcessor(jp, jp)
	case *base64decodeFlag:
		p = processor.NewBase64DecodeProcessor()
	case *base64encodeFlag:
		p = processor.NewBase64EncodeProcessor()
	default:
		var in, out processor.ObjectProcessor
		switch *inputProcessor {
		case "json":
			in = processor.NewJSONProcessor(*queryEngine, *queryString, *jsonCompressMode)
		default:
			log.Fatalln("Invalid input processor - " + *inputProcessor)
		}

		switch *outputProcessor {
		case "json":
			out = processor.NewJSONProcessor(*queryEngine, *queryString, *jsonCompressMode)
		default:
			log.Fatalln("Invalid output processor - " + *outputProcessor)
		}
		p = processor.NewMixedProcessor(in, out)
	}

	var dst = os.Stdout
	var buf []byte
	s := bufio.NewScanner(src)
	for s.Scan() {
		if buf, err = p.Process(buf, s.Bytes()); err != nil {
			log.Fatalln(err)
		}
		if _, err = dst.Write(buf); err != nil {
			log.Fatalln(err)
		}
		if _, err = dst.Write(lineBreakBytes); err != nil {
			log.Fatalln(err)
		}
	}
	if err = s.Err(); err != nil {
		log.Fatalln(err)
	}
}
