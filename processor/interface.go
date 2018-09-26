package processor

// Processor is the interface of stream processors
type Processor interface {
	// Process processes a single line in bytes from the stream,
	// produces a result in bytes along with anything error from it.
	Process(buf, src []byte) ([]byte, error)
}

type UnmarshalFunction func(s []byte, v interface{}) error
type QueryEvalFunction func(v interface{}) (interface{}, error)
type MarshalFunction func(v interface{}) ([]byte, error)

type MixedProcessor struct {
	unmarshal UnmarshalFunction
	filter    QueryEvalFunction
	queryEval QueryEvalFunction
	marshal   MarshalFunction
}

// NewMixedProcessor returns a MixedProcessor
func NewMixedProcessor(unmarshal UnmarshalFunction, filter, queryEval QueryEvalFunction, marshal MarshalFunction) *MixedProcessor {
	return &MixedProcessor{
		unmarshal: unmarshal,
		filter:    filter,
		queryEval: queryEval,
		marshal:   marshal,
	}
}

func (m *MixedProcessor) Process(_, src []byte) ([]byte, error) {
	var v interface{}
	err := m.unmarshal(src, &v)
	if err != nil {
		return nil, err
	}
	if m.filter != nil {
		if q, err := m.filter(v); err != nil || q != true {
			return nil, err
		}
	}
	if v, err = m.queryEval(v); err != nil {
		return nil, err
	}
	return m.marshal(v)
}
