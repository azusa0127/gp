package processor

import (
	"gopkg.in/yaml.v2"
)

func NewYAMLMarshalFunction() MarshalFunction {
	return yaml.Marshal
}
