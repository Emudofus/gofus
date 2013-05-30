package server

import (
	"flag"
	"github.com/Nyasu/gofus/login/net"
	"github.com/Nyasu/gofus/shared"
	"log"
)

var (
	debug = flag.Bool("debug", false, "will modify default behavior in some circumstances")
)

type server struct {
	net shared.Server
}

// Create a new login server
func NewServer() shared.StartStopper {
	return &server{
		net.NewServer(*debug),
	}
}

// Start login server without blocking
func (server *server) Start() (err error) {
	if err = shared.Start(server.net); err != nil {
		return
	}

	log.Print("login server started without error")

	return nil
}

// Stop login server without blocking
func (server *server) Stop() (err error) {
	if err = shared.Stop(server.net); err != nil {
		return
	}

	log.Print("login server stopped without error")

	return nil
}
