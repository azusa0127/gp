package pipeline

import "errors"

var (
	// ErrNilArgV is the error returned when a nil buffer pointer v is passed in to Unmarshal()
	ErrNilArgV = errors.New("Value of buffer pointer argument v for Unmarshal() cannot be nil")

	// ErrInvalidTypeArgV is the error returned when the value type of argument v is unsupported
	ErrInvalidTypeArgV = errors.New("Value of argument v for Marshal()/Unmarshal() is passed in a type unsupported")
)

// Unmarshaller unmarshals a serialized data into some interface.
type Unmarshaller interface {
	Unmarshal(s []byte, v interface{}) error
}

// Marshaller serializes data byte array.
type Marshaller interface {
	Marshal(v interface{}) ([]byte, error)
}

// Transformer operates on data and changing it along the way.
type Transformer interface {
	Transform(v interface{}) interface{}
}

var (
	// ErrNullUnmarshallerInvalidTypeArgV is the error returned when the value type of argument v is unsupported
	ErrNullUnmarshallerInvalidTypeArgV = errors.New("NullUnmarshaller.Unmarshal() accepts buffer pointer v only as a non nil pointer in type *[]byte")

	// ErrNullMarshallerInvalidTypeArgV is the error returned when the value type of argument v is unsupported
	ErrNullMarshallerInvalidTypeArgV = errors.New("NullMarshaller.Marshal() accepts argument v only in type []byte")
)

type NullUnmarshaller struct{}

func (n *NullUnmarshaller) Unmarshal(s []byte, v interface{}) error {
	if v == nil {
		return ErrNilArgV
	}
	ptr, ok := v.(*[]byte)
	if !ok {
		return ErrNullUnmarshallerInvalidTypeArgV
	}
	*ptr = s
	return nil
}

type NullMarshaller struct{}

func (n *NullMarshaller) Marshal(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, ErrNilArgV
	}
	b, ok := v.([]byte)
	if !ok {
		return nil, ErrNullMarshallerInvalidTypeArgV
	}
	return b, nil
}

type NullTransformer struct{}

func (n *NullTransformer) Transform(v interface{}) interface{} {
	return v
}
