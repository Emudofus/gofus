package network

import (
	"log"
)

// TODO: Interface
type task struct {
	client Client
	data   []byte
}

// TODO: Interface
type event struct {
	client Client
	login  bool
}

func worker(ctx *context) {
	log.Print("[network-worker] spawned")

	for ctx.running {
		select {
		case task := <-ctx.tasks:
			handle_client_data(ctx, task.client, string(task.data))
		case event := <-ctx.events:
			if event.login {
				handle_client_connection(ctx, event.client)
			} else {
				handle_client_disconnection(ctx, event.client)
			}
		}
	}
}
