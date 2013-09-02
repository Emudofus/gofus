package msg

import (
	protocol "github.com/Blackrush/gofus/protocol/frontend"
)

// Converts a boolean value to an integer (1 if true or 0 if false)
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

var producers = make(map[string]func() protocol.MessageContainer)

func init() {
	// adds only that are received
	producers["Af"] = func() protocol.MessageContainer { return new(QueueStatusRequest) }
	producers["AX"] = func() protocol.MessageContainer { return new(RealmServerSelectionRequest) }
	producers["AT"] = func() protocol.MessageContainer { return new(RealmLoginReq) }
	producers["Ak"] = func() protocol.MessageContainer { return new(ClientUseKeyReq) }
	producers["AV"] = func() protocol.MessageContainer { return new(RegionalVersionReq) }
	producers["Ag"] = func() protocol.MessageContainer { return new(PlayersGiftsReq) }
	producers["AG"] = func() protocol.MessageContainer { return new(SetPlayerGiftReq) }
	producers["AL"] = func() protocol.MessageContainer { return new(PlayersReq) }
}

func New(opcode string) (protocol.MessageContainer, bool) {
	if producer, ok := producers[opcode]; ok {
		return producer(), true
	}
	return nil, false
}
