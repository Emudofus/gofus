package main

import (
	"flag"
	"fmt"
	bnetwork "github.com/Blackrush/gofus/realm/network/backend"
	fnetwork "github.com/Blackrush/gofus/realm/network/frontend"
	"github.com/Blackrush/gofus/shared/db"
	"os"
	"os/signal"
)

var (
	fport      = flag.Int("fport", 5556, "the port the frontend server will listen on")
	bid        = flag.Uint("id", 1, "the id of the realm server")
	addr       = flag.String("addr", "127.0.0.1", "the address of the realm server")
	completion = flag.Int("completion", 0, "the completion of the realm server")
	community  = flag.Int("community", 0, "the community id of the realm server")
	bladdr     = flag.String("bladdr", ":5554", "the address and port the backend service will connect to")
	bpass      = flag.String("bpass", "", "the password used to secure backend service")
	workers    = flag.Int("workers", 1, "the number of workers to spawn")

	dbuser = flag.String("dbuser", "postgres", "the username used to connect to the PostgreSQL database")
	dbname = flag.String("dbname", "gofus", "the name of the PostgreSQL database")
	dbpass = flag.String("dbpass", "", "the password used to connect to the PostgreSQL database")
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

	database := db.Open(&db.Configuration{
		Driver:         "postgres",
		DataSourceName: fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", *dbuser, *dbname, *dbpass),
	})

	bnet := bnetwork.New(bnetwork.Configuration{
		ServerId:         *bid,
		ServerAddr:       *addr,
		ServerPort:       uint16(*fport),
		ServerCompletion: *completion,
		Laddr:            *bladdr,
		Password:         *bpass,
	})

	go bnet.Start()
	defer bnet.Stop()

	fnet := fnetwork.New(database, bnet, fnetwork.Configuration{
		Port:        uint16(*fport),
		Workers:     *workers,
		CommunityId: *community,
	})

	go fnet.Start()
	defer fnet.Stop()

	<-wait_for_input()
}
