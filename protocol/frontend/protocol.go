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

type MessageContainer interface {
	Serializer
	Deserializer
	Opcode() string
}

type Sender interface {
	Send(msg MessageContainer) (int, error)
}

type CloseWither interface {
	CloseWith(msg MessageContainer) error
}
