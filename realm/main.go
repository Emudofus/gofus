package main

import (
	"flag"
	"fmt"
	realmdb "github.com/Blackrush/gofus/realm/db"
	staticdb "github.com/Blackrush/gofus/realm/db/static"
	bnetwork "github.com/Blackrush/gofus/realm/network/backend"
	fnetwork "github.com/Blackrush/gofus/realm/network/frontend"
	"github.com/Blackrush/gofus/shared/db"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
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
	cfg := load_config()

	staticDb := db.Open(&cfg.StaticDatabase)
	defer staticDb.Close()

	maps := staticdb.NewMaps(staticDb)
	maps.Load()

	database := db.Open(&cfg.Database)
	defer database.Close()

	players := realmdb.NewPlayers(database, cfg.Players, maps)
	players.Load()

	bnet := bnetwork.New(players, cfg.Backend)

	go bnet.Start()
	defer bnet.Stop()

	fnet := fnetwork.New(bnet, players, cfg.Frontend)

	go fnet.Start()
	defer fnet.Stop()

	<-wait_for_input()
}
