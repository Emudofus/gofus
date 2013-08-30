package network

import (
	"bytes"
	"github.com/Blackrush/gofus/protocol/msg"
	"io"
	"log"
)

type Handler func(ctx *context, client Client, in io.Reader)
type SubHandlers map[string]Handler
type Handlers map[ClientState]SubHandlers

var (
	handlers = make(Handlers)
)

func get_handler(state ClientState, opcode string) (Handler, bool) {
	sub, ok := handlers[state]
	if !ok {
		return nil, false
	}
	handler, ok := sub[opcode]
	if !ok {
		return nil, false
	}
	return handler, true
}

func set_handler(state ClientState, opcode string, handler Handler) {
	sub, ok := handlers[state]
	if !ok {
		sub = make(SubHandlers)
		handlers[state] = sub
	} else if _, ok := sub[opcode]; ok {
		log.Printf("[network-handler] %s handler already exists: overriding", opcode)
	}

	sub[opcode] = handler
}

func handle_client_connection(ctx *context, client Client) {
	log.Printf("[network-client-%04d] connected", client.Id())

	client.Send(&msg.HelloGame{})
	client.SetState(ClientAuthState)
}

func handle_client_disconnection(ctx *context, client Client) {
	log.Printf("[network-client-%04d] disconnected", client.Id())
}

func handle_client_data(ctx *context, client Client, data []byte) {
	log.Printf("[network-client-%04d] RCV(%03d) %s", client.Id(), len(data), data)

	opcode := string(data[:2])

	if handler, ok := get_handler(client.State(), opcode); ok {
		handler(ctx, client, bytes.NewReader(data))
	} else {
		log.Printf("[network-client-%04d] unknown opcode %s", opcode)
	}
}

func init() {
	// TODO yea, this is really odd. needs work...
	set_handler(ClientAuthState, "AT", func(ctx *context, client Client, in io.Reader) {
		msg := &msg.RealmLogin{}
		msg.Deserialize(in)
		handle_client_realm_login(ctx, client, msg)
	})
}

func handle_client_realm_login(ctx *context, client Client, msg *msg.RealmLogin) {
	log.Printf("%+v", *msg)
}
