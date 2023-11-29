package lib

type (
	Codec[T any] interface {
		Encode(*T) ([]byte, error)
		Decode([]byte) (*T, error)
	}
)
