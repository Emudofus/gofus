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
	log.Print("new client connected with id=", client.Id())

	client.Send(&msg.HelloConnect{ Ticket: client.Ticket() })
	client.SetState(VersionState)
}

func handle_client_disconnection(ctx *context, client Client) {
	delete(ctx.clients, client.Id())
	log.Print("client #", client.Id(), " disconnected")
}

func handle_client_data(ctx *context, client Client, data string) {
	log.Print("received ", len(data), " bytes `", data, "` from client #", client.Id())

	users := db.Users{ctx.db}

	switch client.State() {
	case VersionState:
		if clientVersion == data {
			client.SetState(LoginState)
		} else {
			client.Send(&msg.BadVersion{ Required: clientVersion })
			client.Close()
		}
	case LoginState:
		username, _ := shared.Split2(data, "\n")
		log.Print("user '", username, "' wants to login")

		_, err := users.FindByName(username)
		if err == sql.ErrNoRows {
			client.Send(&msg.LoginError{})
			client.Close()
		} else if err != nil {
			log.Print("can't log '", username, "' because: ", err)
			client.Close()
		} else {

		}
	case RealmState:
	default:
		// No need to panic, it's only one client who got lost; just log and kick him out
		log.Print("unknown state ", client.State(), " of client #", client.Id())
		client.Close()
	}
}
