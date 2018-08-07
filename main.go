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
var base64encodeFlag = flag.Bool("base64e", false, "Flag to encode the result string with base64")
var base64decodeFlag = flag.Bool("base64d", false, "Flag to decode the result string with base64")

func main() {
	flag.Parse()
	var src io.Reader
	var dst = os.Stdout
	var err error
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

	var p processor.Processor

	switch {
	case *base64encodeFlag:
		p = processor.NewBase64EncodeProcessor()
	case *base64decodeFlag:
		p = processor.NewBase64DecodeProcessor()
	default:
		p = processor.NewJSONProcessor(*jsonpathQuery)
	}

	if err = p.Process(bufio.NewScanner(src), dst); err != nil {
		log.Fatalln(err)
	}
}
