package handler

import (
	protocol "github.com/Blackrush/gofus/protocol/frontend"
	"github.com/Blackrush/gofus/protocol/frontend/msg"
	"github.com/Blackrush/gofus/realm/db"
	"github.com/Blackrush/gofus/realm/network/frontend"
	"log"
)

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

func handle_client_game_context_create(ctx frontend.Service, client frontend.Client, m *msg.GameContextCreateReq) {
	switch m.Type {
	case msg.SoloContextType:

	case msg.FightContextType:
		fallthrough
	default:
		log.Print("context type %d is not implemented yet", m.Type)
	}
}
