package main

import (
	"flag"
	"fmt"
	"github.com/Blackrush/gofus/realm/network"
	"os"
	"os/signal"
)

var (
	port    = flag.Int("port", 5556, "the port to listen on")
	workers = flag.Int("workers", 1, "the number of workers to spawn")
)

func wait_for_input() <-chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, os.Interrupt)
	return c
}

func main() {
	flag.Parse()
	fmt.Println(` _______  _______  _______           _______           .-.        .-.           
(  ____ \(  ___  )(  ____ \|\     /|(  ____ \          : :        : :           
| (    \/| (   ) || (    \/| )   ( || (    \/    .--.  : :  .---. : -..  .--.   
| |      | |   | || (__    | |   | || (_____    ; .; ; : :_ :  .; : .. :' .; ;  
| | ____ | |   | ||  __)   | |   | |(_____  )   .__,_; .___;:._.' :_;:_;.__,_; 
| | \_  )| |   | || (      | |   | |      ) |    Blackrush  : :  REALM               
| (___) || (___) || )      | (___) |/\____) |               :_;                 
(_______)(_______)|/       (_______)\_______) 
`)

	networkService := network.New(network.Configuration{
		Port:    uint16(*port),
		Workers: *workers,
	})

	go networkService.Start()
	defer networkService.Stop()

	<-wait_for_input()
}
