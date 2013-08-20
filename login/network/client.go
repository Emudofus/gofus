package network

import (
	"bytes"
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

	Id() int
	Ticket() string
	State() ClientState
	SetState(state ClientState)
}

type netClient struct {
	net.Conn

	id int
	ticket string
	state ClientState
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

func (client *netClient) Write(data []byte) (int, error) {
	log.Printf("outcoming %d bytes to client #%d (%s)\n", len(data), client.id, data)
	return client.Conn.Write(data)
}

func (client *netClient) Send(msg protocol.MessageContainer) (int, error) {
	buf := bytes.NewBuffer(nil)
	if err := msg.Serialize(buf); err != nil {
		return 0, err
	}
	buf.WriteString(messageDelimiter)

	n, err := buf.WriteTo(client)
	return int(n), err
}

func NewNetClient(conn net.Conn, id int, ticket string) Client {
	return &netClient {
		conn,
		id,
		ticket,
		NoneState,
	}
}
