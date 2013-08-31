package msg

import (
	"fmt"
	"github.com/Blackrush/gofus/protocol/frontend/types"
	"io"
	"strings"
	"time"
)

type HelloConnect struct {
	Ticket string
}

func (msg *HelloConnect) Opcode() string { return "HC" }
func (msg *HelloConnect) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Ticket)
	return nil
}
func (msg *HelloConnect) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "HC%s", &msg.Ticket)
	return nil
}

type BadVersion struct {
	Required string
}

func (msg *BadVersion) Opcode() string { return "AlEv" }
func (msg *BadVersion) Serialize(out io.Writer) error {
	fmt.Print(msg.Required)
	return nil
}
func (msg *BadVersion) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "AlEv%s", &msg.Required)
	return nil
}

type LoginError struct{}

func (msg *LoginError) Opcode() string                 { return "AlEf" }
func (msg *LoginError) Serialize(out io.Writer) error  { return nil }
func (msg *LoginError) Deserialize(in io.Reader) error { return nil }

type BannedUser struct{}

func (msg *BannedUser) Opcode() string                 { return "AlEb" }
func (msg *BannedUser) Serialize(out io.Writer) error  { return nil }
func (msg *BannedUser) Deserialize(in io.Reader) error { return nil }

type QueueStatusRequest struct{}

func (msg *QueueStatusRequest) Opcode() string                 { return "Af" }
func (msg *QueueStatusRequest) Serialize(out io.Writer) error  { return nil }
func (msg *QueueStatusRequest) Deserialize(in io.Reader) error { return nil }

type SetNickname struct {
	Nickname string
}

func (msg *SetNickname) Opcode() string { return "Ad" }
func (msg *SetNickname) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Nickname)
	return nil
}
func (msg *SetNickname) Deserialize(in io.Reader) error { return nil }

type SetCommunity struct {
	CommunityId int
}

func (msg *SetCommunity) Opcode() string { return "Ac" }
func (msg *SetCommunity) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.CommunityId)
	return nil
}
func (msg *SetCommunity) Deserialize(in io.Reader) error { return nil }

type LoginSuccess struct {
	IsAdmin bool
}

func (msg *LoginSuccess) Opcode() string { return "AlK" }
func (msg *LoginSuccess) Serialize(out io.Writer) error {
	fmt.Fprint(out, btoi(msg.IsAdmin))
	return nil
}
func (msg *LoginSuccess) Deserialize(in io.Reader) error { return nil }

type SetSecretQuestion struct {
	SecretQuestion string
}

func (msg *SetSecretQuestion) Opcode() string { return "AQ" }
func (msg *SetSecretQuestion) Serialize(out io.Writer) error {
	fmt.Fprint(out, strings.Replace(msg.SecretQuestion, " ", "+", -1))
	return nil
}
func (msg *SetSecretQuestion) Deserialize(in io.Reader) error { return nil }

type SetRealmServers struct {
	Servers []*types.RealmServer
}

func (msg *SetRealmServers) Opcode() string { return "AH" }
func (msg *SetRealmServers) Serialize(out io.Writer) error {
	// golang sucks
	values := make([]interface{}, len(msg.Servers))
	for i, server := range msg.Servers {
		values[i] = server
	}

	fmt.Fprint(out, types.NewMultipleFormatter("|", values...))
	return nil
}
func (msg *SetRealmServers) Deserialize(in io.Reader) error { return nil }

type SetPlayersError struct{}

func (msg *SetPlayersError) Opcode() string                 { return "AxE" }
func (msg *SetPlayersError) Serialize(out io.Writer) error  { return nil }
func (msg *SetPlayersError) Deserialize(in io.Writer) error { return nil }

type SetRealmServerPlayers struct {
	SubscriptionEnd time.Time
	Players         []*types.RealmServerPlayers
}

func (msg *SetRealmServerPlayers) Opcode() string { return "AxK" }
func (msg *SetRealmServerPlayers) Serialize(out io.Writer) error {
	values := make([]interface{}, len(msg.Players))
	for i, players := range msg.Players {
		values[i] = players
	}

	sub := uint64(msg.SubscriptionEnd.Sub(time.Now()).Nanoseconds()) / 1e6

	fmt.Fprint(out, sub, "|", types.NewMultipleFormatter("|", values...))

	return nil
}
func (msg *SetRealmServerPlayers) Deserialize(in io.Reader) error { return nil }

type RealmServerSelectionRequest struct {
	ServerId int
}

func (msg *RealmServerSelectionRequest) Opcode() string { return "AX" }
func (msg *RealmServerSelectionRequest) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.ServerId)
	return nil
}
func (msg *RealmServerSelectionRequest) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "AX%d", &msg.ServerId)
	return nil
}

type RealmServerSelectionResponse struct {
	Address string
	Port    uint16
	Ticket  string
}

func (msg *RealmServerSelectionResponse) Opcode() string { return "AYK" }
func (msg *RealmServerSelectionResponse) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%s:%d;%s", msg.Address, msg.Port, msg.Ticket)
	return nil
}
func (msg *RealmServerSelectionResponse) Deserialize(in io.Reader) error { return nil }

type RealmServerSelectionError struct{}

func (msg *RealmServerSelectionError) Opcode() string                 { return "AXE" }
func (msg *RealmServerSelectionError) Serialize(out io.Writer) error  { return nil }
func (msg *RealmServerSelectionError) Deserialize(in io.Reader) error { return nil }
