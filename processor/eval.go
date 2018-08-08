package processor

import (
	"context"
	"strconv"

	"github.com/PaesslerAG/jsonpath"
	"github.com/jmespath/go-jmespath"
)

func NewJMESPathEvalFunction(query string) QueryEvalFunction {
	return jmespath.MustCompile(query).Search
}

func NewJSONPathEvalFunction(query string) QueryEvalFunction {
	f, err := jsonpath.New(query)
	if err != nil {
		panic(`jsonpath: Compile(` + strconv.Quote(query) + `): ` + err.Error())
	}
	return func(v interface{}) (interface{}, error) {
		return f(context.Background(), v)
	}
}
