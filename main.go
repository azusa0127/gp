package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
	"golang.org/x/crypto/ssh/terminal"
)

var jsonpathQuery = flag.String("q", "", "Query in jsonpath")
var jsonpathEvalFn gval.Evaluable = func(c context.Context, v interface{}) (interface{}, error) { return v, nil }

func main() {
	flag.Parse()
	var src io.Reader
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

	if *jsonpathQuery != "" {
		jsonpathEvalFn, err = jsonpath.New(*jsonpathQuery)
		if err != nil {
			log.Fatalln(err)
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		var buf interface{}
		if err = json.Unmarshal(scanner.Bytes(), &buf); err != nil {
			log.Fatalln(err)
		}
		buf, err = jsonpathEvalFn(context.Background(), buf)
		if err != nil {
			log.Fatalln(err)
		}
		enc.Encode(buf)
	}
	if err = scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
