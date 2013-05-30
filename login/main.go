package main

import (
	"flag"
	"fmt"
	login "github.com/Nyasu/gofus/login/server"
	"github.com/Nyasu/gofus/shared"
	"log"
	"os"
	"os/signal"
)

var Server shared.StartStopper

func main() {
	flag.Parse()

	Server = login.NewServer()

	check_error(Server.Start())
	wait(true)
	check_error(Server.Stop())
}

func check_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func wait(print bool) {
	hook := make(chan os.Signal, 1)
	signal.Notify(hook, os.Kill, os.Interrupt)

	if print {
		fmt.Println("press C-c to shutdown")
	}

	<-hook
}
