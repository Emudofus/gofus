package login

import (
	"flag"
	"github.com/Nyasu/gofus/shared"
	login "github.com/Nyasu/gofus/login/server"
	"log"
	"os"
	"os/signal"
)

var Server shared.StartStopper

func main() {
	flag.Parse()

	Server = login.NewServer()

	check_error(Server.Start())
	defer check_error(Server.Stop())

	wait()
}

func check_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func wait() {
	hook := make(chan os.Signal, 1)
	signal.Notify(hook, os.Kill, os.Interrupt)
	<-hook
}
