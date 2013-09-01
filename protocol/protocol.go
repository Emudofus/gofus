package protocol

import (
	"io"
)

type Serializer interface {
	Serialize(out io.Writer) error
}

type Deserializer interface {
	Deserialize(in io.Reader) error
}
