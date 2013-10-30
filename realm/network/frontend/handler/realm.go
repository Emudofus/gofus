package handler

import (
	"github.com/Blackrush/gofus/protocol/frontend/msg"
	"github.com/Blackrush/gofus/realm/db"
	"github.com/Blackrush/gofus/realm/network/frontend"
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
