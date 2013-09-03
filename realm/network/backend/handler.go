package backend

import (
	"github.com/Blackrush/gofus/protocol/backend"
	frontend "github.com/Blackrush/gofus/protocol/frontend/types"
	"log"
)

func client_connection(ctx *context) {

}

func client_disconnection(ctx *context) {

}

func client_handle_data(ctx *context, arg backend.Message) {
	switch msg := arg.(type) {
	case *backend.HelloConnectMsg:
		client_handle_hello_connect(ctx, msg)
	case *backend.AuthRespMsg:
		client_handle_auth(ctx, msg)
	case *backend.UserPlayersReqMsg:
		client_handle_user_players(ctx, msg)
	case *backend.ClientConnMsg:
		client_handle_client_connection(ctx, msg)
	}
}

func client_handle_hello_connect(ctx *context, msg *backend.HelloConnectMsg) {
	ctx.send(&backend.AuthReqMsg{
		Id:          uint16(ctx.config.ServerId),
		Credentials: ctx.get_password_hash(msg.Salt),
	})
}

func client_handle_auth(ctx *context, msg *backend.AuthRespMsg) {
	if msg.Success {
		log.Printf("[backend-net] successfully synchronized")

		ctx.send(&backend.SetInfosMsg{
			Address:    ctx.config.ServerAddr,
			Port:       ctx.config.ServerPort,
			Completion: uint32(ctx.config.ServerCompletion),
		})

		ctx.SetState(frontend.RealmOnlineState)
	} else {
		panic("can't authenticate to login server")
	}
}

func client_handle_user_players(ctx *context, msg *backend.UserPlayersReqMsg) {
	ctx.send(&backend.UserPlayersRespMsg{ // TODO
		UserId:  msg.UserId,
		Players: 1,
	})
}

func client_handle_client_connection(ctx *context, msg *backend.ClientConnMsg) {
	ctx.pendingUsers[msg.Ticket] = &msg.User

	ctx.send(&backend.ClientConnReadyMsg{
		Ticket: msg.Ticket,
	})
}
