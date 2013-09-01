package network

import (
	protocol "github.com/Blackrush/gofus/protocol/frontend"
	"github.com/Blackrush/gofus/protocol/frontend/msg"
	"log"
)

func handle_client_connection(ctx *context, client *net_client) {
	log.Printf("[frontend-net-client-%04d] CONN", client.Id())

	client.Send(&msg.HelloGame{})
}

func handle_client_disconnection(ctx *context, client *net_client) {
	log.Printf("[frontend-net-client-%04d] DCONN", client.Id())
}

func handle_client_data(ctx *context, client *net_client, arg protocol.MessageContainer) {
	log.Printf("[frontend-net-client-%04d] RCV(%s) %+v", client.Id(), arg.Opcode(), arg)

	switch m := arg.(type) {
	case *msg.RealmLoginReq:
		handle_client_login(ctx, client, m)
	}
}

func handle_client_login(ctx *context, client *net_client, m *msg.RealmLoginReq) {
	if infos, ok := ctx.backend.GetUserInfos(m.Ticket); ok {
		client.SetUserInfos(*infos)

		client.Send(&msg.RealmLoginSuccess{})
	} else {
		client.CloseWith(&msg.RealmLoginError{})
	}
}
