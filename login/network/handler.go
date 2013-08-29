package network

import (
	"database/sql"
	"github.com/Blackrush/gofus/login/db"
	"github.com/Blackrush/gofus/protocol/msg"
	"github.com/Blackrush/gofus/protocol/types"
	"github.com/Blackrush/gofus/shared"
	"log"
	"time"
)

func handle_client_connection(ctx *context, client Client) {
	ctx.clients[client.Id()] = client
	log.Printf("client #%04d [connected]", client.Id())

	client.Send(&msg.HelloConnect{Ticket: client.Ticket()})
	client.SetState(VersionState)
}

func handle_client_disconnection(ctx *context, client Client) {
	delete(ctx.clients, client.Id())
	log.Printf("client #%04d [disconnected]", client.Id())
}

func handle_client_data(ctx *context, client Client, data string) {
	log.Printf("client #%04d [%03d bytes]<<<%s", client.Id(), len(data), data)

	switch client.State() {
	case VersionState:
		if clientVersion == data {
			client.SetState(LoginState)
		} else {
			client.CloseWith(&msg.BadVersion{Required: clientVersion})
		}
	case LoginState:
		username, pass := shared.Split2(data, "\n")
		handle_client_auth(ctx, client, username, pass)
	case RealmState:
		switch data[:2] {
		case "Af":
			handle_client_queue_status(ctx, client)
		case "Ax":
			handle_client_player_list(ctx, client)
		case "AX":
		}
	default:
		// No need to panic, it's only one client who got lost; just log and kick him out
		log.Printf("client #%04d [unknown state]", client.Id())
		client.Close()
	}
}

func authenticate(ctx *context, client Client, user *db.User, pass string) bool {
	clear := shared.DecryptDofusPassword(pass[2:], client.Ticket())

	if clear != user.Password {
		client.CloseWith(&msg.LoginError{})
		return false
	} else if !user.Rights.Has(db.NoneRight) {
		client.CloseWith(&msg.BannedUser{})
		return false
	}

	client.SetUser(user)
	return true
}

func handle_client_auth(ctx *context, client Client, username, pass string) {
	user, err := ctx.users.FindByName(username)
	if err == sql.ErrNoRows {
		client.CloseWith(&msg.LoginError{})
	} else if err != nil {
		log.Print("can't log '", username, "' because: ", err)
		client.Close()
	} else if authenticate(ctx, client, user, pass) {
		client.Send(&msg.SetNickname{user.Nickname})
		client.Send(&msg.SetCommunity{0}) // (fr)
		client.Send(&msg.SetRealmServers{[]*types.RealmServer{
			&types.RealmServer{
				Id:         1, // jiva (fr)
				State:      types.RealmOnlineState,
				Completion: 0,
				Joinable:   true,
			},
		}})
		client.Send(&msg.LoginSuccess{IsAdmin: false})
		client.Send(&msg.SetSecretQuestion{user.SecretQuestion})

		client.SetState(RealmState)
	}
}

func handle_client_queue_status(ctx *context, client Client) {
}

func handle_client_player_list(ctx *context, client Client) {
	client.Send(&msg.SetRealmServerPlayers{
		SubscriptionEnd: time.Now().AddDate(1, 0, 0), // 1 year 0 month 0 day
		Players: []*types.RealmServerPlayers{
			&types.RealmServerPlayers{
				Id:      1,
				Players: 1,
			},
		},
	})
}
