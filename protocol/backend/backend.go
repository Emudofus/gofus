package backend

import (
	"github.com/Blackrush/gofus/protocol"
)

type Opcoder interface {
	Opcode() uint16
}

type Message interface {
	protocol.Serializer
	protocol.Deserializer
	Opcoder
}

type Sender interface {
	Send(msg Message) (int, error)
}

type CloseWither interface {
	CloseWith(msg Message) error
}
