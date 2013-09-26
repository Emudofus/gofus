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

	if client.UserInfos() != nil {
		ctx.backend.NotifyUserConnection(client.UserInfos().Id, false)
	}
}

func handle_client_data(ctx *context, client *net_client, arg protocol.MessageContainer) {
	switch m := arg.(type) {
	case *msg.RealmLoginReq:
		handle_client_login(ctx, client, m)
	case *msg.RegionalVersionReq:
		handle_client_regional_version(ctx, client, m)
	case *msg.PlayersReq:
		handle_client_players(ctx, client, m)
	}
}

func handle_client_login(ctx *context, client *net_client, m *msg.RealmLoginReq) {
	if infos, ok := ctx.backend.GetUserInfos(m.Ticket); ok {
		client.SetUserInfos(*infos)
		ctx.backend.NotifyUserConnection(infos.Id, true)

		client.Send(&msg.RealmLoginSuccess{ctx.config.CommunityId})
	} else {
		client.CloseWith(&msg.RealmLoginError{})
	}
}

func handle_client_regional_version(ctx *context, client *net_client, m *msg.RegionalVersionReq) {
	client.Send(&msg.RegionalVersionResp{ctx.config.CommunityId})
}

func handle_client_players(ctx *context, client *net_client, m *msg.PlayersReq) {
	if players, ok := ctx.players.GetByOwnerId(uint(client.userInfos.Id)); ok {
		client.Send(&msg.PlayersResp{
			ServerId:        ctx.config.ServerId,
			SubscriptionEnd: client.userInfos.SubscriptionEnd,
			Players:         players,
		})
	}
}
