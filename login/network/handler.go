package network

import (
	"github.com/Blackrush/gofus/protocol/msg"
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

	switch client.State() {
	case VersionState:
		if clientVersion == data {
			client.SetState(LoginState)
		} else {
			client.Send(&msg.BadVersion{ Required: clientVersion })
			client.Close()
		}
	case LoginState:
	case RealmState:
	default:
		// No need to panic, it's only one client who got lost; just log and kick him out
		log.Print("unknown state ", client.State(), " of client #", client.Id())
		client.Close()
	}
}
