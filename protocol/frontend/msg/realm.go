package msg

import (
	"fmt"
	"github.com/Blackrush/gofus/realm/db"
	"io"
	"strings"
	"time"
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

type SetIdentity struct {
	Identity string
}

func (msg *SetIdentity) Opcode() string { return "Ai" }
func (msg *SetIdentity) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Identity)
	return nil
}
func (msg *SetIdentity) Deserialize(in io.Reader) error {
	fmt.Fscanf(in, "Ai%s", &msg.Identity)
	return nil
}

type PlayersReq struct{}

func (msg *PlayersReq) Opcode() string                 { return "AL" }
func (msg *PlayersReq) Serialize(out io.Writer) error  { return nil }
func (msg *PlayersReq) Deserialize(in io.Reader) error { return nil }

type PlayersResp struct {
	ServerId        uint
	SubscriptionEnd time.Time
	Players         []*db.Player
}

func (msg *PlayersResp) Opcode() string { return "ALK" }
func (msg *PlayersResp) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%d|%d", msg.SubscriptionEnd.Sub(time.Now()).Nanoseconds()/1e6, len(msg.Players))
	for _, player := range msg.Players {
		fmt.Fprintf(out, "|%d;%s;%d;%d;%d;%d;%d;%s;%d;;;",
			player.Id,
			player.Name,
			player.Experience.Level,
			player.Appearance.Skin,
			player.Appearance.Colors.First,
			player.Appearance.Colors.Second,
			player.Appearance.Colors.Third,
			player.Appearance.Accessories,
			msg.ServerId,
		)
	}
	return nil
}
func (msg *PlayersResp) Deserialize(in io.Reader) error {
	return nil
}

type RandNameReq struct{}

func (msg *RandNameReq) Opcode() string                 { return "AP" }
func (msg *RandNameReq) Serialize(out io.Writer) error  { return nil }
func (msg *RandNameReq) Deserialize(in io.Reader) error { return nil }

type RandNameResp struct {
	Name string
}

func (msg *RandNameResp) Opcode() string { return "AP" }
func (msg *RandNameResp) Serialize(out io.Writer) error {
	fmt.Fprint(out, msg.Name)
	return nil
}
func (msg *RandNameResp) Deserialize(in io.Reader) error {
	return nil
}

type CreateUserReq struct {
	Name   string
	Breed  int
	Gender bool
	Colors struct {
		First  int
		Second int
		Third  int
	}
}

func (msg *CreateUserReq) Opcode() string { return "AA" }
func (msg *CreateUserReq) Serialize(out io.Writer) error {
	fmt.Fprintf(out, "%s|%d|%d|%d|%d|%d", msg.Name, msg.Breed, btoi(msg.Gender), msg.Colors.First, msg.Colors.Second, msg.Colors.Third)
	return nil
}
func (msg *CreateUserReq) Deserialize(in io.Reader) error {
	var body string
	fmt.Fscanf(in, "AA%s", &body)

	args := strings.SplitN(body, "|", 6)

	msg.Name = args[0]
	msg.Breed = atoi(args[1])
	msg.Gender = aitob(args[2])
	msg.Colors.First = atoi(args[3])
	msg.Colors.Second = atoi(args[4])
	msg.Colors.Third = atoi(args[5])

	return nil
}
