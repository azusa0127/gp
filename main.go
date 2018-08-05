package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
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

var bash64EncodeFlag = flag.Bool("base64e", false, "Flag to encode the result string with base64")
var bash64DecodeFlag = flag.Bool("base64d", false, "Flag to decode the result string with base64")

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

func process(v interface{}) (interface{}, error) {
	var err error
	v, err = jsonpathEvalFn(context.Background(), v)
	if err != nil {
		return v, err
	}

	return v, nil
}

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
	jsonEnc := json.NewEncoder(os.Stdout)
	jsonEnc.SetIndent("", "  ")

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

		switch {
		case *bash64EncodeFlag:
			err = base64helper(buf, os.Stdout, true)
		case *bash64DecodeFlag:
			err = base64helper(buf, os.Stdout, false)
		default:
			err = jsonEnc.Encode(buf)
		}
		if err != nil {
			log.Fatalln(err)
		}
	}
	if err = scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
