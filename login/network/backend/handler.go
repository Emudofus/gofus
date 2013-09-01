package backend

import (
	"github.com/Blackrush/gofus/protocol/backend"
)

func client_handle_data(ctx *context, client *Client, arg backend.Message) {
	switch msg := arg.(type) {
	case *backend.AuthReqMsg:
		client_handle_auth(ctx, client, msg)
	}
}

func client_handle_auth(ctx *context, client *Client, msg *backend.AuthReqMsg) {

}
