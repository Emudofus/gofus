package msg

import (
	"fmt"
	"io"
)

type HelloGame struct{}

func (msg *HelloGame) Opcode() string                 { return "HG" }
func (msg *HelloGame) Serialize(out io.Writer) error  { return nil }
func (msg *HelloGame) Deserialize(in io.Reader) error { return nil }

type RealmLoginReq struct {
	Ticket string
}

func (msg *RealmLoginReq) Opcode() string { return "AT" }
func (msg *RealmLoginReq) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Ticket)
	return nil
}
func (msg *RealmLoginReq) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "AT%s", &msg.Ticket)
	return nil
}

type RealmLoginSuccess struct {
	CommunityId int
}

func (msg *RealmLoginSuccess) Opcode() string { return "ATK" }
func (msg *RealmLoginSuccess) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.CommunityId)
	return nil
}
func (msg *RealmLoginSuccess) Deserialize(in io.Reader) error { return nil }

type RealmLoginError struct{}

func (msg *RealmLoginError) Opcode() string                 { return "ATE" }
func (msg *RealmLoginError) Serialize(out io.Writer) error  { return nil }
func (msg *RealmLoginError) Deserialize(in io.Reader) error { return nil }

type ClientUseKeyReq struct {
	KeyId int
}

func (msg *ClientUseKeyReq) Opcode() string { return "Ak" }
func (msg *ClientUseKeyReq) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%x", msg.KeyId)
	return nil
}
func (msg *ClientUseKeyReq) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "Ak%x", &msg.KeyId)
	return nil
}

type RegionalVersionReq struct{}

func (msg *RegionalVersionReq) Opcode() string                 { return "AV" }
func (msg *RegionalVersionReq) Serialize(out io.Writer) error  { return nil }
func (msg *RegionalVersionReq) Deserialize(in io.Reader) error { return nil }

type RegionalVersionResp struct {
	CommunityId int
}

func (msg *RegionalVersionResp) Opcode() string { return "AV" }
func (msg *RegionalVersionResp) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.CommunityId)
	return nil
}
func (msg *RegionalVersionResp) Deserialize(in io.Reader) error { return nil }

type PlayersGiftsReq struct {
	Language string
}

func (msg *PlayersGiftsReq) Opcode() string { return "Ag" }
func (msg *PlayersGiftsReq) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Language)
	return nil
}
func (msg *PlayersGiftsReq) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "Ag%s", &msg.Language)
	return nil
}

type SetPlayerGiftReq struct {
	GiftId   int
	PlayerId int
}

func (msg *SetPlayerGiftReq) Opcode() string { return "AG" }
func (msg *SetPlayerGiftReq) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%d|%d", msg.GiftId, msg.PlayerId)
	return nil
}
func (msg *SetPlayerGiftReq) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "%d|%d", &msg.GiftId, &msg.PlayerId)
	return nil
}

type PlayersReq struct{}

func (msg *PlayersReq) Opcode() string                 { return "AL" }
func (msg *PlayersReq) Serialize(out io.Writer) error  { return nil }
func (msg *PlayersReq) Deserialize(in io.Reader) error { return nil }
