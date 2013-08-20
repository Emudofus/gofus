package db

import (
//	"database/sql"
	"fmt"
)

type UserRight uint64

func (right UserRight) With(with UserRight) UserRight {
	return right | with
}

func (right UserRight) Without(without UserRight) UserRight {
	return right & (without & 0xFF)
}

func (right UserRight) Has(other UserRight) bool {
	return (right & other) != 0
}

const (
	NoneRight UserRight = 1 << iota
	LoginRight = 1 << iota
)

const AllRight = NoneRight | LoginRight

type User struct {
	Id uint
	Name string
	Password string
	Nickname string
	SecretQuestion string
	SecretAnswer string
	Rights UserRight
}

func (user *User) String() string {
	return fmt.Sprintf("User{ Id: %d, Name: %s, Nickname: %s, Rights: %d }", user.Id, user.Name, user.Nickname, user.Rights)
}
