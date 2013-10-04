package protocol

import (
	"github.com/Blackrush/gofus/protocol"
)

type MessageContainer interface {
	protocol.Serializer
	protocol.Deserializer
	Opcode() string
}

type Sender interface {
	Send(msg MessageContainer) (int, error)
}

type CloseWither interface {
	CloseWith(msg MessageContainer) error
}
