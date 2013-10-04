package db

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
	NoneRight  UserRight = 1 << iota
	LoginRight           = 1 << iota
)

const AllRight = NoneRight | LoginRight
