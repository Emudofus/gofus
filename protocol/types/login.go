package types

import (
	"fmt"
)

type RealmServerState int

const (
	RealmOfflineState RealmServerState = iota
	RealmOnlineState
	RealmSavingState
)

type RealmServer struct {
	Id         int
	State      RealmServerState
	Completion int
	Joinable   bool
}

func (server *RealmServer) Format(state fmt.State, c rune) {
	fmt.Fprintf(state, "%d;%d;%d;%d", server.Id, server.State, server.Completion, btoi(server.Joinable))
}

type RealmServerPlayers struct {
	Id      int
	Players int
}

func (server *RealmServerPlayers) Format(state fmt.State, c rune) {
	fmt.Fprintf(state, "%d,%d", server.Id, server.Players)
}
