package db

import (
	"database/sql"
	"fmt"
	"log"
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

	persisted bool
}

type Players struct {
	db *sql.DB

	players             []*Player
	players_by_owner_id map[uint][]*Player
}

func NewPlayers(db *sql.DB) *Players {
	return &Players{
		db:                  db,
		players:             make([]*Player, 0, 10),
		players_by_owner_id: make(map[uint][]*Player),
	}
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

func player_values(player *Player, with_id, id_last bool) []interface{} {
	result := []interface{}
	if with_id && !id_last {
		result = append(result, player.Id)
	}
	result = append(result,
		player.Id,
		player.OwnerId,
		player.Name,
		player.Appearance.Skin,
		player.Appearance.Colors.First,
		player.Appearance.Colors.Second,
		player.Appearance.Colors.Third,
		player.Experience.Level,
		player.Experience.Experience,
	)
	if with_id && id_last {
		result = append(result, player.Id)
	}
	return result
}

func player_ptrvalues(player *Player) []interface{} {
	return []interface{}{
		&player.Id,
		&player.OwnerId,
		&player.Name,
		&player.Appearance.Skin,
		&player.Appearance.Colors.First,
		&player.Appearance.Colors.Second,
		&player.Appearance.Colors.Third,
		&player.Experience.Level,
		&player.Experience.Experience,
	}
}

func (p *Players) FindAll() ([]*Player, bool) {
	rows, err := p.db.Query("select id, owner_id, name, skin, first_color, second_color, third_color, level, experience from players")
	if err != nil {
		log.Print(err)
		return nil, false
	}

	result := []*Player{}
	for rows.Next() {
		player := &Player{persisted: true}
		if err := rows.Scan(player_ptrvalues(player)); err != nil {
			log.Print(err)
			return nil, false
		}
		result = append(result, player)
	}

	return result, true
}

func (p *Players) Persist(player *Player) (inserted bool, success bool) {
	if !player.persisted {
		stmt, err := p.db.Prepare("insert into players(owner_id, name, skin, first_color, second_color, third_color, level, experience) values($1, $2, $3, $4, $5, $6, $7, $8")
		if err != nil {
			log.Print(err)
			return
		}

		res, err := stmt.Exec(player_values(player, false, false))
		if err != nil {
			log.Print(err)
			return
		}

		if id, err := res.LastInsertId(); err == nil {
			player.Id = id
		} else {
			log.Print(err)
			return
		}

		player.persisted = true
		inserted = true
	} else {
		stmt, err := p.db.Prepare("update players set owner_id=$1, name=$2, skin=$3, first_color=$4, second_color=$5, third_color=$6, level=$7, experience=$8 where id=$9")
		if err != nil {
			log.Print(err)
			return
		}

		if _, err := stmt.Exec(player_values(player, true, true)); err != nil {
			log.Print(err)
			return
		}
	}

	success = true
	return
}

func (p *Players) Remove(player *Player) (success bool) {
	if !player.persisted {
		return false
	}

	if stmt, err := p.db.Prepare("delete from players where id=$1"); err == nil {
		if _, err := stmt.Exec(player.Id); err == nil {
			return true
		} else {
			log.Print(err)
		}
	} else {
		log.Print(err)
	}
	return false
}