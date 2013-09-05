package main

import (
	"flag"
	"fmt"
	bnetwork "github.com/Blackrush/gofus/login/network/backend"
	fnetwork "github.com/Blackrush/gofus/login/network/frontend"
	"github.com/Blackrush/gofus/shared/db"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
)

var (
	fport   = flag.Int("fport", 5555, "the port the frontend server will listen on")
	bport   = flag.Int("bport", 5554, "the port the backend server will listen on")
	bpass   = flag.String("bpass", "", "the password used to secure the backend server")
	workers = flag.Int("workers", 1, "the number of workers to start")

	dbuser = flag.String("dbuser", "postgres", "the username used to connect to the PostgreSQL database")
	dbname = flag.String("dbname", "gofus", "the name of the PostgreSQL database")
	dbpass = flag.String("dbpass", "", "the password used to connect to the PostgreSQL database")
)

func wait_user_input() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Kill, os.Interrupt)
	return sig
}

func maybe_empty(str *string) string {
	if str == nil || len(*str) <= 0 {
		return "''"
	} else {
		return *str
	}
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

	database := db.Open(&db.Configuration{
		Driver:         "postgres",
		DataSourceName: fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", maybe_empty(dbuser), maybe_empty(dbname), maybe_empty(dbpass)),
	})
	defer database.Close()

	bnet := bnetwork.New(database, bnetwork.Configuration{
		Port:     uint16(*bport),
		Password: *bpass,
	})

	go bnet.Start()
	defer bnet.Stop()

	fnet := fnetwork.New(database, bnet, fnetwork.Configuration{
		Port:    uint16(*fport),
		Workers: *workers,
	})

	go fnet.Start()
	defer fnet.Stop()

	<-wait_user_input()
}
