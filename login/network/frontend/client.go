package frontend

import (
	"bytes"
	"github.com/Blackrush/gofus/login/db"
	protocol "github.com/Blackrush/gofus/protocol/frontend"
	"io"
	"log"
	"net"
)

type ClientState uint8

func (state *ClientState) Inc() {
	*state = ClientState(*state + 1)
}

const (
	ClientNoneState ClientState = iota
	ClientLoginState
	ClientRealmState
)

type Client interface {
	io.WriteCloser
	protocol.Sender
	protocol.CloseWither
	Alive() bool

	Id() uint64
	Ticket() string
	State() *ClientState
	User() *db.User
	SetUser(user *db.User)
}

type net_client struct {
	conn  net.Conn
	alive bool

	id     uint64
	ticket string
	state  ClientState
	user   *db.User
}

func NewNetClient(conn net.Conn, id uint64, ticket string) Client {
	return &net_client{
		conn:   conn,
		alive:  true,
		id:     id,
		ticket: ticket,
		state:  ClientNoneState,
		user:   nil,
	}
}

func (client *net_client) Write(b []byte) (int, error) {
	log.Printf("[frontend-net-client-%04d] SND(%03d) %s", client.id, len(b), b)
	return client.conn.Write(b)
}

func (client *net_client) Close() error {
	client.alive = false
	return client.conn.Close()
}

func (client *net_client) Send(msg protocol.MessageContainer) (int, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(msg.Opcode())
	msg.Serialize(buf)
	buf.WriteString(output_msg_delimiter)

	n, err := buf.WriteTo(client)
	return int(n), err
}

func (client *net_client) CloseWith(msg protocol.MessageContainer) error {
	client.Send(msg)
	return client.Close()
}

func (client *net_client) Alive() bool {
	return client.alive
}

func (client *net_client) Id() uint64 {
	return client.id
}

func (client *net_client) Ticket() string {
	return client.ticket
}

func (client *net_client) State() *ClientState {
	return &client.state
}

func (client *net_client) User() *db.User {
	return client.user
}

func (client *net_client) SetUser(user *db.User) {
	client.user = user
}
