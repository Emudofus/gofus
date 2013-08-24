package main

import (
	"database/sql"
	"os"
	"os/signal"
	"fmt"
	"flag"
	"github.com/Blackrush/gofus/login/network"
	_ "github.com/lib/pq"
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

	db, err := sql.Open("postgres", "user=root dbname=gofus")
	if err != nil {
		panic(err.Error())
	}

	networkService := network.New(db, network.Configuration {
		Port: uint16(*port),
		NbWorkers: *nbWorkers,
	})

	go networkService.Start()

	wait_user_input()

	networkService.Stop()
}
