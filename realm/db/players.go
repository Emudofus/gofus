package db

import (
	"database/sql"
	"fmt"
)

type PlayerExperience struct {
	Level      int
	Experience uint64
}

type PlayerColor int

func (color PlayerColor) String() string {
	if int(color) != -1 {
		return fmt.Sprintf("%x", color)
	}
	return ""
}

type PlayerAccessory int

func (accessory PlayerAccessory) String() string {
	if int(accessory) != -1 {
		return fmt.Sprintf("%x", accessory)
	}
	return ""
}

type PlayerAccessories [5]int

func (a *PlayerAccessories) Weapon() int {
	return a[0]
}
func (a *PlayerAccessories) SetWeapon(id int) {
	a[0] = id
}
func (a *PlayerAccessories) Hat() int {
	return a[1]
}
func (a *PlayerAccessories) SetHat(id int) {
	a[1] = id
}
func (a *PlayerAccessories) Cloak() int {
	return a[2]
}
func (a *PlayerAccessories) SetCloak(id int) {
	a[2] = id
}
func (a *PlayerAccessories) Pet() int {
	return a[3]
}
func (a *PlayerAccessories) SetPet(id int) {
	a[3] = id
}
func (a *PlayerAccessories) Shield() int {
	return a[4]
}
func (a *PlayerAccessories) SetShield(id int) {
	a[4] = id
}

func (a PlayerAccessories) String() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s", a[0], a[1], a[2], a[3], a[4])
}

type PlayerAppearance struct {
	Skin   int
	Colors struct {
		First  PlayerColor
		Second PlayerColor
		Third  PlayerColor
	}
	Accessories PlayerAccessories
}

type Player struct {
	Id         uint64
	OwnerId    uint
	Name       string
	Appearance PlayerAppearance
	Experience PlayerExperience
}

type Players struct {
	db *sql.DB

	players             []*Player
	players_by_owner_id map[uint][]*Player
}

func (p *Players) GetByOwnerId(ownerId uint) ([]*Player, bool) {
	player, ok := p.players_by_owner_id[ownerId]
	return player, ok
}

func players_index_of(players []*Player, player *Player) (int, bool) {
	for i, p := range players {
		if p == player {
			return i, true
		}
	}
	return 0, false
}

func players_remove(players *[]*Player, player *Player) (ok bool) {
	if index, ok := players_index_of(*players, player); ok {
		//https://code.google.com/p/go-wiki/wiki/SliceTricks
		copy((*players)[index:], (*players)[index+1:])
		(*players)[len(*players)-1] = nil
		*players = (*players)[:len(*players)-1]
	}
	return
}

func (p *Players) add_player(player *Player) {
	p.players = append(p.players, player)

	var players []*Player
	var ok bool
	if players, ok = p.players_by_owner_id[player.OwnerId]; !ok {
		players = make([]*Player, 0, 10)
	}
	players = append(players, player)
	p.players_by_owner_id[player.OwnerId] = players
}

func (p *Players) rem_player(player *Player) {
	players_remove(&p.players, player)
	if players, ok := p.players_by_owner_id[player.OwnerId]; ok {
		players_remove(&players, player)
	}
}
