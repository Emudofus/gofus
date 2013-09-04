package backend

import (
	"bytes"
	"fmt"
	"github.com/Blackrush/gofus/login/db"
	"github.com/Blackrush/gofus/protocol/backend"
	frontend "github.com/Blackrush/gofus/protocol/frontend/types"
	"log"
)

type Realm struct {
	frontend.RealmServer
	Address string
	Port    uint16

	client            *Client
	conn_callbacks    map[string]chan bool // ticket -> callback
	players_callbacks map[uint]chan frontend.RealmServerPlayers
}

func (realm *Realm) AssertJoinable() {
	if !realm.Joinable {
		log.Panic(fmt.Sprintf("realm %d is not joinable", realm.Id))
	}
}

func (realm *Realm) NotifyUserConnection(ticket string, user *db.User) (callback chan bool) {
	realm.AssertJoinable()

	if _, exists := realm.conn_callbacks[ticket]; exists {
		return // TODO prevent any flood
	}

	realm.client.Send(&backend.ClientConnMsg{
		Ticket: ticket,
		User: backend.UserInfos{
			Id:              uint64(user.Id),
			SecretQuestion:  user.SecretQuestion,
			SecretAnswer:    user.SecretAnswer,
			SubscriptionEnd: user.SubscriptionEnd,
			Rights:          user.Rights,
		},
	})

	callback = make(chan bool, 1)
	realm.conn_callbacks[ticket] = callback
	return
}

func (realm *Realm) AskPlayers(userId uint, callback chan frontend.RealmServerPlayers) {
	if !realm.Joinable {
		callback <- frontend.RealmServerPlayers{
			Id:      realm.Id,
			Players: 0,
		}
		return
	}

	if _, exists := realm.players_callbacks[userId]; exists {
		return // TODO prevent any flood
	}

	realm.client.Send(&backend.UserPlayersReqMsg{uint64(userId)})
	realm.players_callbacks[userId] = callback
}

func client_connection(ctx *context, client *Client) {
	client.Send(&backend.HelloConnectMsg{client.salt})
}

func client_disconnection(ctx *context, client *Client) {
	if client.realm != nil {
		log.Printf("[realm-%02d] is now offline", client.realm.Id)

		client.realm.State = frontend.RealmOfflineState
		client.realm.Joinable = false

		for _, callback := range client.realm.conn_callbacks {
			close(callback)
		}
		client.realm.conn_callbacks = nil

		for _, callback := range client.realm.players_callbacks {
			close(callback)
		}
		client.realm.players_callbacks = nil

		client.realm.client = nil
		client.realm = nil
	}
}

func client_handle_data(ctx *context, client *Client, arg backend.Message) {
	switch msg := arg.(type) {
	case *backend.AuthReqMsg:
		client_handle_auth(ctx, client, msg)
	case *backend.SetInfosMsg:
		client_handle_set_infos(ctx, client, msg)
	case *backend.SetStateMsg:
		client_handle_set_state(ctx, client, msg)
	case *backend.ClientConnReadyMsg:
		client_handle_client_conn_ready(ctx, client, msg)
	case *backend.UserPlayersRespMsg:
		client_handle_user_players(ctx, client, msg)
	case *backend.UserConnectedMsg:
		client_handle_user_connected(ctx, client, msg)
	}
}

func client_authenticate(ctx *context, client *Client, credentials []byte) bool {
	return bytes.Equal(ctx.get_password_hash(client.salt), credentials)
}

func client_handle_auth(ctx *context, client *Client, msg *backend.AuthReqMsg) {
	if client.realm != nil {
		log.Printf("[realm-%02d] tried to reauth", client.realm.Id)
		return
	}

	realm, exists := ctx.realms[int(msg.Id)]

	if exists && realm.Joinable {
		goto failure
	}

	if client_authenticate(ctx, client, msg.Credentials) {
		if !exists {
			realm = new(Realm)
			realm.Id = int(msg.Id)
			ctx.realms[realm.Id] = realm
		}
		realm.client = client
		realm.conn_callbacks = make(map[string]chan bool)
		realm.players_callbacks = make(map[uint]chan frontend.RealmServerPlayers)

		client.realm = realm

		client.Send(&backend.AuthRespMsg{Success: true})

		log.Printf("[realm-%02d] is now synchronized", client.realm.Id)
		return
	}

failure: // maybe there is a better way to do
	client.Send(&backend.AuthRespMsg{Success: false})
	client.Close()
}

func client_handle_set_infos(ctx *context, client *Client, msg *backend.SetInfosMsg) {
	client.realm.Address = msg.Address
	client.realm.Port = msg.Port
	client.realm.Completion = int(msg.Completion)

	log.Printf("[realm-%02d] updated his infos", client.realm.Id)
}

func client_handle_set_state(ctx *context, client *Client, msg *backend.SetStateMsg) {
	client.realm.State = msg.State
	client.realm.Joinable = msg.State == frontend.RealmOnlineState

	log.Printf("[realm-%02d] updated his state, now %d", client.realm.Id, client.realm.State)
}

func client_handle_client_conn_ready(ctx *context, client *Client, msg *backend.ClientConnReadyMsg) {
	if callback, ok := client.realm.conn_callbacks[msg.Ticket]; ok {
		callback <- true
		delete(client.realm.conn_callbacks, msg.Ticket)
	} else {
		log.Printf("[realm-%02d] tried to allow a unknown client connection", client.realm.Id)
	}
}

func client_handle_user_players(ctx *context, client *Client, msg *backend.UserPlayersRespMsg) {
	id := uint(msg.UserId)

	if callback, ok := client.realm.players_callbacks[id]; ok {
		callback <- frontend.RealmServerPlayers{
			Id:      client.realm.Id,
			Players: int(msg.Players),
		}
		delete(client.realm.players_callbacks, id)
	} else {
		log.Printf("[realm-%02d] tried to give players but wasn't necessary", client.realm.Id)
	}
}

func client_handle_user_connected(ctx *context, client *Client, msg *backend.UserConnectedMsg) {
	if user, err := ctx.users.FindById(int(msg.UserId)); err == nil {
		if msg.Connected {
			if user.IsConnected() {
				log.Printf("[realm-%02d] user %d is already connected", user.Id)
			} else {
				user.CurrentRealmServerId = client.realm.Id
				ctx.users.Update(user)
			}
		} else {
			if user.CurrentRealmServerId != client.realm.Id {
				log.Printf("[realm-%02d] has no right to mark user %d as disconnected", client.realm.Id, user.Id)
			} else {
				user.CurrentRealmServerId = -1
				ctx.users.Update(user)
			}
		}
	} else {
		log.Panic(err.Error())
	}
}
