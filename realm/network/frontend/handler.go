package frontend

import (
	"github.com/Blackrush/gofus/protocol/backend"
	protocol "github.com/Blackrush/gofus/protocol/frontend"
	"github.com/Blackrush/gofus/protocol/frontend/msg"
	"github.com/Blackrush/gofus/realm/db"
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
	case *msg.RandNameReq:
		handle_client_rand_name(ctx, client, m)
	case *msg.CreatePlayerReq:
		handle_client_create_player(ctx, client, m)
	case *msg.PlayerSelectionReq:
		handle_client_player_selection(ctx, client, m)
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
	client_send_player_list(ctx, client)
}

func handle_client_rand_name(ctx *context, client *net_client, m *msg.RandNameReq) {
	client.Send(&msg.RandNameResp{Name: "Gofus"})
}

func create_default_player(user *backend.UserInfos, m *msg.CreatePlayerReq) *db.Player {
	player := &db.Player{
		OwnerId: uint(user.Id),
		Name:    m.Name,
	}

	player.Appearance.Skin = m.Breed * 10
	player.Appearance.Colors.First = db.PlayerColor(m.Colors.First)
	player.Appearance.Colors.Second = db.PlayerColor(m.Colors.Second)
	player.Appearance.Colors.Third = db.PlayerColor(m.Colors.Third)

	player.Experience = db.PlayerExperience{
		Level:      1,
		Experience: 0,
	}

	return player
}

func handle_client_create_player(ctx *context, client *net_client, m *msg.CreatePlayerReq) {
	player := create_default_player(client.UserInfos(), m)

	if inserted, success := ctx.players.Persist(player); inserted && success {
		client_send_player_list(ctx, client)
	} else {
		client.Send(&msg.CreatePlayerErrorResp{})
	}
}

func handle_client_player_selection(ctx *context, client *net_client, m *msg.PlayerSelectionReq) {

}
