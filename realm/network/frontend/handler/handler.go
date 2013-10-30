package handler

import (
	protocol "github.com/Blackrush/gofus/protocol/frontend"
	"github.com/Blackrush/gofus/protocol/frontend/msg"
	"github.com/Blackrush/gofus/realm/db"
	"github.com/Blackrush/gofus/realm/network/frontend"
	"log"
)

func client_send_player_list(ctx frontend.Service, client frontend.Client) {
	var players []*db.Player
	var ok bool
	if players, ok = ctx.Players().GetByOwnerId(uint(client.UserInfos().Id)); !ok {
		players = make([]*db.Player, 0, 0)
	}

	client.Send(&msg.PlayersResp{
		ServerId:        ctx.Config().ServerId,
		SubscriptionEnd: client.UserInfos().SubscriptionEnd,
		Players:         players,
	})
}

func HandleClientConnection(ctx frontend.Service, client frontend.Client) {
	log.Printf("[frontend-net-client-%04d] CONN", client.Id())

	client.Send(&msg.HelloGame{})
}

func HandleClientDisconnection(ctx frontend.Service, client frontend.Client) {
	log.Printf("[frontend-net-client-%04d] DCONN", client.Id())

	if client.UserInfos() != nil {
		ctx.Backend().NotifyUserConnection(client.UserInfos().Id, false)
	}
}

func HandleClientData(ctx frontend.Service, client frontend.Client, arg protocol.MessageContainer) {
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
	case *msg.GameContextCreateReq:
		handle_client_game_context_create(ctx, client, m)
	}
}

func handle_client_login(ctx frontend.Service, client frontend.Client, m *msg.RealmLoginReq) {
	if infos, ok := ctx.Backend().GetUserInfos(m.Ticket); ok {
		client.SetUserInfos(*infos)
		ctx.Backend().NotifyUserConnection(infos.Id, true)

		client.Send(&msg.RealmLoginSuccess{ctx.Config().CommunityId})
	} else {
		client.CloseWith(&msg.RealmLoginError{})
	}
}

func handle_client_regional_version(ctx frontend.Service, client frontend.Client, m *msg.RegionalVersionReq) {
	client.Send(&msg.RegionalVersionResp{ctx.Config().CommunityId})
}

func handle_client_players(ctx frontend.Service, client frontend.Client, m *msg.PlayersReq) {
	client_send_player_list(ctx, client)
}

func handle_client_rand_name(ctx frontend.Service, client frontend.Client, m *msg.RandNameReq) {
	client.Send(&msg.RandNameResp{Name: "Gofus"})
}

func handle_client_create_player(ctx frontend.Service, client frontend.Client, m *msg.CreatePlayerReq) {
	player := ctx.Players().NewPlayer(uint(client.UserInfos().Id), m.Name, m.Breed, m.Gender, m.Colors.First, m.Colors.Second, m.Colors.Third)

	if inserted, success := ctx.Players().Persist(player); inserted && success {
		client_send_player_list(ctx, client)
	} else {
		client.Send(&msg.CreatePlayerErrorResp{})
	}
}

func handle_client_player_selection(ctx frontend.Service, client frontend.Client, m *msg.PlayerSelectionReq) {
	player, ok := ctx.Players().GetById(m.PlayerId)
	if !ok {
		client.CloseWith(&msg.PlayerSelectionErrorResp{})
		return
	}

	client.SetPlayer(player)
	client.Send(&msg.PlayerSelectionResp{player})
}

func handle_client_game_context_create(ctx frontend.Service, client frontend.Client, m *msg.GameContextCreateReq) {
	log.Print("YOSH")
}
