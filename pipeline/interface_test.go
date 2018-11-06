package pipeline

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNullUnmarshaller_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr bool
	}{
		{
			"Unmarshals valid bytes correctly",
			[]byte("abcde"),
			[]byte("abcde"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NullUnmarshaller{}
			var got []byte
			if err := n.Unmarshal(tt.input, &got); (err != nil) != tt.wantErr {
				t.Errorf("NullUnmarshaller.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if bytes.Compare(got, tt.want) != 0 {
				t.Errorf("NullUnmarshaller.Unmarshal() result mismatch: Actual = %s, Wanted %s", got, tt.want)
			}
		})
	}
}

func TestNullMarshaller_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    []byte
		wantErr bool
	}{
		{
			"Marshal() should return the exact input back if in []byte",
			[]byte("abcde"),
			[]byte("abcde"),
			false,
		},
		{
			"Marshal() should return the error if input in type other than []byte",
			"abcde",
			nil,
			true,
		},
		{
			"Marshal() should return nil if input in is nil with ErrNilArgV",
			nil,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NullMarshaller{}
			got, err := n.Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NullMarshaller.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NullMarshaller.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullTransformer_Transform(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
	}{
		{
			"Transform() should return the exact input back if in []byte",
			[]byte("abcde"),
		},
		{
			"Transform() should return the error if input in string",
			"abcde",
		},
		{
			"Transform() should return nil if input in is nil with no error",
			nil,
		},
		{
			"Transform() should return the exact input back with an input object",
			&NullMarshaller{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NullTransformer{}
			if got := n.Transform(tt.v); !reflect.DeepEqual(got, tt.v) {
				t.Errorf("NullTransformer.Transform() = %v, want %v", got, tt.v)
			}
		})
	}
}

func TestNullDeserializer_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr bool
	}{
		{
			"Deserialize valid bytes correctly",
			[]byte("abcde"),
			[]byte("abcde"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &NullDeserializer{}
			got, err := n.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NullDeserializer.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NullDeserializer.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
