package processor

// Processor is the interface of stream processors
type Processor interface {
	// Process processes a single line in bytes from the stream,
	// produces a result in bytes along with anything error from it.
	Process(buf, src []byte) ([]byte, error)
}

// ObjectProcessor is the interface of object processors that unmashals bytes into object
// and marshal it back into bytes.
type ObjectProcessor interface {
	Unmarshal(s []byte) (interface{}, error)
	Marshal(v interface{}) ([]byte, error)
}

// MixedProcessor combines and converts 2 ObjectProcessor into a Processor
type MixedProcessor struct {
	Marshal   func(v interface{}) ([]byte, error)
	Unmarshal func(s []byte) (interface{}, error)
}

// Process processes a single line in bytes from the stream,
// produces a result in bytes along with anything error from it.
func (m *MixedProcessor) Process(buf, src []byte) ([]byte, error) {
	v, err := m.Unmarshal(src)
	if err != nil {
		return nil, err
	}
	return m.Marshal(v)
}

// NewMixedProcessor returns a MixedProcessor
func NewMixedProcessor(in, out ObjectProcessor) *MixedProcessor {
	return &MixedProcessor{
		Marshal:   out.Marshal,
		Unmarshal: in.Unmarshal,
	}
}
