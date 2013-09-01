package backend

import (
	"github.com/Blackrush/gofus/protocol/frontend/types"
	"github.com/Blackrush/gofus/shared/db"
	"io"
	"time"
)

var producers = make(map[uint16]func() Message)

func init() {
	producers[1] = func() Message { return new(HelloConnectMsg) }
	producers[2] = func() Message { return new(AuthReqMsg) }
	producers[3] = func() Message { return new(AuthRespMsg) }
	producers[4] = func() Message { return new(SetInfosMsg) }
	producers[5] = func() Message { return new(SetStateMsg) }
	producers[6] = func() Message { return new(ClientConnMsg) }
	producers[7] = func() Message { return new(ClientConnReadyMsg) }
	producers[8] = func() Message { return new(UserPlayersReqMsg) }
	producers[9] = func() Message { return new(UserPlayersRespMsg) }
}

func NewMsg(opcode uint16) (Message, bool) {
	if producer, ok := producers[opcode]; ok {
		return producer(), true
	}
	return nil, false
}

type HelloConnectMsg struct {
	Salt string
}

func (msg *HelloConnectMsg) Opcode() uint16 { return 1 }
func (msg *HelloConnectMsg) Serialize(out io.Writer) error {
	Put(out, msg.Salt)
	return nil
}
func (msg *HelloConnectMsg) Deserialize(in io.Reader) error {
	Read(in, &msg.Salt)
	return nil
}

type AuthReqMsg struct {
	Id          uint16
	Credentials []byte
}

func (msg *AuthReqMsg) Opcode() uint16 { return 2 }
func (msg *AuthReqMsg) Serialize(out io.Writer) error {
	Put(out, msg.Id)
	Put(out, uint32(len(msg.Credentials)))
	Put(out, msg.Credentials)
	return nil
}
func (msg *AuthReqMsg) Deserialize(in io.Reader) error {
	Read(in, &msg.Id)
	var tmp uint32
	Read(in, &tmp)
	msg.Credentials = make([]byte, tmp)
	Read(in, msg.Credentials)
	return nil
}

type AuthRespMsg struct {
	Success bool
}

func (msg *AuthRespMsg) Opcode() uint16 { return 3 }
func (msg *AuthRespMsg) Serialize(out io.Writer) error {
	Put(out, msg.Success)
	return nil
}
func (msg *AuthRespMsg) Deserialize(in io.Reader) error {
	Read(in, &msg.Success)
	return nil
}

type SetInfosMsg struct {
	Address    string
	Port       uint16
	Completion uint32
}

func (msg *SetInfosMsg) Opcode() uint16 { return 4 }
func (msg *SetInfosMsg) Serialize(out io.Writer) error {
	Put(out, msg.Address)
	Put(out, msg.Port)
	Put(out, msg.Completion)
	return nil
}
func (msg *SetInfosMsg) Deserialize(in io.Reader) error {
	Read(in, &msg.Address)
	Read(in, &msg.Port)
	Read(in, &msg.Completion)
	return nil
}

type SetStateMsg struct {
	State types.RealmServerState
}

func (msg *SetStateMsg) Opcode() uint16 { return 5 }
func (msg *SetStateMsg) Serialize(out io.Writer) error {
	Put(out, uint8(msg.State))
	return nil
}
func (msg *SetStateMsg) Deserialize(in io.Reader) error {
	var tmp uint8
	Read(in, &tmp)
	msg.State = types.RealmServerState(tmp)
	return nil
}

type UserInfos struct {
	Id              uint64
	SecretQuestion  string
	SecretAnswer    string // TODO security
	SubscriptionEnd time.Time
	Rights          db.UserRight
}

func (msg *UserInfos) Serialize(out io.Writer) error {
	Put(out, msg.Id)
	Put(out, msg.SecretQuestion)
	Put(out, msg.SecretAnswer)
	Put(out, msg.SubscriptionEnd)
	Put(out, int64(msg.Rights))
	return nil
}
func (msg *UserInfos) Deserialize(in io.Reader) error {
	Read(in, &msg.Id)
	Read(in, &msg.SecretQuestion)
	Read(in, &msg.SecretAnswer)
	Read(in, &msg.SubscriptionEnd)
	var tmp uint64
	Read(in, &tmp)
	msg.Rights = db.UserRight(tmp)
	return nil
}

type ClientConnMsg struct {
	Ticket string
	User   UserInfos
}

func (msg *ClientConnMsg) Opcode() uint16 { return 6 }
func (msg *ClientConnMsg) Serialize(out io.Writer) error {
	Put(out, msg.Ticket)
	Put(out, &msg.User)
	return nil
}
func (msg *ClientConnMsg) Deserialize(in io.Reader) error {
	Read(in, &msg.Ticket)
	Read(in, &msg.User)
	return nil
}

type ClientConnReadyMsg struct {
	Ticket string
}

func (msg *ClientConnReadyMsg) Opcode() uint16 { return 7 }
func (msg *ClientConnReadyMsg) Serialize(out io.Writer) error {
	Put(out, msg.Ticket)
	return nil
}
func (msg *ClientConnReadyMsg) Deserialize(in io.Reader) error {
	Read(in, &msg.Ticket)
	return nil
}

type UserPlayersReqMsg struct {
	UserId uint64
}

func (msg *UserPlayersReqMsg) Opcode() uint16 { return 8 }
func (msg *UserPlayersReqMsg) Serialize(out io.Writer) error {
	Put(out, msg.UserId)
	return nil
}
func (msg *UserPlayersReqMsg) Deserialize(in io.Reader) error {
	Read(in, &msg.UserId)
	return nil
}

type UserPlayersRespMsg struct {
	UserId  uint64
	Players uint8
}

func (msg *UserPlayersRespMsg) Opcode() uint16 { return 9 }
func (msg *UserPlayersRespMsg) Serialize(out io.Writer) error {
	Put(out, msg.UserId)
	Put(out, msg.Players)
	return nil
}
func (msg *UserPlayersRespMsg) Deserialize(in io.Reader) error {
	Read(in, &msg.UserId)
	Read(in, &msg.Players)
	return nil
}
