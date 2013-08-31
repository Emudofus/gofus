package realm

import (
	"github.com/Blackrush/gofus/protocol/types"
)

type Realm struct {
	types.RealmServer
	client *client
}
