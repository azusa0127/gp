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

var queryString = flag.String("q", "", "Query string to be passed in specified query engine")
var filterString = flag.String("filter", "", "Filter query to be used")
var jsonCompressMode = flag.Bool("c", false, "Compress Mode")
var noColorMode = flag.Bool("nc", false, "Flag to prettify output without color")

var jsonFlag = flag.Bool("json", false, "Flag to process JSON stream")
var base64encodeFlag = flag.Bool("base64e", false, "Flag to encode the result string with base64")
var base64decodeFlag = flag.Bool("base64d", false, "Flag to decode the result string with base64")
var toYAMLFlag = flag.Bool("toyaml", false, "Flag to convert JSON steam to yaml")

var inputProcessor = flag.String("i", "json", "InputProcessor [json]")
var queryEngine = flag.String("qe", "jmespath", "Query engine selection [jmespath|jsonpath] default `jmespath`")
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
	case *base64decodeFlag:
		p = processor.NewBase64DecodeProcessor()
	case *base64encodeFlag:
		p = processor.NewBase64EncodeProcessor()
	default:
		switch {
		case *jsonFlag:
			*inputProcessor = "json"
			*queryEngine = "jmespath"
			*outputProcessor = "json"
		case *toYAMLFlag:
			*inputProcessor = "json"
			*queryEngine = "jmespath"
			*outputProcessor = "yaml"
		}

		var in processor.UnmarshalFunction
		var filter processor.QueryEvalFunction
		var eval processor.QueryEvalFunction
		var out processor.MarshalFunction

		switch *inputProcessor {
		case "json":
			in = processor.NewJSONUnmarshalFunction()
		default:
			log.Fatalln("invalid input processor - " + *inputProcessor)
		}

		eval = func(v interface{}) (interface{}, error) { return v, nil }
		if *queryString != "" {
			switch *queryEngine {
			case "jmespath":
				eval = processor.NewJMESPathEvalFunction(*queryString)
			case "jsonpath":
				eval = processor.NewJSONPathEvalFunction(*queryString)
			}
		}

		if *filterString != "" {
			switch *queryEngine {
			case "jmespath":
				filter = processor.NewJMESPathEvalFunction(*filterString)
			case "jsonpath":
				filter = processor.NewJSONPathEvalFunction(*filterString)
			}
		}

		switch *outputProcessor {
		case "json":
			out = processor.NewJSONMarshalFunction(*jsonCompressMode, *noColorMode)
		case "yaml":
			out = processor.NewYAMLMarshalFunction()
		default:
			log.Fatalln("invalid output processor - " + *outputProcessor)
		}

		p = processor.NewMixedProcessor(in, filter, eval, out)
	}

	var dst = os.Stdout
	var buf []byte
	s := bufio.NewScanner(src)
	for s.Scan() {
		if buf, err = p.Process(buf, s.Bytes()); err != nil {
			log.Fatalln(err)
		}
		if buf != nil {
			if _, err = dst.Write(buf); err != nil {
				log.Fatalln(err)
			}
			if _, err = dst.Write(lineBreakBytes); err != nil {
				log.Fatalln(err)
			}
		}
	}
	if err = s.Err(); err != nil {
		log.Fatalln(err)
	}
}
