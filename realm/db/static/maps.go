package static

import (
	"database/sql"
	"github.com/Blackrush/gofus/realm/core"
	"log"
)

type Map struct {
	Id  int
	Pos struct {
		X int
		Y int
	}
	Width   int
	Height  int
	Subarea int
	Data    string
	Key     string
	Date    string
	Premium bool
	Places  string
	State   *MapState
	Cells   []*MapCell
}

func (m *Map) GetCell(id uint16) *MapCell {
	return m.Cells[id-1]
}

func (m *Map) valptrs() []interface{} {
	return []interface{}{
		&m.Id,
		&m.Pos.X,
		&m.Pos.Y,
		&m.Width,
		&m.Height,
		&m.Subarea,
		&m.Data,
		&m.Key,
		&m.Date,
		&m.Premium,
		&m.Places,
	}
}

type MapCell struct {
	Map *Map
	Id  uint16
}

func ParseMapCell(m *Map, id uint16, data string) *MapCell {
	return &MapCell{
		Map: m,
		Id:  id,
	}
}

func ParseMapCells(m *Map) (res []*MapCell) {
	for i := 0; i < len(m.Data); i += 10 {
		res = append(res, ParseMapCell(m, uint16(i/10)+1, m.Data[i:i+10]))
	}
	return
}

type Actors map[uint64]core.Actor

type MapState struct {
	Map    *Map
	Actors Actors
}

func NewMapState(mapp *Map) *MapState {
	return &MapState{
		Map:    mapp,
		Actors: make(Actors),
	}
}

type Maps struct {
	db *sql.DB

	cache map[int]*Map
}

func NewMaps(db *sql.DB) *Maps {
	return &Maps{
		db: db,
	}
}

func (maps *Maps) Load() {
	if maps.cache != nil {
		panic("map repository already loaded")
	}

	maps.cache = make(map[int]*Map)

	if all, err := maps.find_all(); err == nil {
		for _, m := range all {
			m.Cells = ParseMapCells(m)
			m.State = NewMapState(m)
			maps.cache[m.Id] = m
		}
	} else {
		panic(err.Error())
	}

	log.Printf("[static-database] %d maps available", len(maps.cache))
}

func (maps *Maps) GetById(id int) (*Map, bool) {
	m, ok := maps.cache[id]
	return m, ok
}

func (maps *Maps) find_all() (res []*Map, err error) {
	rows, err := maps.db.Query("select id, posx, posy, width, height, subarea, data, key, date, premium, places from maps")
	if err != nil {
		return
	}

	for rows.Next() {
		var m Map
		if err = rows.Scan(m.valptrs()...); err != nil {
			return nil, err
		}

		res = append(res, &m)
	}

	return
}
