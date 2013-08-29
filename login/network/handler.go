package network

import (
	"database/sql"
	"github.com/Blackrush/gofus/login/db"
	"github.com/Blackrush/gofus/protocol/msg"
	"github.com/Blackrush/gofus/shared"
	"log"
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

func authenticate(ctx *context, client Client, user *db.User, pass string) bool {
	clear := shared.DecryptDofusPassword(pass[2:], client.Ticket())
	println(clear)
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

func handle_client_data(ctx *context, client Client, data string) {
	log.Printf("client #%04d (%03d bytes)<<<%s", client.Id(), len(data), data)

	users := db.Users{ctx.db}

	switch client.State() {
	case VersionState:
		if clientVersion == data {
			client.SetState(LoginState)
		} else {
			client.Send(&msg.BadVersion{Required: clientVersion})
			client.Close()
		}
	case LoginState:
		username, pass := shared.Split2(data, "\n")
		log.Print("user '", username, "' wants to login")

		user, err := users.FindByName(username)
		if err == sql.ErrNoRows {
			client.CloseWith(&msg.LoginError{})
		} else if err != nil {
			log.Print("can't log '", username, "' because: ", err)
			client.Close()
		} else if authenticate(ctx, client, user, pass) {
			client.Send(&msg.SetNickname{user.Nickname})
			client.Send(&msg.SetCommunity{0})
			client.Send(&msg.LoginSuccess{false})
			client.Send(&msg.SetSecretQuestion{user.SecretQuestion})

			client.SetState(RealmState)
		}
	case RealmState:
	default:
		// No need to panic, it's only one client who got lost; just log and kick him out
		log.Print("unknown state ", client.State(), " of client #", client.Id())
		client.Close()
	}
}
