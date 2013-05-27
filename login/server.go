package login

import (
	"github.com/Blackrush/gofus/shared"
	"flag"
)

var (
	debug = flag.Bool("debug", false, "will modify default behavior in some circumstances")
)

type server struct {
	net *network
}

// Create a new login server
func NewServer() shared.StartStopper {
	return &server{
		new_network(),
	}
}

// Start login server without blocking
func (server *server) Start() (err error) {
	if err = server.net.Start(); err != nil {
		return
	}

	return nil
}

// Stop login server without blocking
func (server *server) Stop() (err error) {
	if err = server.net.Stop(); err != nil {
		return
	}

	return nil
}
