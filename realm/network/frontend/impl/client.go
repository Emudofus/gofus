package impl

import (
	"bytes"
	"github.com/Blackrush/gofus/protocol/backend"
	protocol "github.com/Blackrush/gofus/protocol/frontend"
	"github.com/Blackrush/gofus/realm/db"
	"log"
	"net"
)

type net_client struct {
	net.Conn
	alive bool

	id        uint64
	userInfos *backend.UserInfos
	player    *db.Player
}

func new_net_client(conn net.Conn, id uint64) *net_client {
	return &net_client{
		Conn:  conn,
		alive: true,
		id:    id,
	}
}

func (client *net_client) Write(b []byte) (int, error) {
	log.Printf("[frontend-net-client-%04d] SND(%03d) %s", client.Id(), len(b), b)
	return client.Conn.Write(b)
}

func (client *net_client) Close() error {
	client.alive = false
	return client.Conn.Close()
}

func (client *net_client) Send(msg protocol.MessageContainer) (int, error) {
	buf := new(bytes.Buffer)
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

func (client *net_client) UserInfos() *backend.UserInfos {
	return client.userInfos
}

func (client *net_client) SetUserInfos(userInfos backend.UserInfos) {
	client.userInfos = &userInfos
}

func (client *net_client) Player() *db.Player {
	return client.player
}

func (client *net_client) SetPlayer(player *db.Player) {
	client.player = player
}
