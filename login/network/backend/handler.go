package backend

import (
	"bytes"
	"github.com/Blackrush/gofus/protocol/backend"
	frontend "github.com/Blackrush/gofus/protocol/frontend/types"
	"log"
)

type Realm struct {
	frontend.RealmServer
	Address string
	Port    uint16

	client *Client
}

func NewRealm(client *Client) *Realm {
	realm := new(Realm)
	realm.client = client
	return realm
}

func client_handle_data(ctx *context, client *Client, arg backend.Message) {
	switch msg := arg.(type) {
	case *backend.AuthReqMsg:
		client_handle_auth(ctx, client, msg)
	case *backend.SetInfosMsg:
		client_handle_set_infos(ctx, client, msg)
	}
}

func client_authenticate(ctx *context, client *Client, credentials []byte) (*Realm, bool) {
	if bytes.Equal(ctx.get_password_hash(client.salt), credentials) {
		return NewRealm(client), true
	}
	return nil, false
}

func client_handle_auth(ctx *context, client *Client, msg *backend.AuthReqMsg) {
	if client.realm != nil {
		log.Printf("[realm-%04d] tried to reauth", client.realm.Id)
		return
	}

	if _, exists := ctx.realms[int(msg.Id)]; exists {
		goto failure
	} else if realm, ok := client_authenticate(ctx, client, msg.Credentials); ok {
		ctx.realms[realm.Id] = realm
		client.realm = realm

		client.Send(&backend.AuthRespMsg{Success: true})
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

	log.Printf("[realm-%04d] updated his infos", client.realm.Id)
}
