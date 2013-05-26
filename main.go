package main

import (
	"flag"
	"github.com/Blackrush/gofus/game"
	"github.com/Blackrush/gofus/login"
	"log"
	"os"
	"os/signal"
)

type StartStopper interface {
	Start() error
	Stop() error
}

var (
	Server StartStopper
)

func main() {
	flag.Parse()

	Server = get_server()

	if err := Server.Start(); err != nil {
		log.Fatal(err)
	}

	<-wait()

	if err := Server.Stop(); err != nil {
		log.Fatal(err)
	}
}

func get_server() (server StartStopper) {
	switch flag.Arg(0) {
	case "login":
		server = login.NewServer()
	case "game":
		server = game.NewServer()
	}
	return
}

func wait() <-chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, os.Interrupt)
	return c
}
