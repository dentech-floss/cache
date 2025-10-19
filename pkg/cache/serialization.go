package cache

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
)

// Serializer defines the interface for serializing and deserializing data.
type Serializer interface {
	// Serialize converts a value to bytes.
	Serialize(v interface{}) ([]byte, error)
	// Deserialize converts bytes back to a value.
	Deserialize(data []byte, v interface{}) error
}

// JSONSerializer implements JSON serialization.
type JSONSerializer struct{}

// NewJSONSerializer creates a new JSON serializer.
func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

// Serialize converts a value to JSON bytes.
func (j *JSONSerializer) Serialize(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Deserialize converts JSON bytes back to a value.
func (j *JSONSerializer) Deserialize(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// GobSerializer implements Go binary serialization.
type GobSerializer struct{}

// NewGobSerializer creates a new gob serializer.
func NewGobSerializer() *GobSerializer {
	return &GobSerializer{}
}

// Serialize converts a value to gob bytes.
func (g *GobSerializer) Serialize(v interface{}) ([]byte, error) {
	// We need to use a buffer to get the bytes
	var buf []byte
	err := gob.NewEncoder(&gobBuffer{&buf}).Encode(v)
	return buf, err
}

// Deserialize converts gob bytes back to a value.
func (g *GobSerializer) Deserialize(data []byte, v interface{}) error {
	return gob.NewDecoder(&gobBuffer{&data}).Decode(v)
}

// gobBuffer is a simple buffer implementation for gob encoding/decoding.
type gobBuffer struct {
	data *[]byte
}

func (b *gobBuffer) Write(p []byte) (n int, err error) {
	*b.data = append(*b.data, p...)
	return len(p), nil
}

func (b *gobBuffer) Read(p []byte) (n int, err error) {
	if len(*b.data) == 0 {
		return 0, io.EOF
	}
	n = copy(p, *b.data)
	*b.data = (*b.data)[n:]
	return n, nil
}

// NewSerializer creates a serializer based on the specified type.
func NewSerializer(serializationType SerializationType) (Serializer, error) {
	switch serializationType {
	case SerializationJSON:
		return &JSONSerializer{}, nil
	case SerializationGob:
		return &GobSerializer{}, nil
	case SerializationProtobuf:
		return nil, errors.New("protobuf serialization requires special handling - use NewDistributed")
	default:
		return nil, errors.New("unknown serialization type")
	}
}
