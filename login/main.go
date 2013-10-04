package main

import (
	"flag"
	"fmt"
	logindb "github.com/Blackrush/gofus/login/db"
	bnetwork "github.com/Blackrush/gofus/login/network/backend"
	fnetwork "github.com/Blackrush/gofus/login/network/frontend"
	"github.com/Blackrush/gofus/shared/db"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
)

func wait_user_input() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Kill, os.Interrupt)
	return sig
}

func main() {
	flag.Parse()
	fmt.Println(` _______  _______  _______           _______           .-.        .-.           
(  ____ \(  ___  )(  ____ \|\     /|(  ____ \          : :        : :           
| (    \/| (   ) || (    \/| )   ( || (    \/    .--.  : :  .---. : -..  .--.   
| |      | |   | || (__    | |   | || (_____    ; .; ; : :_ :  .; : .. :' .; ;  
| | ____ | |   | ||  __)   | |   | |(_____  )   .__,_; .___;:._.' :_;:_;.__,_; 
| | \_  )| |   | || (      | |   | |      ) |    Blackrush  : :  LOGIN                
| (___) || (___) || )      | (___) |/\____) |               :_;                 
(_______)(_______)|/       (_______)\_______) 
`)
	cfg := load_config()

	database := db.Open(&cfg.Database)
	defer database.Close()

	users := &logindb.Users{database}
	if err := users.ResetCurrentRealm(); err != nil {
		panic(err)
	}

	bnet := bnetwork.New(users, cfg.Backend)

	go bnet.Start()
	defer bnet.Stop()

	fnet := fnetwork.New(users, bnet, cfg.Frontend)

	go fnet.Start()
	defer fnet.Stop()

	<-wait_user_input()
}
