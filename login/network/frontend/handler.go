package frontend

import (
	"bytes"
	"database/sql"
	"github.com/Blackrush/gofus/login/db"
	"github.com/Blackrush/gofus/protocol/frontend/msg"
	"github.com/Blackrush/gofus/protocol/frontend/types"
	"github.com/Blackrush/gofus/shared"
	db2 "github.com/Blackrush/gofus/shared/db"
	"log"
)

func client_connection(ctx *context, client Client) {
	log.Printf("[frontend-net-client-%04d] CONN", client.Id())

	client.Send(&msg.HelloConnect{client.Ticket()})
}

func client_disconnection(ctx *context, client Client) {
	log.Printf("[frontend-net-client-%04d] DCONN", client.Id())
}

func client_handle_data(ctx *context, client Client, data []byte) {
	log.Printf("[frontend-net-client-%04d] RCV(%03d) %s", client.Id(), len(data), data)

	switch *client.State() {
	case ClientNoneState:
		client_handle_none_state(ctx, client, data)
	case ClientLoginState:
		client_handle_login_state(ctx, client, data)
	case ClientRealmState:
		switch string(data[:2]) {
		case "Ax":
			client_handle_players_list(ctx, client, data)
		case "Af":
			client_handle_queue_status(ctx, client, data)
		case "AX":
			client_handle_realm_selection(ctx, client, data)
		}
	}
}

func client_handle_none_state(ctx *context, client Client, data []byte) {
	version := string(data)

	if client_version != version {
		client.CloseWith(&msg.BadVersion{Required: client_version})
	} else {
		client.State().Inc()
	}
}

func client_authenticate(ctx *context, client Client, username, password string) (*db.User, bool) {
	user, err := ctx.users.FindByName(username)

	if err != nil {
		if err == sql.ErrNoRows {
			client.Send(&msg.LoginError{})
		} else {
			log.Print(err)
		}
		return nil, false
	} else if !user.ValidPassword(shared.DecryptDofusPassword(password, client.Ticket())) {
		client.Send(&msg.LoginError{})
		return nil, false
	} else if !user.Rights.Has(db2.LoginRight) {
		client.Send(&msg.BannedUser{})
		return nil, false
	}

	return user, true
}

// TODO perform some caching
func get_realm_servers(ctx *context) []*types.RealmServer {
	left := ctx.backend.GetRealms()
	right := make([]*types.RealmServer, len(left))
	for i, realm := range left {
		right[i] = &realm.RealmServer
	}
	return right
}

func client_handle_login_state(ctx *context, client Client, data []byte) {
	username, password := shared.Splits2(data, []byte("\n#1"))

	if user, ok := client_authenticate(ctx, client, username, password); ok {
		client.Send(&msg.SetCommunity{user.CommunityId})
		client.Send(&msg.SetNickname{user.Nickname})
		client.Send(&msg.SetSecretQuestion{user.SecretQuestion})
		client.Send(&msg.SetRealmServers{get_realm_servers(ctx)})
		client.Send(&msg.LoginSuccess{IsAdmin: false})

		client.SetUser(user)
		client.State().Inc()
	} else {
		client.Close()
	}
}

func client_handle_players_list(ctx *context, client Client, data []byte) {
	realms := ctx.backend.GetRealms()
	players := make([]*types.RealmServerPlayers, len(realms))
	callback := make(chan types.RealmServerPlayers, len(realms))

	for _, realm := range realms {
		realm.AskPlayers(client.User().Id, callback)
	}

	for i := 0; i < len(realms); i++ {
		p := <-callback
		players = append(players, &p)
	}

	client.Send(&msg.SetRealmServerPlayers{client.User().SubscriptionEnd, players})
}

func client_handle_queue_status(ctx *context, client Client, data []byte) {
	// TODO send queue status
}

func client_handle_realm_selection(ctx *context, client Client, data []byte) {
	m := msg.RealmServerSelectionRequest{}
	m.Deserialize(bytes.NewReader(data))

	if realm, ok := ctx.backend.GetRealm(m.ServerId); ok {
		callback := realm.NotifyUserConnection(client.Ticket(), client.User())

		go func() { // maybe it should be blocking ?
			<-callback
			client.CloseWith(&msg.RealmServerSelectionResponse{
				Address: realm.Address,
				Port:    realm.Port,
				Ticket:  client.Ticket(),
			})
		}()
	} else {
		client.Send(&msg.RealmServerSelectionError{})
	}
}
