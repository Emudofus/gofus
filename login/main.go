package main

import (
	"flag"
	"fmt"
	"github.com/Blackrush/gofus/login/db"
	netlogin "github.com/Blackrush/gofus/login/network/login"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
)

var (
	port      = flag.Int("port", 5555, "the port the server will listen on")
	nbWorkers = flag.Int("nbWorkers", 1, "the number of workers to start")
)

func wait_user_input() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Kill, os.Interrupt)
	return sig
}

func main() {
	fmt.Println(` _______  _______  _______           _______           .-.        .-.           
(  ____ \(  ___  )(  ____ \|\     /|(  ____ \          : :        : :           
| (    \/| (   ) || (    \/| )   ( || (    \/    .--.  : :  .---. : -..  .--.   
| |      | |   | || (__    | |   | || (_____    ; .; ; : :_ :  .; : .. :' .; ;  
| | ____ | |   | ||  __)   | |   | |(_____  )   .__,_; .___;:._.' :_;:_;.__,_; 
| | \_  )| |   | || (      | |   | |      ) |    Blackrush  : :  LOGIN                
| (___) || (___) || )      | (___) |/\____) |               :_;                 
(_______)(_______)|/       (_______)\_______) 
`)

	database := db.Open(&db.Configuration{
		Driver:         "postgres",
		DataSourceName: "user=postgres dbname=gofus password=bla sslmode=disable",
	})
	defer database.Close()

	networkService := netlogin.New(database, netlogin.Configuration{
		Port:      uint16(*port),
		NbWorkers: *nbWorkers,
	})

	go networkService.Start()
	defer networkService.Stop()

	<-wait_user_input()
}
