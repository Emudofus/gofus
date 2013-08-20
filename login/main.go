package main

import (
	"os"
	"os/signal"
	"fmt"
	"flag"
	"github.com/Blackrush/gofus/login/network"
)

var (
	port = flag.Int("port", 5555, "the port the server will listen on")
	nbWorkers = flag.Int("nbWorkers", 1, "the number of workers to start")
)

func wait_user_input() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)
	<-sig
}

func main() {
	fmt.Println("/==============\\")
	fmt.Println("| PHOTON ALPHA |")
	fmt.Println("\\==============/")
	fmt.Println()

	networkService := network.New(network.Configuration {
		Port: uint16(*port),
		NbWorkers: *nbWorkers,
	})

	go networkService.Start()

	wait_user_input()

	networkService.Stop()
}
