package db

import (
	"database/sql"
	"fmt"
	"github.com/Blackrush/gofus/realm/db/static"
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
		return fmt.Sprintf("%x", int(accessory))
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
	return fmt.Sprintf("%v,%v,%v,%v,%v", a[0], a[1], a[2], a[3], a[4])
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

type PlayerPosition struct {
	Map  *static.Map
	Cell *static.MapCell
}

type Player struct {
	Id         uint64
	OwnerId    uint
	Name       string
	Appearance PlayerAppearance
	Experience PlayerExperience
	Position   PlayerPosition

	persisted bool
}

type PlayersConfig struct {
	StartMap  int
	StartCell uint16
}

type Players struct {
	db     *sql.DB
	config PlayersConfig
	maps   *static.Maps

	players             []*Player
	players_by_owner_id playersByOwnerId
}

type playersByOwnerId map[uint][]*Player

func (p *playersByOwnerId) add_player(player *Player) {
	players, _ := (*p)[player.OwnerId]
	(*p)[player.OwnerId] = append(players, player)
}

func (p *playersByOwnerId) rem_player(player *Player) {
	players, _ := (*p)[player.OwnerId]
	if players_remove(&players, player) {
		(*p)[player.OwnerId] = players
	}
}

func NewPlayers(db *sql.DB, config PlayersConfig, maps *static.Maps) *Players {
	return &Players{
		db:     db,
		config: config,
		maps:   maps,
	}
}

func (p *Players) Load() {
	if p.players_by_owner_id != nil {
		panic("player repository already loaded")
	}

	p.players_by_owner_id = make(map[uint][]*Player)

	if players, success := p.find_all(); success {
		for _, player := range players {
			p.add_player(player)
		}
	} else {
		panic("can't load player repository")
	}

	log.Printf("[database] %d players available", len(p.players))
}

func (p *Players) GetById(id uint64) (*Player, bool) {
	for _, player := range p.players {
		if player.Id == id {
			return player, true
		}
	}
	return nil, false
}

func (p *Players) GetByOwnerId(ownerId uint) ([]*Player, bool) {
	players, ok := p.players_by_owner_id[ownerId]
	return players, ok
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
	p.players_by_owner_id.add_player(player)
}

func (p *Players) rem_player(player *Player) {
	players_remove(&p.players, player)
	p.players_by_owner_id.rem_player(player)
}

func player_values(player *Player, with_id, id_last bool) []interface{} {
	var result []interface{}
	if with_id && !id_last {
		result = append(result, player.Id)
	}
	result = append(result,
		player.OwnerId,
		player.Name,
		player.Appearance.Skin,
		player.Appearance.Colors.First,
		player.Appearance.Colors.Second,
		player.Appearance.Colors.Third,
		player.Experience.Level,
		player.Experience.Experience,
		player.Position.Map.Id,
		player.Position.Cell.Id,
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

func (p *Players) find_all() ([]*Player, bool) {
	rows, err := p.db.Query("select id, owner_id, name, skin, first_color, second_color, third_color, level, experience, current_map, current_cell from players")
	if err != nil {
		log.Print(err)
		return nil, false
	}

	var result []*Player
	for rows.Next() {
		player := &Player{persisted: true}

		var currentMap int
		var currentCell uint16

		if err := rows.Scan(append(player_ptrvalues(player), &currentMap, &currentCell)...); err != nil {
			log.Print(err)
			return nil, false
		}

		if m, ok := p.maps.GetById(currentMap); ok {
			player.Position.Map = m
			player.Position.Cell = m.GetCell(currentCell)
		} else {
			panic(fmt.Sprintf("unknown map %d", currentMap))
		}

		result = append(result, player)
	}

	return result, true
}

func (p *Players) NewPlayer(userId uint, name string, breed int, firstColor, secondColor, thirdColor int) *Player {
	player := &Player{
		OwnerId: userId,
		Name:    name,
	}

	player.Appearance.Skin = breed * 10
	player.Appearance.Colors.First = PlayerColor(firstColor)
	player.Appearance.Colors.Second = PlayerColor(secondColor)
	player.Appearance.Colors.Third = PlayerColor(thirdColor)

	player.Experience = PlayerExperience{
		Level:      1,
		Experience: 0,
	}
	m, _ := p.maps.GetById(p.config.StartMap)
	player.Position = PlayerPosition{
		Map:  m,
		Cell: m.GetCell(p.config.StartCell),
	}

	return player
}

func (p *Players) Persist(player *Player) (inserted bool, success bool) {
	if !player.persisted {
		stmt, err := p.db.Prepare("insert into players(owner_id, name, skin, first_color, second_color, third_color, level, experience, current_map, current_cell) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id;")
		//defer stmt.Close()

		if err != nil {
			log.Print(err)
			return
		}

		res, err := stmt.Query(player_values(player, false, false)...)

		if err != nil {
			log.Print(err)
			return
		}
		if res.Err() != nil {
			log.Print(err.Error())
			return
		}
		if !res.Next() {
			log.Print("the database did not returned any values")
			return
		}

		var id uint64
		if err := res.Scan(&id); err == nil {
			player.Id = id
		} else {
			log.Print(err)
			return
		}

		player.persisted = true
		p.add_player(player)
		inserted = true
	} else {
		stmt, err := p.db.Prepare("update players set owner_id=$1, name=$2, skin=$3, first_color=$4, second_color=$5, third_color=$6, level=$7, experience=$8, current_map=$9, current_cell=$10 where id=$11")
		//defer stmt.Close()

		if err != nil {
			log.Print(err)
			return
		}

		if _, err := stmt.Exec(player_values(player, true, true)...); err != nil {
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
		defer stmt.Close()

		if _, err := stmt.Exec(player.Id); err == nil {
			player.persisted = false
			p.rem_player(player)

			return true
		} else {
			log.Print(err)
		}
	} else {
		log.Print(err)
	}
	return false
}
