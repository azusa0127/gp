package pipeline

type SimplePipeline struct {
	dese Deserializer
	tran Transformer
	mars Marshaller
}

func NewSimplePipeline(d Deserializer, t Transformer, m Marshaller) *SimplePipeline {
	if d == nil {
		d = NullDeserializer
	}
	if t == nil {
		t == NullTransformer
	}
	if m == nil {
		m = NullMarshaller
	}
	return &SimplePipeline{dese: d, tran: t, mars: m}
}

func (s *SimplePipeline) Process(s []byte) ([]byte, error) {
	raw, err := s.dese.Parse(s)
	if err != nil {
		return
	}
	trans, err := s.tran.Transform(raw)
	if err != nil {
		return
	}
	return s.mars(trans)
}
