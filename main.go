package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/azusa0127/jg/processor"
	"golang.org/x/crypto/ssh/terminal"
)

var jsonpathQuery = flag.String("q", "", "Query in jsonpath")

var base64encodeFlag = flag.Bool("base64e", false, "Flag to encode the result string with base64")
var base64decodeFlag = flag.Bool("base64d", false, "Flag to decode the result string with base64")

func base64helper(v interface{}, w io.Writer, encode bool) error {
	var err error
	var buf *bytes.Buffer
	switch v.(type) {
	case string:
		buf = bytes.NewBufferString(v.(string))
	case []byte:
		buf = bytes.NewBuffer(v.([]byte))
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		buf = bytes.NewBuffer(b)
	}

	if encode {
		_, err = io.Copy(base64.NewEncoder(base64.StdEncoding, w), buf)
	} else {
		_, err = io.Copy(w, base64.NewDecoder(base64.StdEncoding, buf))
	}
	return err
}

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
		p = processor.NewBase64Processor(src, dst, true)
	case *base64decodeFlag:
		p = processor.NewBase64Processor(src, dst, false)
	default:
		p = processor.NewJSONProcessor(src, dst, *jsonpathQuery)
	}

	if err = p.Process(); err != nil {
		log.Fatalln(err)
	}
}
