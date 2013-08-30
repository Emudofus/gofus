package network

import (
	"bytes"
	"github.com/Blackrush/gofus/protocol"
	"io"
	"log"
	"net"
)

type ClientState int

func (state *ClientState) Increment() {
	*state = ClientState(*state + 1)
}

const (
	ClientNoneState ClientState = iota
	ClientAuthState
	ClientRealmState
	ClientWorldState
)

type Client interface {
	io.WriteCloser
	protocol.Sender
	CloseWith(msg protocol.MessageContainer) error

	Id() uint64
	Alive() bool
	State() ClientState
	SetState(state ClientState)
}

type netClient struct {
	net.Conn

	id    uint64
	alive bool
	state ClientState
}

func NewNetClient(conn net.Conn, id uint64) Client {
	return &netClient{
		Conn:  conn,
		id:    id,
		alive: true,
		state: ClientNoneState,
	}
}

func (client *netClient) Write(data []byte) (int, error) {
	log.Printf("[network-client-%04d] SND(%03d) %s", client.id, len(data), data)
	return client.Conn.Write(data)
}

func (client *netClient) Close() error {
	client.alive = false
	return client.Conn.Close()
}

func (client *netClient) Send(msg protocol.MessageContainer) (int, error) {
	buf := bytes.NewBuffer(nil)

	if n, err := buf.WriteString(msg.Opcode()); err != nil {
		return n, err
	}
	if err := msg.Serialize(buf); err != nil {
		return 0, err
	}
	if n, err := buf.WriteString(output_msg_delimiter); err != nil {
		return n, err
	}

	n, err := buf.WriteTo(client)
	return int(n), err
}

func (client *netClient) CloseWith(msg protocol.MessageContainer) error {
	if _, err := client.Send(msg); err != nil {
		return err
	}
	return client.Close()
}

func (client *netClient) Id() uint64 {
	return client.id
}

func (client *netClient) Alive() bool {
	return client.alive
}

func (client *netClient) State() ClientState {
	return client.state
}

func (client *netClient) SetState(state ClientState) {
	client.state = state
}
