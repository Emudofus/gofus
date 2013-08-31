package main

import (
	"flag"
	"fmt"
	lnetwork "github.com/Blackrush/gofus/login/network/login"
	rnetwork "github.com/Blackrush/gofus/login/network/realm"
	"github.com/Blackrush/gofus/shared/db"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
)

var (
	lport     = flag.Int("lport", 5555, "the port the login server will listen on")
	rport     = flag.Int("rport", 5554, "the port the realm server will listen on")
	nbWorkers = flag.Int("nbWorkers", 1, "the number of workers to start")

	dbuser = flag.String("dbuser", "postgres", "the username used to connect to the PostgreSQL database")
	dbname = flag.String("dbname", "gofus", "the name of the PostgreSQL database")
	dbpass = flag.String("dbpass", "", "the password used to connect to the PostgreSQL database")
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

	database := db.Open(&db.Configuration{
		Driver:         "postgres",
		DataSourceName: fmt.Sprint("user=%s dbname=%s password=%s sslmode=disable", *dbuser, *dbname, *dbpass),
	})
	defer database.Close()

	rnet := rnetwork.New(rnetwork.Configuration{
		Port: uint16(*rport),
	})

	go rnet.Start()
	defer rnet.Stop()

	lnet := lnetwork.New(database, lnetwork.Configuration{
		Port:      uint16(*lport),
		NbWorkers: *nbWorkers,
	})

	go lnet.Start()
	defer lnet.Stop()

	<-wait_user_input()
}
