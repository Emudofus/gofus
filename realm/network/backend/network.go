package backend

import (
	digest "crypto/sha512"
	"fmt"
	"github.com/Blackrush/gofus/protocol/backend"
	frontend "github.com/Blackrush/gofus/protocol/frontend/types"
	"github.com/Blackrush/gofus/realm/db"
	"github.com/Blackrush/gofus/shared"
	"hash"
	"io"
	"log"
	"net"
)

type Configuration struct {
	ServerId         uint
	ServerAddr       string
	ServerPort       uint16
	ServerCompletion int
	Laddr            string
	Password         string
}

type Service interface {
	shared.StartStopper

	SetState(state frontend.RealmServerState)
	GetUserInfos(ticket string) (*backend.UserInfos, bool)
	NotifyUserConnection(userId uint64, connected bool)
}

type context struct {
	config  Configuration
	players *db.Players

	running      bool
	conn         net.Conn
	pendingUsers map[string]*backend.UserInfos // ticket -> user infos
}

func New(players *db.Players, config Configuration) Service {
	return &context{
		config:       config,
		players:      players,
		pendingUsers: make(map[string]*backend.UserInfos),
	}
}

func (ctx *context) Start() {
	if ctx.running {
		panic("backend service is already running")
	}
	ctx.running = true

	go client_connect(ctx)

	log.Print("[backend-net] successfully started")
}

func (ctx *context) Stop() {
	ctx.running = false
}

func shexdigest(digest hash.Hash, input string) []byte {
	return digest.Sum([]byte(input))
}

func (ctx *context) get_password_hash(salt string) []byte {
	d := digest.New()
	return shexdigest(d, fmt.Sprintf("%x%s", shexdigest(d, ctx.config.Password), salt))
}

func (ctx *context) send(msg backend.Message) error {
	log.Printf("[backend-net] SND(%d)", msg.Opcode())
	backend.Put(ctx.conn, msg.Opcode())
	return msg.Serialize(ctx.conn)
}

func (ctx *context) SetState(state frontend.RealmServerState) {
	go ctx.send(&backend.SetStateMsg{state}) // make it fast
}

func (ctx *context) GetUserInfos(ticket string) (*backend.UserInfos, bool) {
	if infos, ok := ctx.pendingUsers[ticket]; ok {
		delete(ctx.pendingUsers, ticket)
		return infos, true
	}
	return nil, false
}

func (ctx *context) NotifyUserConnection(userId uint64, connected bool) {
	go ctx.send(&backend.UserConnectedMsg{userId, connected})
}

func conn_rcv(conn net.Conn) (backend.Message, bool) {
	var opcode uint16

	if n, err := backend.Read(conn, &opcode); n <= 0 || err == io.EOF {
		return nil, false
	} else if err != nil {
		panic(err.Error())
	}

	if msg, ok := backend.NewMsg(opcode); ok {
		if err := msg.Deserialize(conn); err == io.EOF {
			return nil, false
		} else if err != nil {
			panic(err.Error())
		}

		return msg, true
	}

	log.Printf("[backend-net] unknown opcode %d", opcode)

	return nil, false
}

func client_connect(ctx *context) {
	conn, err := net.Dial("tcp", ctx.config.Laddr)
	if err != nil {
		panic(err.Error())
	}

	ctx.conn = conn
	defer client_close(ctx)

	log.Printf("[backend-net] connected to %s", ctx.config.Laddr)
	client_connection(ctx)

	for ctx.running {
		if msg, ok := conn_rcv(conn); ok {
			log.Printf("[backend-net] RCV(%d)", msg.Opcode())
			client_handle_data(ctx, msg)
		} else {
			break
		}
	}
}

func client_close(ctx *context) {
	ctx.conn.Close()
	client_disconnection(ctx)
}
