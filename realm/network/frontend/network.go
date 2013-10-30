package frontend

import (
	pback "github.com/Blackrush/gofus/protocol/backend"
	pfront "github.com/Blackrush/gofus/protocol/frontend"
	"github.com/Blackrush/gofus/realm/db"
	"github.com/Blackrush/gofus/realm/network/backend"
	"github.com/Blackrush/gofus/shared"
	"io"
)

type Configuration struct {
	Port        uint16
	Workers     int
	CommunityId int
	ServerId    uint
}

type Service interface {
	shared.StartStopper

	Config() Configuration
	Backend() backend.Service
	Players() *db.Players
}

type Client interface {
	io.WriteCloser
	pfront.Sender
	pfront.CloseWither
	Alive() bool

	Id() uint64
	UserInfos() *pback.UserInfos
	SetUserInfos(userInfos pback.UserInfos)
	Player() *db.Player
	SetPlayer(player *db.Player)
}
