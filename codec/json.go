package codec

import (
	"encoding/json"

	"github.com/Meduzz/abstractions.go/lib"
)

type (
	jsonCodec[T any] struct{}
)

func NewJsonCodec[T any]() lib.Codec[T] {
	return &jsonCodec[T]{}
}

func (j *jsonCodec[T]) Encode(it *T) ([]byte, error) {
	return json.Marshal(it)
}

func (j *jsonCodec[T]) Decode(data []byte) (*T, error) {
	it := new(T)
	err := json.Unmarshal(data, it)

	return it, err
}
