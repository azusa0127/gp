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
	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

var lineBreak = []byte("\n")

var jsonpathQuery = flag.String("q", "", "Query in jsonpath")
var jsonpathEvalFn gval.Evaluable = func(c context.Context, v interface{}) (interface{}, error) { return v, nil }
var jsonFormatter = &colorjson.Formatter{
	KeyColor:        color.New(color.FgWhite),
	StringColor:     color.New(color.FgGreen),
	BoolColor:       color.New(color.FgYellow),
	NumberColor:     color.New(color.FgCyan),
	NullColor:       color.New(color.FgMagenta),
	StringMaxLength: 0,
	DisabledColor:   false,
	Indent:          2,
	RawStrings:      false,
}

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

func jsonHelper(v interface{}, w io.Writer) error {
	b, err := jsonFormatter.Marshal(v)
	if err != nil {
		return err
	}
	w.Write(b)
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

	if *jsonpathQuery != "" {
		jsonpathEvalFn, err = jsonpath.New(*jsonpathQuery)
		if err != nil {
			log.Fatalln(err)
		}
	}

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
			err = base64helper(buf, dst, true)
		case *bash64DecodeFlag:
			err = base64helper(buf, dst, false)
		default:
			err = jsonHelper(buf, dst)
		}
		if err != nil {
			log.Fatalln(err)
		}
		dst.Write(lineBreak)
	}
	if err = scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
