package network

import (
	"bytes"
	"fmt"
	"github.com/Blackrush/gofus/login/db"
	"github.com/Blackrush/gofus/protocol"
	"io"
	"log"
	"net"
)

type ClientState int

func (state ClientState) Next() ClientState {
	return ClientState(int(state) + 1)
}

func (state *ClientState) Increment() {
	*state = state.Next()
}

const (
	NoneState ClientState = iota
	VersionState
	LoginState
	RealmState
)

type Client interface {
	io.WriteCloser
	protocol.Sender

	Alive() bool
	CloseWith(msg protocol.MessageContainer) (n int, err error)

	Id() int
	Ticket() string
	State() ClientState
	SetState(state ClientState)
	User() *db.User
	SetUser(user *db.User)
}

type netClient struct {
	net.Conn
	alive bool

	id     int
	ticket string
	state  ClientState
	user   *db.User
}

func (client *netClient) Alive() bool {
	return client.alive
}

func (client *netClient) Id() int {
	return client.id
}

func (client *netClient) Ticket() string {
	return client.ticket
}

func (client *netClient) State() ClientState {
	return client.state
}

func (client *netClient) SetState(state ClientState) {
	client.state = state
}

func (client *netClient) User() *db.User {
	return client.user
}

func (client *netClient) SetUser(user *db.User) {
	client.user = user
}

func (client *netClient) Close() error {
	client.alive = false
	return client.Conn.Close()
}

func (client *netClient) Write(data []byte) (int, error) {
	log.Printf("client #%04d (%03d bytes)>>>%s\n", client.id, len(data), data)
	return client.Conn.Write(data)
}

func (client *netClient) Send(msg protocol.MessageContainer) (int, error) {
	buf := bytes.NewBuffer(nil)
	fmt.Fprint(buf, msg.Opcode())
	if err := msg.Serialize(buf); err != nil {
		return 0, err
	}
	buf.WriteString(messageDelimiter)

	n, err := buf.WriteTo(client)
	return int(n), err
}

func (client *netClient) CloseWith(msg protocol.MessageContainer) (n int, err error) {
	if n, err = client.Send(msg); err != nil {
		return
	} else if err = client.Close(); err != nil {
		return
	} else {
		return n, nil
	}
}

func NewNetClient(conn net.Conn, id int, ticket string) Client {
	return &netClient{
		Conn:   conn,
		alive:  true,
		id:     id,
		ticket: ticket,
		state:  NoneState,
		user:   nil,
	}
}
