package static

import (
	"database/sql"
	"github.com/Blackrush/gofus/realm/core"
)

type Map struct {
	Id    int
	Date  string
	Key   string
	State *State
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

	if maps, err := maps.find_all(); err == nil {
		for _, m := range maps {
			m.State = NewMapState(m)
			maps.cache[m.Id] = m
		}
	}
}

func (maps *Maps) find_all() (res []*Map, err error) {
	rows, err := maps.db.Query("select id, date, key from maps")
	if err != nil {
		return
	}

	for rows.Next() {
		var mapp Map
		err = rows.Scan(&mapp.Id, &mapp.Date, &mapp.Key)
		if err != nil {
			return nil, err
		}

		res = append(res, &mapp)
	}

	return
}
