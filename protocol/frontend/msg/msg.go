package msg

import (
	protocol "github.com/Blackrush/gofus/protocol/frontend"
	"strconv"
)

// Converts a boolean value to an integer (1 if true or 0 if false)
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Converts an integer value to a boolean (true if 1, false otherwise)
func itob(i int) bool {
	return i == 1
}

// Converts an integer represented as a string to a boolean (true if "1", false otherwise)
func aitob(str string) bool {
	return str == "1"
}

// Converts a string to an integer
func atoi(str string) int {
	if result, err := strconv.Atoi(str); err == nil {
		return result
	} else {
		panic(err.Error())
	}
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
	producers["Ai"] = func() protocol.MessageContainer { return new(SetIdentity) }
	producers["AL"] = func() protocol.MessageContainer { return new(PlayersReq) }
	producers["AP"] = func() protocol.MessageContainer { return new(RandNameReq) }
	producers["AA"] = func() protocol.MessageContainer { return new(CreatePlayerReq) }
	producers["AS"] = func() protocol.MessageContainer { return new(PlayerSelectionReq) }
}

func New(opcode string) (protocol.MessageContainer, bool) {
	if producer, ok := producers[opcode]; ok {
		return producer(), true
	}
	return nil, false
}
