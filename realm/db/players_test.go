package db

import (
	"fmt"
	"github.com/Blackrush/gofus/shared/db"
	_ "github.com/lib/pq"
	"log"
	"testing"
	"time"
)

const (
	dbuser = "gofus"
	dbname = "gofus"
	dbpass = "lel"
)

func maybe_empty(str string) string {
	if len(str) <= 0 {
		return "''"
	} else {
		return str
	}
}

func TestInsertUser(t *testing.T) {
	database := db.Open(&db.Configuration{
		Driver:         "postgres",
		DataSourceName: fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", maybe_empty(dbuser), maybe_empty(dbname), maybe_empty(dbpass)),
	})
	defer database.Close()
	players := NewPlayers(database)

	player := &Player{OwnerId: 1, Name: fmt.Sprintf("gofus_test_%d", time.Now().Unix())}

	if inserted, success := players.Persist(player); !inserted || !success {
		t.Fatal("can't insert player")
	}

	if !player.persisted {
		t.Fatal("player is not marked as persisted")
	}

	if player.Id == 0 { // player.Id is unsigned
		t.Fatal("player id has not been set")
	}

	log.Printf("player id = %d", player.Id)
}
